package goja_runtime

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/dop251/goja"
	"github.com/rs/zerolog"
)

type Manager struct {
	pools     map[string]*Pool // Map of extension ID to pool
	mu        sync.RWMutex
	logger    *zerolog.Logger
	maxPerExt int32
}

func NewManager(logger *zerolog.Logger, maxPerExt int32) *Manager {
	return &Manager{
		pools:     make(map[string]*Pool),
		logger:    logger,
		maxPerExt: maxPerExt,
	}
}

// GetOrCreatePool returns an existing pool or creates a new one for an extension
func (m *Manager) GetOrCreatePool(extID string, initFn func() (*goja.Runtime, error)) (*Pool, error) {
	m.mu.RLock()
	if pool, exists := m.pools[extID]; exists {
		m.mu.RUnlock()
		return pool, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check after acquiring write lock
	if pool, exists := m.pools[extID]; exists {
		return pool, nil
	}

	pool := NewPool(m.maxPerExt, initFn, m.logger)
	m.pools[extID] = pool
	return pool, nil
}

// Cleanup cleans up all pools
func (m *Manager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, pool := range m.pools {
		pool.Cleanup()
	}
	m.pools = make(map[string]*Pool)
}

// Pool represents a pool of Goja runtimes for a specific extension
type Pool struct {
	runtimes chan *goja.Runtime
	initFn   func() (*goja.Runtime, error)
	logger   *zerolog.Logger
	active   int32
	maxSize  int32
	metrics  metrics
}

type metrics struct {
	created  atomic.Int64
	reused   atomic.Int64
	errors   atomic.Int64
	timeouts atomic.Int64
}

func NewPool(maxSize int32, initFn func() (*goja.Runtime, error), logger *zerolog.Logger) *Pool {
	return &Pool{
		runtimes: make(chan *goja.Runtime, maxSize),
		initFn:   initFn,
		logger:   logger,
		maxSize:  maxSize,
	}
}

// Get gets a runtime from the pool or creates a new one
func (p *Pool) Get(ctx context.Context) (*goja.Runtime, error) {
	// If an idle runtime exists, use it.
	select {
	case runtime := <-p.runtimes:
		p.metrics.reused.Add(1)
		atomic.AddInt32(&p.active, 1)
		return runtime, nil
	default:
		// Try to reserve a spot for a new runtime atomically.
		for {
			cur := atomic.LoadInt32(&p.active)
			if cur >= p.maxSize {
				// Wait for an idle runtime to become available.
				select {
				case runtime := <-p.runtimes:
					p.metrics.reused.Add(1)
					atomic.AddInt32(&p.active, 1)
					return runtime, nil
				case <-ctx.Done():
					p.metrics.timeouts.Add(1)
					return nil, ctx.Err()
				}
			}
			// Attempt to reserve a slot.
			if atomic.CompareAndSwapInt32(&p.active, cur, cur+1) {
				break
			}
		}
		// Create a new runtime.
		runtime, err := p.initFn()
		if err != nil {
			p.metrics.errors.Add(1)
			atomic.AddInt32(&p.active, -1)
			return nil, err
		}
		p.metrics.created.Add(1)
		return runtime, nil
	}
}

// Put returns a runtime to the pool
func (p *Pool) Put(runtime *goja.Runtime) {
	if runtime == nil {
		return
	}

	runtime.ClearInterrupt()
	atomic.AddInt32(&p.active, -1)

	// Try to put back in pool or discard if full
	select {
	case p.runtimes <- runtime:
	default:
		// Pool is full, discard runtime
	}
}

// Cleanup cleans up the pool
func (p *Pool) Cleanup() {
	close(p.runtimes)
	for runtime := range p.runtimes {
		runtime.ClearInterrupt()
	}
}

// Stats returns the pool's metrics.
func (p *Pool) Stats() map[string]int64 {
	return map[string]int64{
		"created":  p.metrics.created.Load(),
		"reused":   p.metrics.reused.Load(),
		"errors":   p.metrics.errors.Load(),
		"timeouts": p.metrics.timeouts.Load(),
	}
}

// PrintMetrics logs metrics for each extension pool managed by the Manager.
func (m *Manager) PrintMetrics() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for extID, pool := range m.pools {
		stats := pool.Stats()
		m.logger.Info().
			Str("extension", extID).
			Int64("created", stats["created"]).
			Int64("reused", stats["reused"]).
			Int64("errors", stats["errors"]).
			Int64("timeouts", stats["timeouts"]).
			Msg("VM Pool Metrics")
	}
}

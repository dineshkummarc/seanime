package extension_repo_test

import (
	"os"
	"seanime/internal/extension"
	"seanime/internal/extension_repo"
	"seanime/internal/util"
	"testing"

	hibikemanga "github.com/5rahim/hibike/pkg/extension/manga"
	hibikeonlinestream "github.com/5rahim/hibike/pkg/extension/onlinestream"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestGojaWithExtension(t *testing.T) {
	// Get the script
	filepath := "./goja_manga_test/my-manga-provider.ts"
	fileB, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	ext := &extension.Extension{
		ID:          "my-manga-provider",
		Name:        "MyMangaProvider",
		Version:     "0.1.0",
		ManifestURI: "",
		Language:    extension.LanguageTypescript,
		Type:        extension.TypeMangaProvider,
		Description: "",
		Author:      "",
		Payload:     string(fileB),
	}

	// Create the provider
	provider, _, err := extension_repo.NewGojaMangaProvider(ext, ext.Language, util.NewLogger())
	require.NoError(t, err)

	// Test the search function
	searchResult, err := provider.Search(hibikemanga.SearchOptions{Query: "dandadan"})
	require.NoError(t, err)

	spew.Dump(searchResult)

	// Should have a result with rating of 1
	var dandadanRes *hibikemanga.SearchResult
	for _, res := range searchResult {
		if res.SearchRating == 1 {
			dandadanRes = res
			break
		}
	}
	require.NotNil(t, dandadanRes)
	spew.Dump(dandadanRes)

	// Test the search function again
	searchResult, err = provider.Search(hibikemanga.SearchOptions{Query: "boku no kokoro no yaibai"})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(searchResult), 1)

	t.Logf("Search results: %d", len(searchResult))

	// Test the findChapters function
	chapters, err := provider.FindChapters("pYN47sZm") // Boku no Kokoro no Yabai Yatsu
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(chapters), 100)

	t.Logf("Chapters: %d", len(chapters))

	// Test the findChapterPages function
	pages, err := provider.FindChapterPages("WLxnx") // Boku no Kokoro no Yabai Yatsu - Chapter 1
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(pages), 10)

	for _, page := range pages {
		t.Logf("Page: %s, Index: %d\n", page.URL, page.Index)
	}
}

func TestGojaOnlinestreamExtension(t *testing.T) {
	// Get the script
	filepath := "./goja_animepahe/animepahe.ts"
	fileB, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	ext := &extension.Extension{
		ID:          "animepahe",
		Name:        "Animepahe",
		Version:     "0.1.0",
		ManifestURI: "",
		Language:    extension.LanguageTypescript,
		Type:        extension.TypeOnlinestreamProvider,
		Description: "",
		Author:      "",
		Payload:     string(fileB),
	}

	// Create the provider
	provider, _, err := extension_repo.NewGojaOnlinestreamProvider(ext, ext.Language, util.NewLogger())
	require.NoError(t, err)

	// Test the search function
	searchResult, err := provider.Search(hibikeonlinestream.SearchOptions{Query: "dandadan"})
	require.NoError(t, err)

	spew.Dump(searchResult)

	// Should have a result with rating of 1
	var dandadanRes *hibikeonlinestream.SearchResult
	dandadanRes = searchResult[0]
	require.NotNil(t, dandadanRes)

	// Test find episodes
	episodes, err := provider.FindEpisodes(dandadanRes.ID)
	require.NoError(t, err)

	util.Spew(episodes)

}

func TestGojaOnlinestreamExtension2(t *testing.T) {
	// Get the script
	filepath := "./goja_animepahe/animepahe.ts"
	fileB, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}

	ext := &extension.Extension{
		ID:          "animepahe",
		Name:        "Animepahe",
		Version:     "0.1.0",
		ManifestURI: "",
		Language:    extension.LanguageTypescript,
		Type:        extension.TypeOnlinestreamProvider,
		Description: "",
		Author:      "",
		Payload:     string(fileB),
	}

	// Create the provider
	provider, _, err := extension_repo.NewGojaOnlinestreamProvider(ext, ext.Language, util.NewLogger())
	require.NoError(t, err)
	// Find first episode server
	server, err := provider.FindEpisodeServer(&hibikeonlinestream.EpisodeDetails{
		Provider: "animepahe",
		ID:       "b3dd4279853b86047935ed3c9f6c26ed23e56c974efe43333aa4019f68412a13$f789184a-0479-abe3-14b9-155f323a125c",
		Number:   1,
		URL:      "",
		Title:    "",
	}, "kwik")
	require.NoError(t, err)

	spew.Dump(server)
}

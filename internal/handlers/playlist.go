package handlers

import (
	"errors"
	"github.com/seanime-app/seanime/internal/library/entities"
	"path/filepath"
	"strings"
)

// HandleCreatePlaylist
//
//	@summary creates a new playlist.
//	@desc This will create a new playlist with the given name and local file paths.
//	@desc The response is ignored, the client should re-fetch the playlists after this.
//	@route /v1/playlist [POST]
//	@returns entities.Playlist
func HandleCreatePlaylist(c *RouteCtx) error {

	type body struct {
		Name  string   `json:"name"`
		Paths []string `json:"paths"`
	}

	var b body
	if err := c.Fiber.BodyParser(&b); err != nil {
		return c.RespondWithError(err)
	}

	// Get the local files
	dbLfs, _, err := c.App.Database.GetLocalFiles()
	if err != nil {
		return c.RespondWithError(err)
	}

	// Filter the local files
	lfs := make([]*entities.LocalFile, 0)
	for _, path := range b.Paths {
		for _, lf := range dbLfs {
			if lf.GetNormalizedPath() == strings.ToLower(filepath.ToSlash(path)) {
				lfs = append(lfs, lf)
				break
			}
		}
	}

	// Create the playlist
	playlist := entities.NewPlaylist(b.Name)
	playlist.SetLocalFiles(lfs)

	// Save the playlist
	if err := c.App.Database.SavePlaylist(playlist); err != nil {
		return c.RespondWithError(err)
	}

	return c.RespondWithData(playlist)
}

// HandleGetPlaylists
//
//	@summary returns all playlists.
//	@route /v1/playlists [GET]
//	@returns []entities.Playlist
func HandleGetPlaylists(c *RouteCtx) error {

	playlists, err := c.App.Database.GetPlaylists()
	if err != nil {
		return c.RespondWithError(err)
	}

	return c.RespondWithData(playlists)
}

// HandleUpdatePlaylist
//
//	@summary updates a playlist.
//	@returns the updated playlist
//	@desc The response is ignored, the client should re-fetch the playlists after this.
//	@route /v1/playlist/{id} [PATCH]
//	@params id - int - true - "The ID of the playlist to update."
//	@returns entities.Playlist
func HandleUpdatePlaylist(c *RouteCtx) error {

	type body struct {
		DbId  uint     `json:"dbId"`
		Name  string   `json:"name"`
		Paths []string `json:"paths"`
	}

	var b body
	if err := c.Fiber.BodyParser(&b); err != nil {
		return c.RespondWithError(err)
	}

	// Get the local files
	dbLfs, _, err := c.App.Database.GetLocalFiles()
	if err != nil {
		return c.RespondWithError(err)
	}

	// Filter the local files
	lfs := make([]*entities.LocalFile, 0)
	for _, path := range b.Paths {
		for _, lf := range dbLfs {
			if lf.GetNormalizedPath() == strings.ToLower(filepath.ToSlash(path)) {
				lfs = append(lfs, lf)
				break
			}
		}
	}

	// Recreate playlist
	playlist := entities.NewPlaylist(b.Name)
	playlist.DbId = b.DbId
	playlist.Name = b.Name
	playlist.SetLocalFiles(lfs)

	// Save the playlist
	if err := c.App.Database.UpdatePlaylist(playlist); err != nil {
		return c.RespondWithError(err)
	}

	return c.RespondWithData(playlist)
}

// HandleDeletePlaylist
//
//	@summary deletes a playlist.
//	@route /v1/playlist [DELETE]
//	@returns bool
func HandleDeletePlaylist(c *RouteCtx) error {

	type body struct {
		DbId uint `json:"dbId"`
	}

	var b body
	if err := c.Fiber.BodyParser(&b); err != nil {
		return c.RespondWithError(err)

	}

	if err := c.App.Database.DeletePlaylist(b.DbId); err != nil {
		return c.RespondWithError(err)
	}

	return c.RespondWithData(true)
}

// HandleGetPlaylistEpisodes
//
//	@summary returns all the local files of a playlist media entry that have not been watched.
//	@route /v1/playlist/episodes/{id}/{progress} [GET]
//	@params id - int - true - "The ID of the media entry."
//	@params progress - int - true - "The progress of the media entry."
func HandleGetPlaylistEpisodes(c *RouteCtx) error {

	lfs, _, err := c.App.Database.GetLocalFiles()
	if err != nil {
		return c.RespondWithError(err)
	}

	lfw := entities.NewLocalFileWrapper(lfs)

	// Params
	mId, err := c.Fiber.ParamsInt("id")
	if err != nil {
		return c.RespondWithError(err)
	}
	progress, err := c.Fiber.ParamsInt("progress")
	if err != nil {
		return c.RespondWithError(err)
	}

	group, found := lfw.GetLocalEntryById(mId)
	if !found {
		return c.RespondWithError(errors.New("media entry not found"))
	}

	toWatch := group.GetUnwatchedLocalFiles(progress)

	return c.RespondWithData(toWatch)
}

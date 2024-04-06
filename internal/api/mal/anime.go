package mal

import (
	"fmt"
	"net/url"
)

const (
	BaseAnimeFields string = "id,title,main_picture,alternative_titles,start_date,end_date,start_season,nsfw,synopsis,num_episodes,mean,rank,popularity,media_type,status"
)

type (
	BasicAnime struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		MainPicture struct {
			Medium string `json:"medium"`
			Large  string `json:"large"`
		} `json:"main_picture"`
		AlternativeTitles struct {
			Synonyms []string `json:"synonyms"`
			En       string   `json:"en"`
			Ja       string   `json:"ja"`
		} `json:"alternative_titles"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		StartSeason struct {
			Year   int    `json:"year"`
			Season string `json:"season"`
		} `json:"start_season"`
		Synopsis    string      `json:"synopsis"`
		NSFW        string      `json:"nsfw"`
		NumEpisodes int         `json:"num_episodes"`
		Mean        float32     `json:"mean"`
		Rank        int         `json:"rank"`
		Popularity  int         `json:"popularity"`
		MediaType   MediaType   `json:"media_type"`
		Status      MediaStatus `json:"status"`
	}
	AnimeListEntry struct {
		Node struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			MainPicture struct {
				Medium string `json:"medium"`
				Large  string `json:"large"`
			} `json:"main_picture"`
		} `json:"node"`
		ListStatus struct {
			Status             MediaListStatus `json:"status"`
			IsRewatching       bool            `json:"is_rewatching"`
			NumEpisodesWatched int             `json:"num_episodes_watched"`
			Score              int             `json:"score"`
			UpdatedAt          string          `json:"updated_at"`
		} `json:"list_status"`
	}
)

func (w *Wrapper) GetAnimeDetails(mId int) (*BasicAnime, error) {
	reqUrl := fmt.Sprintf("%s/anime/%d?fields=%s", ApiBaseURL, mId, BaseAnimeFields)

	if w.AccessToken == "" {
		return nil, fmt.Errorf("access token is empty")
	}

	var anime BasicAnime
	err := w.doQuery("GET", reqUrl, nil, "application/json", &anime)
	if err != nil {
		return nil, err
	}

	return &anime, nil
}

func (w *Wrapper) GetAnimeCollection() ([]*AnimeListEntry, error) {

	reqUrl := fmt.Sprintf("%s/users/@me/animelist?fields=list_status&limit=1000", ApiBaseURL)

	type response struct {
		Data []*AnimeListEntry `json:"data"`
	}

	var data response
	err := w.doQuery("GET", reqUrl, nil, "application/json", &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

type AnimeListProgressParams struct {
	NumEpisodesWatched *int
}

func (w *Wrapper) UpdateAnimeProgress(opts *AnimeListProgressParams, mId int) error {
	// Get anime details
	anime, err := w.GetAnimeDetails(mId)
	if err != nil {
		return err
	}

	status := MediaListStatusWatching
	if anime.Status == MediaStatusFinishedAiring && anime.NumEpisodes == *opts.NumEpisodesWatched {
		status = MediaListStatusCompleted
	}

	// Update MAL list entry
	err = w.UpdateAnimeListStatus(&AnimeListStatusParams{
		Status:             &status,
		NumEpisodesWatched: opts.NumEpisodesWatched,
	}, mId)

	return err
}

type AnimeListStatusParams struct {
	Status             *MediaListStatus
	IsRewatching       *bool
	NumEpisodesWatched *int
	Score              *int
}

func (w *Wrapper) UpdateAnimeListStatus(opts *AnimeListStatusParams, mId int) error {

	reqUrl := fmt.Sprintf("%s/anime/%d/my_list_status", ApiBaseURL, mId)

	// Build URL
	urlData := url.Values{}
	if opts.Status != nil {
		urlData.Set("status", string(*opts.Status))
	}
	if opts.IsRewatching != nil {
		urlData.Set("is_rewatching", fmt.Sprintf("%t", *opts.IsRewatching))
	}
	if opts.NumEpisodesWatched != nil {
		urlData.Set("num_watched_episodes", fmt.Sprintf("%d", *opts.NumEpisodesWatched))
	}
	if opts.Score != nil {
		urlData.Set("score", fmt.Sprintf("%d", *opts.Score))
	}
	encodedData := urlData.Encode()

	err := w.doMutation("PATCH", reqUrl, encodedData)
	if err != nil {
		return err
	}
	return nil
}

func (w *Wrapper) DeleteAnimeListItem(mId int) error {

	reqUrl := fmt.Sprintf("%s/anime/%d/my_list_status", ApiBaseURL, mId)

	err := w.doMutation("DELETE", reqUrl, "")
	if err != nil {
		return err
	}

	return nil
}

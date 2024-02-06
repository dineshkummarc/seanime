package entities

const (
	AutoDownloaderRuleTitleComparisonContains AutoDownloaderRuleTitleComparisonType = "contains"
	AutoDownloaderRuleTitleComparisonLikely   AutoDownloaderRuleTitleComparisonType = "likely"
)

const (
	AutoDownloaderRuleEpisodeNext     AutoDownloaderRuleEpisodeType = "next_episodes"
	AutoDownloaderRuleEpisodeAll      AutoDownloaderRuleEpisodeType = "all_episodes"
	AutoDownloaderRuleEpisodeSelected AutoDownloaderRuleEpisodeType = "selected_episodes"
)

type (
	AutoDownloaderRuleTitleComparisonType string
	AutoDownloaderRuleEpisodeType         string

	AutoDownloaderRule struct {
		DbID                uint                                  `json:"dbId,omitempty"`
		Enabled             bool                                  `json:"enabled"`
		MediaId             int                                   `json:"mediaId"`
		ReleaseGroups       []string                              `json:"releaseGroups,omitempty"`
		Resolutions         []string                              `json:"resolutions"`
		ComparisonTitle     string                                `json:"comparisonTitle"`
		TitleComparisonType AutoDownloaderRuleTitleComparisonType `json:"titleComparisonType"`
		EpisodeType         AutoDownloaderRuleEpisodeType         `json:"episodeType"`
		EpisodeNumbers      []int                                 `json:"episodeNumbers,omitempty"`
	}
)

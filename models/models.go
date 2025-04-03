package models

import (
	"time"
)

type Anime struct {
	ID                string         `json:"id"`
	Slug              string         `json:"slug"`
	Name              string         `json:"name"`
	ReleaseYear       string         `json:"release_year"`
	ReleaseDate       *time.Time     `json:"release_date"`
	Image             string         `json:"image"`
	ImageDesktop      string         `json:"image_desktop"`
	Synopsis          string         `json:"synopsis"`
	Rating            int            `json:"rating"`
	Score             float64        `json:"score"`
	Genres            []string       `json:"genres"`
	AiringDay         string         `json:"airing_day"`
	TotalEpisodes     int            `json:"total_episodes"`
	CurrentSeason     int            `json:"current_season"`
	SeasonNames       map[string]interface{} `json:"season_names"`
	SeasonYears       map[string]interface{} `json:"season_years"`
	AudioType         string         `json:"audio_type"`
	LogoAnime         string         `json:"logo_anime"`
	ThumbnailImage    string         `json:"thumbnail_image"`
	Audio             string         `json:"audio"`
	Subtitles         string         `json:"subtitles"`
	ContentAdvisory   string         `json:"content_advisory"`
	Based             map[string]interface{} `json:"based"`
	IsRelease         bool           `json:"is_release"`
	IsPopularSeason   bool           `json:"is_popular_season"`
	NewReleases       bool           `json:"new_releases"`
	IsPopular         bool           `json:"is_popular"`
	IsNextSeason      bool           `json:"is_next_season"`
	IsThumbnail       bool           `json:"is_thumbnail"`
	IsMovie           bool           `json:"is_movie"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type Episode struct {
	ID          string     `json:"id"`
	AnimeID     string     `json:"anime_id"`
	Season      int        `json:"season"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Duration    string     `json:"duration"`
	Synopsis    string     `json:"synopsis"`
	Image       string     `json:"image"`
	VideoURL    string     `json:"video_url"`
	ReleaseDate *time.Time `json:"release_date"`
	IsLancamento bool      `json:"is_lancamento"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
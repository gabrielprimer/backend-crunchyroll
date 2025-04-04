package models

type Anime struct {
    ID                string         `json:"id"`
    Slug              string         `json:"slug"`
    Name              string         `json:"name"`
    ReleaseYear       string         `json:"release_year"`
    ReleaseDate       string         `json:"release_date,omitempty"` 
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
    CreatedAt         string         `json:"created_at"` 
    UpdatedAt         string         `json:"updated_at"`
}
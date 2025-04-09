package models

import (
	"time"
	"errors"
	"backend-crunchyroll/enums"
)

type Anime struct {
	ID                  string         `json:"id" graphql:"id"`
	PublicCode          string         `json:"public_code" graphql:"publicCode"`
	Slug                string         `json:"slug" graphql:"slug"`
	Name                string         `json:"name" graphql:"name"`
	ReleaseYear         *int           `json:"release_year,omitempty" graphql:"releaseYear"`
	ReleaseDate         *time.Time     `json:"release_date,omitempty" graphql:"releaseDate"`
	ImagePoster         string         `json:"image_poster,omitempty" graphql:"imagePoster"`
	ImageBannerDesktop  string         `json:"image_banner_desktop,omitempty" graphql:"imageBannerDesktop"`
	ImageBannerMobile   string         `json:"image_banner_mobile,omitempty" graphql:"imageBannerMobile"`
	ImageCardCompact    string         `json:"image_card_compact,omitempty" graphql:"imageCardCompact"`
	Synopsis            string         `json:"synopsis,omitempty" graphql:"synopsis"`
	Rating              *int           `json:"rating,omitempty" graphql:"rating"`
	Score               *float64       `json:"score,omitempty" graphql:"score"`
	AiringDay     *enums.AiringDay     `json:"airing_day,omitempty" graphql:"airingDay"`
	TotalEpisodes       *int           `json:"total_episodes,omitempty" graphql:"totalEpisodes"`
	AudioType     *enums.AudioType     `json:"audio_type,omitempty" graphql:"audioType"`
	ImageLogo           string         `json:"image_logo,omitempty" graphql:"imageLogo"`
	ImageThumbnail      string         `json:"image_thumbnail,omitempty" graphql:"imageThumbnail"`
	ContentAdvisory     string         `json:"content_advisory,omitempty" graphql:"contentAdvisory"`
	CreatedAt           time.Time      `json:"created_at" graphql:"createdAt"`
	UpdatedAt           time.Time      `json:"updated_at" graphql:"updatedAt"`
	IsReleasing         bool           `json:"is_releasing" graphql:"isReleasing"`
	IsSeasonPopular     bool           `json:"is_season_popular" graphql:"isSeasonPopular"`
	IsNewRelease        bool           `json:"is_new_release" graphql:"isNewRelease"`
	IsPopular           bool           `json:"is_popular" graphql:"isPopular"`
	HasNextSeason       bool           `json:"has_next_season" graphql:"hasNextSeason"`
	HasThumbnail        bool           `json:"has_thumbnail" graphql:"hasThumbnail"`
	
	// Relacionamentos que podem ser carregados posteriormente
	Genres              []Genre        `json:"genres,omitempty" graphql:"genres"`
	AudioLanguages      []string       `json:"audio_languages,omitempty" graphql:"audioLanguages"` // USER-DEFINED type
	Subtitles           []string       `json:"subtitles,omitempty" graphql:"subtitles"` // USER-DEFINED type
	Seasons             []AnimeSeason  `json:"seasons,omitempty" graphql:"seasons"`
	ContentSources      []ContentSource `json:"content_sources,omitempty" graphql:"contentSources"`
}

type AnimeSeason struct {
	ID               string         `json:"id" graphql:"id"`
	AnimeID          string         `json:"anime_id" graphql:"animeId"`
	SeasonNumber     int            `json:"season_number" graphql:"seasonNumber"`
	SeasonName       *string        `json:"season_name,omitempty" graphql:"seasonName"`
	Status    *enums.AnimeStatus `json:"status,omitempty" graphql:"status"`
	TotalEpisodes    *int           `json:"total_episodes,omitempty" graphql:"totalEpisodes"`
	FirstEpisodeDate *time.Time     `json:"first_episode_date,omitempty" graphql:"firstEpisodeDate"`
	LastEpisodeDate  *time.Time     `json:"last_episode_date,omitempty" graphql:"lastEpisodeDate"`
	Season    *enums.SeasonEnum `json:"season,omitempty" graphql:"season"`
	CreatedAt        time.Time      `json:"created_at" graphql:"createdAt"`
	UpdatedAt        time.Time      `json:"updated_at" graphql:"updatedAt"`
}

type ContentSource struct {
	ID          string    `json:"id" graphql:"id"`
	AnimeID     *string   `json:"anime_id,omitempty" graphql:"animeId"`
	MovieID     *string   `json:"movie_id,omitempty" graphql:"movieId"`
	SourceType enums.SourceType `json:"source_type" graphql:"sourceType"`
	Title       *string   `json:"title,omitempty" graphql:"title"`
	Authors     *string   `json:"authors,omitempty" graphql:"authors"`
	Copyright   *string   `json:"copyright,omitempty" graphql:"copyright"`
	CreatedAt   time.Time `json:"created_at" graphql:"createdAt"`
	UpdatedAt   time.Time `json:"updated_at" graphql:"updatedAt"`
}

type Genre struct {
	ID        string    `json:"id" graphql:"id"`
	Name      string    `json:"name" graphql:"name"`
	CreatedAt time.Time `json:"created_at" graphql:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" graphql:"updatedAt"`
}

// Validação básica
func (a *Anime) Validate() error {
	if a.PublicCode == "" || a.Slug == "" || a.Name == "" {
		return errors.New("public_code, slug e name são obrigatórios")
	}
	return nil
}
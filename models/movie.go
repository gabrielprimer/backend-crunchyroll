package models

import (
	"time"
	"errors"
	"backend-crunchyroll/enums"
)

type Movie struct {
	ID               string                 `json:"id" graphql:"id"`
	PublicCode       string                 `json:"public_code" graphql:"publicCode"`
	Slug             string                 `json:"slug" graphql:"slug"`
	Name             string                 `json:"name" graphql:"name"`
	ReleaseYear      *int                   `json:"release_year,omitempty" graphql:"releaseYear"`
	ReleaseDate      *time.Time             `json:"release_date,omitempty" graphql:"releaseDate"`
	Image            *string                `json:"image,omitempty" graphql:"image"`
	ImageDesktop     *string                `json:"image_desktop,omitempty" graphql:"imageDesktop"`
	Synopsis         *string                `json:"synopsis,omitempty" graphql:"synopsis"`
	Rating           *int                   `json:"rating,omitempty" graphql:"rating"`
	Score            *float64               `json:"score,omitempty" graphql:"score"`
	Duration         *string                `json:"duration,omitempty" graphql:"duration"`
	VideoURL         *string                `json:"video_url,omitempty" graphql:"videoUrl"`
	LogoMovie        *string                `json:"logo_movie,omitempty" graphql:"logoMovie"`
	ThumbnailImage   *string                `json:"thumbnail_image,omitempty" graphql:"thumbnailImage"`
	Audio            *enums.AudioLanguage   `json:"audio,omitempty" graphql:"audio"`
	Subtitles        *string                `json:"subtitles,omitempty" graphql:"subtitles"`
	ContentAdvisory  *string                `json:"content_advisory,omitempty" graphql:"contentAdvisory"`
	IsReleased       bool                   `json:"is_released" graphql:"isReleased"`
	CreatedAt        time.Time              `json:"created_at" graphql:"createdAt"`
	UpdatedAt        time.Time              `json:"updated_at" graphql:"updatedAt"`
	
	// Relacionamentos
	Genres          []Genre         `json:"genres,omitempty" graphql:"genres"`
	ContentSources  []ContentSource `json:"content_sources,omitempty" graphql:"contentSources"`
}

// Validação básica
func (m *Movie) Validate() error {
	if m.PublicCode == "" || m.Slug == "" || m.Name == "" {
		return errors.New("public_code, slug e name são obrigatórios")
	}
	
	if m.Audio != nil && !m.Audio.IsValid() {
		return errors.New("audio language inválido")
	}
	
	return nil
}
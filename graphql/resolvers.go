package graphql

import (
	"context"
	"log"

	"github.com/jefersonprimer/backend-crunchyroll/models"
	"github.com/nedpals/supabase-go"
)

type Resolver struct {
	DB *supabase.Client
}

func (r *Resolver) GetAnimeBySlug(ctx context.Context, args struct{ Slug string }) (*models.Anime, error) {
	var results []models.Anime
	err := r.DB.DB.From("animes").Select("*").Eq("slug", args.Slug).Execute(&results)
	if err != nil {
		log.Printf("Error fetching anime by slug: %v", err)
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return &results[0], nil
}

func (r *Resolver) GetEpisodesByAnime(ctx context.Context, args struct{ AnimeID string }) ([]*models.Episode, error) {
	var episodes []*models.Episode
	err := r.DB.DB.From("episodes").Select("*").Eq("anime_id", args.AnimeID).Execute(&episodes)
	if err != nil {
		log.Printf("Error fetching episodes: %v", err)
		return nil, err
	}
	return episodes, nil
}
package resolvers

import (
	"context"
	"log"

	"github.com/jefersonprimer/backend-crunchyroll/models"
	"github.com/nedpals/supabase-go"
)

type Resolver struct {
	SupabaseClient *supabase.Client
}

func (r *Resolver) GetAnimes(ctx context.Context) ([]*models.Anime, error) {
	var animes []*models.Anime
	err := r.SupabaseClient.DB.From("animes").Select("*").Execute(&animes)
	if err != nil {
		log.Printf("Error fetching animes: %v", err)
		return nil, err
	}
	return animes, nil
}

func (r *Resolver) GetAnimeBySlug(ctx context.Context, args struct{ Slug string }) (*models.Anime, error) {
	var anime models.Anime
	err := r.SupabaseClient.DB.From("animes").Select("*").Eq("slug", args.Slug).Single().Execute(&anime)
	if err != nil {
		log.Printf("Error fetching anime by slug: %v", err)
		return nil, err
	}
	return &anime, nil
}

func (r *Resolver) GetEpisodesByAnime(ctx context.Context, args struct{ AnimeID string }) ([]*models.Episode, error) {
	var episodes []*models.Episode
	err := r.SupabaseClient.DB.From("episodes").Select("*").Eq("anime_id", args.AnimeID).Order("season").Order("release_date").Execute(&episodes)
	if err != nil {
		log.Printf("Error fetching episodes: %v", err)
		return nil, err
	}
	return episodes, nil
}

func (r *Resolver) GetLatestReleases(ctx context.Context) ([]*models.Anime, error) {
	var animes []*models.Anime
	err := r.SupabaseClient.DB.From("animes").Select("*").Eq("is_release", true).Order("updated_at", &supabase.OrderOpts{Descending: true}).Limit(10).Execute(&animes)
	if err != nil {
		log.Printf("Error fetching latest releases: %v", err)
		return nil, err
	}
	return animes, nil
}

func (r *Resolver) GetPopularAnimes(ctx context.Context) ([]*models.Anime, error) {
	var animes []*models.Anime
	err := r.SupabaseClient.DB.From("animes").Select("*").Eq("is_popular", true).Order("score", &supabase.OrderOpts{Descending: true}).Limit(10).Execute(&animes)
	if err != nil {
		log.Printf("Error fetching popular animes: %v", err)
		return nil, err
	}
	return animes, nil
}
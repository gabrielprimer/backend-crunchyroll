package graphql

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"backend-crunchyroll/models"
	"github.com/nedpals/supabase-go"
	"go.uber.org/zap"
)

type Resolver struct {
	DB      *supabase.Client
	logger  *zap.Logger
	cache   *CacheStore
	metrics *ResolverMetrics
}

type CacheStore struct {
	animes        sync.Map // map[string]*models.Anime
	episodes      sync.Map // map[string][]*models.Episode
	animeBySlug   sync.Map // map[string]string (slug -> animeID)
	latestAnimes  []*models.Anime
	popularAnimes []*models.Anime
	latestFetch   time.Time
	popularFetch  time.Time
}

type ResolverMetrics struct {
	cacheHits   int64
	cacheMisses int64
	dbQueries   int64
}

func NewResolver(db *supabase.Client) *Resolver {
	return &Resolver{
		DB:      db,
		logger:  zap.NewExample(),
		cache:   &CacheStore{},
		metrics: &ResolverMetrics{},
	}
}

// GetAnimeBySlug - Versão otimizada com cache e métricas
func (r *Resolver) GetAnimeBySlug(ctx context.Context, args struct{ Slug string }) (*models.Anime, error) {
	if anime, ok := r.getAnimeFromCacheBySlug(args.Slug); ok {
		r.metrics.cacheHits++
		return anime, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var results []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("slug", args.Slug).
		ExecuteWithContext(ctx, &results)

	if err != nil {
		r.logger.Error("Falha ao buscar anime",
			zap.String("slug", args.Slug),
			zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	if len(results) == 0 {
		return nil, nil
	}

	anime := results[0]
	r.cacheAnime(anime)
	r.metrics.dbQueries++

	return anime, nil
}

// GetEpisodesByAnime - Versão otimizada com cache
func (r *Resolver) GetEpisodesByAnime(ctx context.Context, args struct{ AnimeID string }) ([]*models.Episode, error) {
	if episodes, ok := r.cache.episodes.Load(args.AnimeID); ok {
		r.metrics.cacheHits++
		return episodes.([]*models.Episode), nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var episodes []*models.Episode
	err := r.DB.DB.From("episodes").
		Select("*").
		Eq("anime_id", args.AnimeID).
		ExecuteWithContext(ctx, &episodes)

	if err != nil {
		r.logger.Error("Falha ao buscar episódios",
			zap.String("animeID", args.AnimeID),
			zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	r.cache.episodes.Store(args.AnimeID, episodes)
	r.metrics.dbQueries++

	return episodes, nil
}

// GetLatestReleases - Busca animes marcados como lançamentos recentes com cache
func (r *Resolver) GetLatestReleases(ctx context.Context) ([]*models.Anime, error) {
	if !r.cache.latestFetch.IsZero() && time.Since(r.cache.latestFetch) < 5*time.Minute && len(r.cache.latestAnimes) > 0 {
		r.metrics.cacheHits++
		return r.cache.latestAnimes, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var animes []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("new_releases", strconv.FormatBool(true)). // <- CORRIGIDO
		ExecuteWithContext(ctx, &animes)

	if err != nil {
		r.logger.Error("Falha ao buscar lançamentos recentes", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	for _, anime := range animes {
		r.cacheAnime(anime)
	}

	r.cache.latestAnimes = animes
	r.cache.latestFetch = time.Now()
	r.metrics.dbQueries++

	return animes, nil
}

// GetPopularAnimes - Busca animes populares com cache
func (r *Resolver) GetPopularAnimes(ctx context.Context) ([]*models.Anime, error) {
	if !r.cache.popularFetch.IsZero() && time.Since(r.cache.popularFetch) < 5*time.Minute && len(r.cache.popularAnimes) > 0 {
		r.metrics.cacheHits++
		return r.cache.popularAnimes, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var animes []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("is_popular", strconv.FormatBool(true)). // <- CORRIGIDO
		ExecuteWithContext(ctx, &animes)

	if err != nil {
		r.logger.Error("Falha ao buscar animes populares", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	for _, anime := range animes {
		r.cacheAnime(anime)
	}

	r.cache.popularAnimes = animes
	r.cache.popularFetch = time.Now()
	r.metrics.dbQueries++

	return animes, nil
}

// -- Métodos de cache --

func (r *Resolver) cacheAnime(anime *models.Anime) {
	r.cache.animes.Store(anime.ID, anime)
	r.cache.animeBySlug.Store(anime.Slug, anime.ID)
}

func (r *Resolver) getAnimeFromCacheBySlug(slug string) (*models.Anime, bool) {
	if id, ok := r.cache.animeBySlug.Load(slug); ok {
		if anime, ok := r.cache.animes.Load(id.(string)); ok {
			return anime.(*models.Anime), true
		}
	}
	return nil, false
}

// -- Monitoramento --

func (r *Resolver) GetCacheStats() map[string]int64 {
	return map[string]int64{
		"hits":    r.metrics.cacheHits,
		"misses":  r.metrics.cacheMisses,
		"queries": r.metrics.dbQueries,
	}
}

func (r *Resolver) InvalidateCache() {
	r.cache = &CacheStore{}
	r.logger.Info("Cache invalidado")
}

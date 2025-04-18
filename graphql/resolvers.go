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
	releasingAnimes []*models.Anime
	seasonPopularAnimes []*models.Anime
	nextSeasonAnimes []*models.Anime
	hasThumbnail []*models.Anime
	genres        sync.Map // map[string][]*models.Genre
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
		r.logger.Error("Falha ao buscar anime", zap.String("slug", args.Slug), zap.Error(err))
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
		r.logger.Error("Falha ao buscar episódios", zap.String("animeID", args.AnimeID), zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	r.cache.episodes.Store(args.AnimeID, episodes)
	r.metrics.dbQueries++

	return episodes, nil
}

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
		Eq("is_new_release", strconv.FormatBool(true)).
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
		Eq("is_popular", strconv.FormatBool(true)).
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

func (r *Resolver) GetReleasingAnimes(ctx context.Context) ([]*models.Anime, error) {
	if !r.cache.popularFetch.IsZero() && time.Since(r.cache.popularFetch) < 5*time.Minute && len(r.cache.releasingAnimes) > 0 {
		r.metrics.cacheHits++
		return r.cache.releasingAnimes, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var animes []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("is_releasing", strconv.FormatBool(true)).
		ExecuteWithContext(ctx, &animes)

	if err != nil {
		r.logger.Error("Falha ao buscar animes em lancamento", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	for _, anime := range animes {
		r.cacheAnime(anime)
	}

	r.cache.releasingAnimes = animes
	r.cache.popularFetch = time.Now()
	r.metrics.dbQueries++

	return animes, nil
}

func (r *Resolver) GetSeasonPopularAnimes(ctx context.Context) ([]*models.Anime, error) {
	if !r.cache.popularFetch.IsZero() && time.Since(r.cache.popularFetch) < 5*time.Minute && len(r.cache.seasonPopularAnimes) > 0 {
		r.metrics.cacheHits++
		return r.cache.seasonPopularAnimes, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var animes []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("is_season_popular", strconv.FormatBool(true)).
		ExecuteWithContext(ctx, &animes)

	if err != nil {
		r.logger.Error("Falha ao buscar animes de season popular", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	for _, anime := range animes {
		r.cacheAnime(anime)
	}

	r.cache.seasonPopularAnimes = animes
	r.cache.popularFetch = time.Now()
	r.metrics.dbQueries++

	return animes, nil
}

func (r *Resolver) GetNextSeasonAnimes(ctx context.Context) ([]*models.Anime, error) {
	if !r.cache.popularFetch.IsZero() && time.Since(r.cache.popularFetch) < 5*time.Minute && len(r.cache.nextSeasonAnimes) > 0 {
		r.metrics.cacheHits++
		return r.cache.nextSeasonAnimes, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var animes []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("has_next_season", strconv.FormatBool(true)).
		ExecuteWithContext(ctx, &animes)

	if err != nil {
		r.logger.Error("Falha ao buscar animes de temporadas futuras", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	for _, anime := range animes {
		r.cacheAnime(anime)
	}

	r.cache.seasonPopularAnimes = animes
	r.cache.popularFetch = time.Now()
	r.metrics.dbQueries++

	return animes, nil
}

func (r *Resolver) GetHasThumbnail(ctx context.Context) ([]*models.Anime, error) {
	if !r.cache.popularFetch.IsZero() && time.Since(r.cache.popularFetch) < 5*time.Minute && len(r.cache.hasThumbnail) > 0 {
		r.metrics.cacheHits++
		return r.cache.hasThumbnail, nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var animes []*models.Anime
	err := r.DB.DB.From("animes").
		Select("*").
		Eq("has_thumbnail", strconv.FormatBool(true)).
		ExecuteWithContext(ctx, &animes)

	if err != nil {
		r.logger.Error("Falha ao buscar animes de com thumbnails", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	for _, anime := range animes {
		r.cacheAnime(anime)
	}

	r.cache.hasThumbnail = animes
	r.cache.popularFetch = time.Now()
	r.metrics.dbQueries++

	return animes, nil
}

func (r *Resolver) GetAllAnimeNames(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var results []map[string]interface{}
	err := r.DB.DB.From("animes").
		Select("name").
		ExecuteWithContext(ctx, &results)

	if err != nil {
		r.logger.Error("Falha ao buscar nomes de animes", zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	animeNames := make([]string, len(results))
	for i, result := range results {
		if name, ok := result["name"].(string); ok {
			animeNames[i] = name
		} else {
			r.logger.Error("Nome do anime não é uma string", zap.Any("result", result))
			return nil, errors.New("erro interno do servidor: nome do anime inválido")
		}
	}

	r.metrics.dbQueries++

	return animeNames, nil
}

// -- Cache helpers --

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

// -- Monitoring & Cache --

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

func (r *Resolver) GetGenresByAnimeId(ctx context.Context, args struct{ AnimeID string }) ([]*models.Genre, error) {
	if genres, ok := r.cache.genres.Load(args.AnimeID); ok {
		r.metrics.cacheHits++
		return genres.([]*models.Genre), nil
	}
	r.metrics.cacheMisses++

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var genres []*models.Genre
	err := r.DB.DB.From("genres").
		Select("*").
		Eq("anime_id", args.AnimeID).
		ExecuteWithContext(ctx, &genres)

	if err != nil {
		r.logger.Error("Falha ao buscar gêneros", zap.String("animeID", args.AnimeID), zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	r.cache.genres.Store(args.AnimeID, genres)
	r.metrics.dbQueries++

	return genres, nil
}

func (r *Resolver) GetAudioLanguagesByAnimeId(ctx context.Context, args struct{ AnimeID string }) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var results []map[string]interface{}
	err := r.DB.DB.From("animes").
		Select("audio_languages").
		Eq("id", args.AnimeID).
		ExecuteWithContext(ctx, &results)

	if err != nil {
		r.logger.Error("Falha ao buscar audio languages", zap.String("animeID", args.AnimeID), zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	if len(results) == 0 {
		return nil, nil
	}

	audioLanguages, ok := results[0]["audio_languages"].([]string)
	if !ok {
		r.logger.Error("Falha ao converter audio languages", zap.Any("audioLanguages", results[0]["audio_languages"]))
		return nil, errors.New("erro interno do servidor")
	}

	r.metrics.dbQueries++

	return audioLanguages, nil
}

func (r *Resolver) GetSubtitlesByAnimeId(ctx context.Context, args struct{ AnimeID string }) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var results []map[string]interface{}
	err := r.DB.DB.From("animes").
		Select("subtitles").
		Eq("id", args.AnimeID).
		ExecuteWithContext(ctx, &results)

	if err != nil {
		r.logger.Error("Falha ao buscar subtitles", zap.String("animeID", args.AnimeID), zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	if len(results) == 0 {
		return nil, nil
	}

	subtitles, ok := results[0]["subtitles"].([]string)
	if !ok {
		r.logger.Error("Falha ao converter subtitles", zap.Any("subtitles", results[0]["subtitles"]))
		return nil, errors.New("erro interno do servidor")
	}

	r.metrics.dbQueries++

	return subtitles, nil
}

func (r *Resolver) GetSeasonsByAnimeId(ctx context.Context, args struct{ AnimeID string }) ([]*models.AnimeSeason, error) {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // Buscar dados brutos primeiro para evitar problemas de parsing
    var rawSeasons []map[string]interface{}
    err := r.DB.DB.From("anime_seasons").
        Select("*").
        Eq("anime_id", args.AnimeID).
        ExecuteWithContext(ctx, &rawSeasons)

    if err != nil {
        r.logger.Error("Falha ao buscar temporadas", zap.String("animeID", args.AnimeID), zap.Error(err))
        return nil, errors.New("erro interno do servidor")
    }

    // Converter dados brutos para struct AnimeSeason manualmente
    seasons := make([]*models.AnimeSeason, 0, len(rawSeasons))
    for _, rawSeason := range rawSeasons {
        season := &models.AnimeSeason{
            ID:           rawSeason["id"].(string),
            AnimeID:      rawSeason["anime_id"].(string),
            SeasonNumber: int(rawSeason["season_number"].(float64)),
        }

        // Converter campos opcionais com verificação de tipo
        if name, ok := rawSeason["season_name"].(string); ok {
            season.SeasonName = &name
        }
        
        if episodes, ok := rawSeason["total_episodes"].(float64); ok {
            episodesInt := int(episodes)
            season.TotalEpisodes = &episodesInt
        }

        // Pular campos de data problemáticos ou implementar parsing manual
        // A conversão de created_at e updated_at ainda será necessária
        if createdAtStr, ok := rawSeason["created_at"].(string); ok {
            createdAt, err := time.Parse("2006-01-02T15:04:05.999999", createdAtStr)
            if err == nil {
                season.CreatedAt = createdAt
            } else {
                // Fallback para data atual se houver erro
                season.CreatedAt = time.Now()
            }
        }

        if updatedAtStr, ok := rawSeason["updated_at"].(string); ok {
            updatedAt, err := time.Parse("2006-01-02T15:04:05.999999", updatedAtStr)
            if err == nil {
                season.UpdatedAt = updatedAt
            } else {
                // Fallback para data atual se houver erro
                season.UpdatedAt = time.Now()
            }
        }

        seasons = append(seasons, season)
    }

    r.metrics.dbQueries++
    return seasons, nil
}

func (r *Resolver) GetContentSourcesByAnimeId(ctx context.Context, args struct{ AnimeID string }) ([]*models.ContentSource, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var contentSources []*models.ContentSource
	err := r.DB.DB.From("content_sources").
		Select("*").
		Eq("anime_id", args.AnimeID).
		ExecuteWithContext(ctx, &contentSources)

	if err != nil {
		r.logger.Error("Falha ao buscar fontes de conteúdo", zap.String("animeID", args.AnimeID), zap.Error(err))
		return nil, errors.New("erro interno do servidor")
	}

	r.metrics.dbQueries++

	return contentSources, nil
}

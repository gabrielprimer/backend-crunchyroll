package graphql

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast" // <- IMPORT NECESSÁRIO
	"backend-crunchyroll/models"
	"github.com/nedpals/supabase-go"
)


func NewSchema(db *supabase.Client) (graphql.Schema, error) {
	// Use o construtor apropriado para o resolver com cache
	resolver := NewResolver(db)

	// Definição do scalar JSON para dados flexíveis
	jsonScalar := graphql.NewScalar(graphql.ScalarConfig{
		Name:        "JSON",
		Description: "Tipo escalar JSON para representar objetos JSON conforme RFC 7159",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.ObjectValue:
				return valueAST.GetValue()
			}
			return nil
		},
	})

	// Definição de Episode
	episodeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Episode",
		Fields: graphql.Fields{
			"id":           &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"animeId":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"season":       &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
			"title":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"slug":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"duration":     &graphql.Field{Type: graphql.String},
			"synopsis":     &graphql.Field{Type: graphql.String},
			"image":        &graphql.Field{Type: graphql.String},
			"videoUrl":     &graphql.Field{Type: graphql.String},
			"releaseDate":  &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if episode, ok := p.Source.(*models.Episode); ok && episode.ReleaseDate != nil {
					return episode.ReleaseDate.Format(time.RFC3339), nil
				}
				return nil, nil
			}},
			"isLancamento": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"createdAt":    &graphql.Field{Type: graphql.NewNonNull(graphql.String), Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if episode, ok := p.Source.(*models.Episode); ok {
					return episode.CreatedAt.Format(time.RFC3339), nil
				}
				return "", nil
			}},
			"updatedAt":    &graphql.Field{Type: graphql.NewNonNull(graphql.String), Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if episode, ok := p.Source.(*models.Episode); ok {
					return episode.UpdatedAt.Format(time.RFC3339), nil
				}
				return "", nil
			}},
		},
	})

	// Definição de Anime
	animeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Anime",
		Fields: graphql.Fields{
			"id":             &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"slug":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"name":           &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"releaseYear":    &graphql.Field{Type: graphql.String},
			"releaseDate":    &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if anime, ok := p.Source.(*models.Anime); ok && anime.ReleaseDate != nil {
					return anime.ReleaseDate.Format(time.RFC3339), nil
				}
				return nil, nil
			}},
			"image":          &graphql.Field{Type: graphql.String},
			"imageDesktop":   &graphql.Field{Type: graphql.String},
			"synopsis":       &graphql.Field{Type: graphql.String},
			"rating":         &graphql.Field{Type: graphql.Int},
			"score":          &graphql.Field{Type: graphql.Float},
			"genres":         &graphql.Field{Type: graphql.NewList(graphql.NewNonNull(graphql.String))},
			"airingDay":      &graphql.Field{Type: graphql.String},
			"totalEpisodes":  &graphql.Field{Type: graphql.Int},
			"currentSeason":  &graphql.Field{Type: graphql.Int},
			"seasonNames":    &graphql.Field{Type: jsonScalar},
			"seasonYears":    &graphql.Field{Type: jsonScalar},
			"audioType":      &graphql.Field{Type: graphql.String},
			"logoAnime":      &graphql.Field{Type: graphql.String},
			"thumbnailImage": &graphql.Field{Type: graphql.String},
			"audio":          &graphql.Field{Type: graphql.String},
			"subtitles":      &graphql.Field{Type: graphql.String},
			"contentAdvisory": &graphql.Field{Type: graphql.String},
			"based":          &graphql.Field{Type: jsonScalar},
			"isRelease":      &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"isPopularSeason": &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"newReleases":    &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"isPopular":      &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"isNextSeason":   &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"isThumbnail":    &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"isMovie":        &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
			"createdAt":      &graphql.Field{Type: graphql.NewNonNull(graphql.String), Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if anime, ok := p.Source.(*models.Anime); ok {
					return anime.CreatedAt.Format(time.RFC3339), nil
				}
				return "", nil
			}},
			"updatedAt":      &graphql.Field{Type: graphql.NewNonNull(graphql.String), Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if anime, ok := p.Source.(*models.Anime); ok {
					return anime.UpdatedAt.Format(time.RFC3339), nil
				}
				return "", nil
			}},
			"episodes": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(episodeType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					anime := p.Source.(*models.Anime)
					return resolver.GetEpisodesByAnime(p.Context, struct{ AnimeID string }{AnimeID: anime.ID})
				},
			},
		},
	})

	// Tipo para estatísticas de cache
	statsType := graphql.NewObject(graphql.ObjectConfig{
		Name: "CacheStats",
		Fields: graphql.Fields{
			"hits":    &graphql.Field{Type: graphql.Int},
			"misses":  &graphql.Field{Type: graphql.Int},
			"queries": &graphql.Field{Type: graphql.Int},
		},
	})

	// Query root
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"animes": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(animeType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var animes []*models.Anime
					err := resolver.DB.DB.From("animes").Select("*").Execute(&animes)
					if err != nil {
						return nil, err
					}
					// Cache de cada anime
					for _, anime := range animes {
						resolver.cacheAnime(anime)
					}
					resolver.metrics.dbQueries++
					return animes, nil
				},
			},
			"animeBySlug": &graphql.Field{
				Type: animeType,
				Args: graphql.FieldConfigArgument{
					"slug": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					slug := p.Args["slug"].(string)
					return resolver.GetAnimeBySlug(p.Context, struct{ Slug string }{Slug: slug})
				},
			},
			"episodesByAnime": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(episodeType))),
				Args: graphql.FieldConfigArgument{
					"animeId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					animeID := p.Args["animeId"].(string)
					return resolver.GetEpisodesByAnime(p.Context, struct{ AnimeID string }{AnimeID: animeID})
				},
			},
			"latestReleases": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(animeType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.GetLatestReleases(p.Context)
				},
			},
			"popularAnimes": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(animeType))),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.GetPopularAnimes(p.Context)
				},
			},
			"cacheStats": &graphql.Field{
				Type: statsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return resolver.GetCacheStats(), nil
				},
			},
		},
	})

	// Mutation root
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"invalidateCache": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					resolver.InvalidateCache()
					return true, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}
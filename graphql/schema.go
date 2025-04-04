package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/jefersonprimer/backend-crunchyroll/models"
	"github.com/nedpals/supabase-go"
)

func NewSchema(db *supabase.Client) (graphql.Schema, error) {
	resolver := &Resolver{DB: db}

	episodeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Episode",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.String},
			"animeId": &graphql.Field{Type: graphql.String},
			"title":   &graphql.Field{Type: graphql.String},
		},
	})

	animeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Anime",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.String},
			"name":     &graphql.Field{Type: graphql.String},
			"slug":     &graphql.Field{Type: graphql.String},
			"image":    &graphql.Field{Type: graphql.String},
			"episodes": &graphql.Field{
				Type: graphql.NewList(episodeType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					anime := p.Source.(*models.Anime)
					return resolver.GetEpisodesByAnime(p.Context, struct{ AnimeID string }{AnimeID: anime.ID})
				},
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"anime": &graphql.Field{
				Type: animeType,
				Args: graphql.FieldConfigArgument{
					"slug": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					slug, _ := p.Args["slug"].(string)
					return resolver.GetAnimeBySlug(p.Context, struct{ Slug string }{Slug: slug})
				},
			},
			"animes": &graphql.Field{
				Type: graphql.NewList(animeType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var animes []models.Anime
					err := resolver.DB.DB.From("animes").Select("*").Execute(&animes)
					if err != nil {
						return nil, err
					}
					return animes, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
}
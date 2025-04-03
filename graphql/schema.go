package graphql

import (
	"github.com/graphql-go/graphql"
)

func NewSchema() (graphql.Schema, error) {
	// Definindo os tipos GraphQL
	generoType := graphql.NewEnum(graphql.EnumConfig{
		Name: "Genero",
		Values: graphql.EnumValueConfigMap{
			"Acao":         &graphql.EnumValueConfig{Value: "Ação"},
			"Drama":        &graphql.EnumValueConfig{Value: "Drama"},
			"Fantasia":     &graphql.EnumValueConfig{Value: "Fantasia"},
			"Misterio":     &graphql.EnumValueConfig{Value: "Mistério"},
			"Psicologico":  &graphql.EnumValueConfig{Value: "Psicológico"},
			"Suspense":     &graphql.EnumValueConfig{Value: "Suspense"},
		},
	})

	episodioType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Episodio",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"animeId": &graphql.Field{
				Type: graphql.Int,
			},
			"titulo": &graphql.Field{
				Type: graphql.String,
			},
			"numero": &graphql.Field{
				Type: graphql.Int,
			},
			"duracao": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	animeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Anime",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"titulo": &graphql.Field{
				Type: graphql.String,
			},
			"descricao": &graphql.Field{
				Type: graphql.String,
			},
			"ano": &graphql.Field{
				Type: graphql.Int,
			},
			"generos": &graphql.Field{
				Type: graphql.NewList(generoType),
			},
			"episodios": &graphql.Field{
				Type: graphql.NewList(episodioType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					anime := p.Source.(Anime)
					return GetEpisodiosByAnimeID(anime.ID), nil
				},
			},
		},
	})

	// Definindo as queries
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"animes": &graphql.Field{
				Type: graphql.NewList(animeType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetAnimes(), nil
				},
			},
			"anime": &graphql.Field{
				Type: animeType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}
					return GetAnimeByID(id), nil
				},
			},
			"episodios": &graphql.Field{
				Type: graphql.NewList(episodioType),
				Args: graphql.FieldConfigArgument{
					"animeId": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if animeID, ok := p.Args["animeId"]; ok {
						return GetEpisodiosByAnimeID(animeID.(int)), nil
					}
					return GetEpisodios(), nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
}
package graphql

import (
	"github.com/graphql-go/graphql"
)

func NewSchema() (*graphql.Schema, error) {
	// Definir tipos
	episodioType := graphql.NewObject(
		graphql.ObjectConfig{
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
		},
	)

	animeType := graphql.NewObject(
		graphql.ObjectConfig{
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
					Type: graphql.NewList(graphql.String),
				},
				"episodios": &graphql.Field{
					Type: graphql.NewList(episodioType),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						anime := p.Source.(map[string]interface{})
						return GetEpisodiosByAnimeID(int(anime["id"].(float64))), nil
					},
				},
			},
		},
	)

	// Definir queries
	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
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
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id := p.Args["id"].(int)
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
						animeID, ok := p.Args["animeId"].(int)
						if ok {
							return GetEpisodiosByAnimeID(animeID), nil
						}
						return GetEpisodios(), nil
					},
				},
			},
		},
	)

	// Criar schema
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery,
		},
	)

	if err != nil {
		return nil, err
	}

	return &schema, nil
}
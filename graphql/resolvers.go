package graphql

import (
	"encoding/json"
	"os"
)

type Anime struct {
	ID        int      `json:"id"`
	Titulo    string   `json:"titulo"`
	Descricao string   `json:"descricao"`
	Ano       int      `json:"ano"`
	Generos   []string `json:"generos"`
}

type Episodio struct {
	ID      int    `json:"id"`
	AnimeID int    `json:"animeId"`
	Titulo  string `json:"titulo"`
	Numero  int    `json:"numero"`
	Duracao int    `json:"duracao"`
}

func GetAnimes() []Anime {
	file, err := os.ReadFile("data/animes.json")
	if err != nil {
		return []Anime{}
	}

	var animes []Anime
	err = json.Unmarshal(file, &animes)
	if err != nil {
		return []Anime{}
	}

	return animes
}

func GetAnimeByID(id int) *Anime {
	animes := GetAnimes()
	for _, anime := range animes {
		if anime.ID == id {
			return &anime
		}
	}
	return nil
}

func GetEpisodios() []Episodio {
	file, err := os.ReadFile("data/episodios.json")
	if err != nil {
		return []Episodio{}
	}

	var episodios []Episodio
	err = json.Unmarshal(file, &episodios)
	if err != nil {
		return []Episodio{}
	}

	return episodios
}

func GetEpisodiosByAnimeID(animeID int) []Episodio {
	episodios := GetEpisodios()
	var result []Episodio
	for _, episodio := range episodios {
		if episodio.AnimeID == animeID {
			result = append(result, episodio)
		}
	}
	return result
}
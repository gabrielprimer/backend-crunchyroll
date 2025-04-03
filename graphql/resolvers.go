package graphql

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
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

var (
	animesCache    []Anime
	episodiosCache []Episodio
	once           sync.Once
)

func loadData() {
	once.Do(func() {
		// Obter o caminho absoluto para os arquivos JSON
		animesPath := filepath.Join("data", "animes.json")
		episodiosPath := filepath.Join("data", "episodios.json")

		// Carrega animes
		file, err := os.ReadFile(animesPath)
		if err != nil {
			animesCache = []Anime{}
		} else {
			json.Unmarshal(file, &animesCache)
		}

		// Carrega episódios
		file, err = os.ReadFile(episodiosPath)
		if err != nil {
			episodiosCache = []Episodio{}
		} else {
			json.Unmarshal(file, &episodiosCache)
		}
	})
}

func GetAnimes() []Anime {
	loadData()
	return animesCache
}

func GetAnimeByID(id int) *Anime {
	loadData()
	for _, anime := range animesCache {
		if anime.ID == id {
			return &anime
		}
	}
	return nil
}

func GetEpisodios() []Episodio {
	loadData()
	return episodiosCache
}

func GetEpisodiosByAnimeID(animeID int) []Episodio {
	loadData()
	var result []Episodio
	for _, episodio := range episodiosCache {
		if episodio.AnimeID == animeID {
			result = append(result, episodio)
		}
	}
	return result
}
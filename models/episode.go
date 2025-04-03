package models

type Episode struct {
    ID          string `json:"id"`
    AnimeID     string `json:"anime_id"`
    Season      int    `json:"season"`
    Title       string `json:"title"`
    Slug        string `json:"slug"`
    Duration    string `json:"duration"`
    Synopsis    string `json:"synopsis"`
    Image       string `json:"image"`
    VideoURL    string `json:"video_url"`
    ReleaseDate string `json:"release_date,omitempty"` // Alterado para string
    IsLancamento bool  `json:"is_lancamento"`
    CreatedAt   string `json:"created_at"` // Alterado para string
    UpdatedAt   string `json:"updated_at"` // Alterado para string
}
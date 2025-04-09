package enums

// AiringDay representa os dias da semana para lançamento de episódios
type AiringDay string

const (
	Monday    AiringDay = "monday"
	Tuesday   AiringDay = "tuesday"
	Wednesday AiringDay = "wednesday"
	Thursday  AiringDay = "thursday"
	Friday    AiringDay = "friday"
	Saturday  AiringDay = "saturday"
	Sunday    AiringDay = "sunday"
)

// AnimeStatus representa o status de um anime
type AnimeStatus string

const (
	Ongoing    AnimeStatus = "ongoing"
	Completed  AnimeStatus = "completed"
	Announced  AnimeStatus = "announced"
	Cancelled  AnimeStatus = "cancelled"
)

// AudioLanguage representa os idiomas de áudio disponíveis
type AudioLanguage string

const (
	Portuguese AudioLanguage = "portuguese"
	Japanese   AudioLanguage = "japanese"
	Chinese    AudioLanguage = "chinese"
	Korean     AudioLanguage = "korean"
)

// AudioType representa o tipo de áudio (dublado/legendado)
type AudioType string

const (
	Sub  AudioType = "sub"
	Dub  AudioType = "dub"
	Both AudioType = "both"
)

// EpisodeLanguageType representa o tipo de linguagem do episódio
type EpisodeLanguageType string

const (
	Subtitled EpisodeLanguageType = "SUBTITLED"
	Dubbed    EpisodeLanguageType = "DUBBED"
	Raw       EpisodeLanguageType = "RAW"
)

// SeasonEnum representa as estações do ano
type SeasonEnum string

const (
	Winter SeasonEnum = "winter"
	Spring SeasonEnum = "spring"
	Summer SeasonEnum = "summer"
	Fall   SeasonEnum = "fall"
)

// SourceType representa o tipo de fonte original do conteúdo
type SourceType string

const (
	Manga       SourceType = "manga"
	LightNovel  SourceType = "light_novel"
	VisualNovel SourceType = "visual_novel"
	WebComic    SourceType = "web_comic"
	Original    SourceType = "original"
)

// SubtitleLanguage representa os idiomas de legenda disponíveis
type SubtitleLanguage string

const (
	SubPortuguese SubtitleLanguage = "portuguese"
)

// Genre representa os gêneros de anime disponíveis
type Genre string

const (
	Action       Genre = "Action"
	Adventure    Genre = "Adventure"
	Comedy       Genre = "Comedy"
	Drama        Genre = "Drama"
	Fantasy      Genre = "Fantasy"
	Music        Genre = "Music"
	Romance      Genre = "Romance"
	SciFi        Genre = "Sci-Fi"
	Seinen       Genre = "Seinen"
	Shojo        Genre = "Shojo"
	Shonen       Genre = "Shonen"
	SliceOfLife  Genre = "Slice of life"
	Sports       Genre = "Sports"
	Supernatural Genre = "Supernatural"
	Thriller     Genre = "Thriller"
)

// Métodos auxiliares para conversão e validação

func (ad AiringDay) IsValid() bool {
	switch ad {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	}
	return false
}

func (as AnimeStatus) IsValid() bool {
	switch as {
	case Ongoing, Completed, Announced, Cancelled:
		return true
	}
	return false
}

func (al AudioLanguage) IsValid() bool {
	switch al {
	case Portuguese, Japanese, Chinese, Korean:
		return true
	}
	return false
}

func (at AudioType) IsValid() bool {
	switch at {
	case Sub, Dub, Both:
		return true
	}
	return false
}

func (elt EpisodeLanguageType) IsValid() bool {
	switch elt {
	case Subtitled, Dubbed, Raw:
		return true
	}
	return false
}

func (se SeasonEnum) IsValid() bool {
	switch se {
	case Winter, Spring, Summer, Fall:
		return true
	}
	return false
}

func (st SourceType) IsValid() bool {
	switch st {
	case Manga, LightNovel, VisualNovel, WebComic, Original:
		return true
	}
	return false
}

func (sl SubtitleLanguage) IsValid() bool {
	return sl == SubPortuguese
}

func (g Genre) IsValid() bool {
	switch g {
	case Action, Adventure, Comedy, Drama, Fantasy, Music, Romance,
		SciFi, Seinen, Shojo, Shonen, SliceOfLife, Sports, Supernatural, Thriller:
		return true
	}
	return false
}
package enums

type EpisodeLanguageType string

const (
	Subtitled EpisodeLanguageType = "SUBTITLED"
	Dubbed    EpisodeLanguageType = "DUBBED"
	Raw       EpisodeLanguageType = "RAW"
)

func (elt EpisodeLanguageType) IsValid() bool {
	switch elt {
	case Subtitled, Dubbed, Raw:
		return true
	}
	return false
}
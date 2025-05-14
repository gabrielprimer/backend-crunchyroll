package enums

type SubtitleLanguage string

const (
	SubPortuguese SubtitleLanguage = "portuguese"
)

func (sl SubtitleLanguage) IsValid() bool {
	return sl == SubPortuguese
}
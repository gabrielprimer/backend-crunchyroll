package enums

type AudioLanguage string

const (
	Portuguese AudioLanguage = "portuguese"
	Japanese   AudioLanguage = "japanese"
	Chinese    AudioLanguage = "chinese"
	Korean     AudioLanguage = "korean"
)

type AudioType string

const (
	Sub  AudioType = "sub"
	Dub  AudioType = "dub"
	Both AudioType = "both"
)

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
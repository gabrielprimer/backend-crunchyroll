package enums

type SourceType string

const (
	Manga       SourceType = "manga"
	LightNovel  SourceType = "light_novel"
	VisualNovel SourceType = "visual_novel"
	WebComic    SourceType = "web_comic"
	Original    SourceType = "original"
)

func (st SourceType) IsValid() bool {
	switch st {
	case Manga, LightNovel, VisualNovel, WebComic, Original:
		return true
	}
	return false
}
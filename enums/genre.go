package enums

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

func (g Genre) IsValid() bool {
	switch g {
	case Action, Adventure, Comedy, Drama, Fantasy, Music, Romance,
		SciFi, Seinen, Shojo, Shonen, SliceOfLife, Sports, Supernatural, Thriller:
		return true
	}
	return false
}
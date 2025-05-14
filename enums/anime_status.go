package enums

type AnimeStatus string

const (
	Ongoing    AnimeStatus = "ongoing"
	Completed  AnimeStatus = "completed"
	Announced  AnimeStatus = "announced"
	Cancelled  AnimeStatus = "cancelled"
)

func (as AnimeStatus) IsValid() bool {
	switch as {
	case Ongoing, Completed, Announced, Cancelled:
		return true
	}
	return false
}
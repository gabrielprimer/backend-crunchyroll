package enums

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

func (ad AiringDay) IsValid() bool {
	switch ad {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	}
	return false
}
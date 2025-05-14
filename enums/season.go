package enums

type SeasonEnum string

const (
	Winter SeasonEnum = "winter"
	Spring SeasonEnum = "spring"
	Summer SeasonEnum = "summer"
	Fall   SeasonEnum = "fall"
)

func (se SeasonEnum) IsValid() bool {
	switch se {
	case Winter, Spring, Summer, Fall:
		return true
	}
	return false
}
package gowork

const (
	STALE_ENTITY_MSG = "Stale Entity. It has been updated in another session! Please reload and try again."
	CANNOT_BE_ZERO_VALUE = "value cannot be 'zero' value"
)

type NotFoundError string

func (s NotFoundError) Error() string {
	str := string(s)
	if str == "" {
		return "not found"
	}

	return str
}

func NewNotFoundError() NotFoundError {
	return NotFoundError("")
}

type StaleEntityError string

func (s StaleEntityError) Error() string {
	str := string(s)
	if str == "" {
		return STALE_ENTITY_MSG
	}

	return str
}

func NewStaleEntityError() NotFoundError {
	return NotFoundError(STALE_ENTITY_MSG)
}

package errors

const (
	ErrCannotCreateObservation     = "observation.errCannotCreate"
	ErrObservationNameMustBeUnique = "observation.errNameMustBeUnique"
	ErrObservationURLMustBeUnique  = "observation.errURLMustBeUnique"
	ErrCannotFetchObservations     = "observation.errCannotFetch"
	ErrObservationNotFound         = "observation.errNotFound"
	ErrCannotUpdateObservation     = "observation.errCannotUpdate"
	ErrCannotDeleteObservations    = "observation.errCannotDelete"
	ErrInvalidObservationName      = "observation.errInvalidName"
	ErrInvalidObservationURL       = "observation.errInvalidURL"
	ErrInvalidExcludedFor          = "excluded.errInvalidFor"
	ErrInvalidExcludedValue        = "excluded.errInvalidValue"
	ErrInvalidOneOfFor             = "oneOf.errInvalidFor"
	ErrInvalidOneOfValue           = "oneOf.errInvalidValue"
)

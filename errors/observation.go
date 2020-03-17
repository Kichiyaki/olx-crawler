package errors

const (
	ErrCannotCreateObservation     = "observation.errCannotCreate"
	ErrObservationNameMustBeUnique = "observation.errNameMustBeUnique"
	ErrObservationURLMustBeUnique  = "observation.errURLMustBeUnique"
	ErrCannotFetchObservations     = "observation.errCannotFetch"
	ErrObservationNotFound         = "observation.errObservationNotFound"
	ErrCannotUpdateObservation     = "observation.errCannotUpdate"
	ErrCannotDeleteObservations    = "observation.errCannotDelete"
	ErrInvalidObservationName      = "observation.errInvalidName"
	ErrInvalidObservationURL       = "observation.errInvalidURL"
	ErrInvalidExcludeFor           = "exclude.errInvalidFor"
	ErrInvalidExcludeValue         = "exclude.errInvalidValue"
	ErrInvalidOneOfFor             = "oneOf.errInvalidFor"
	ErrInvalidOneOfValue           = "oneOf.errInvalidValue"
)

package models

type Error struct {
	Message string  `json:"message"`
	Details []Error `json:"details,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

package models

type Error struct {
	Message  string  `json:"message"`
	Detailed []Error `json:"detailed,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

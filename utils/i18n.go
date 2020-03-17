package utils

import "github.com/nicksnyder/go-i18n/v2/i18n"

func NewI18NMsg(id, msg string) *i18n.Message {
	return &i18n.Message{
		ID:    id,
		One:   msg,
		Other: msg,
	}
}

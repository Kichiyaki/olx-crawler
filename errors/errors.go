package errors

import (
	"olx-crawler/models"
)

func Wrap(msg string, details []error) models.Error {
	return models.Error{
		Message: msg,
		Details: convertToErrorStruct(details...),
	}
}

func ToErrorModel(err error) models.Error {
	return convertToErrorStruct(err)[0]
}

func convertToErrorStruct(errors ...error) []models.Error {
	es := []models.Error{}
	for _, err := range errors {
		e, ok := err.(models.Error)
		if !ok {
			e = models.Error{
				Message: err.Error(),
			}
		}
		es = append(es, e)
	}
	return es
}

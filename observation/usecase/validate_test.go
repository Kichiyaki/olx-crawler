package usecase

import (
	"olx-crawler/errors"
	"olx-crawler/models"
	"testing"
)

func TestValidate(t *testing.T) {
	cfg := newConfig()
	t.Run("name is required", func(t *testing.T) {
		err := cfg.validate(models.Observation{})
		if err.Error() != errors.ErrInvalidObservationName {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidObservationName)
		}
	})

	t.Run("url is required", func(t *testing.T) {
		err := cfg.validate(models.Observation{Name: "asdasda"})
		if err.Error() != errors.ErrInvalidObservationURL {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidObservationURL)
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		err := cfg.validate(models.Observation{Name: "asdasda", URL: "sdoslx.pl"})
		if err.Error() != errors.ErrInvalidObservationURL {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidObservationURL)
		}
	})

	t.Run("keyword type is invalid", func(t *testing.T) {
		err := cfg.validate(models.Observation{Name: "asdasda",
			URL:      "https://www.olx.pl",
			Keywords: []models.Keyword{models.Keyword{Type: "eloszka"}}})
		if err.Error() != errors.ErrInvalidKeywordType {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidKeywordType)
		}
	})

	t.Run("keyword group is invalid", func(t *testing.T) {
		err := cfg.validate(models.Observation{Name: "asdasda",
			URL:      "https://www.olx.pl",
			Keywords: []models.Keyword{models.Keyword{Type: "one_of"}}})
		if err.Error() != errors.ErrInvalidKeywordGroup {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidKeywordGroup)
		}
	})

	t.Run("keyword for is invalid", func(t *testing.T) {
		cfg := newConfig()
		err := cfg.validate(models.Observation{Name: "asdasda",
			URL:      "https://www.olx.pl",
			Keywords: []models.Keyword{models.Keyword{Type: "excluded"}}})
		if err.Error() != errors.ErrInvalidKeywordFor {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidKeywordFor)
		}
	})

	t.Run("keyword value is required", func(t *testing.T) {
		err := cfg.validate(models.Observation{Name: "asdasda",
			URL:      "https://www.olx.pl",
			Keywords: []models.Keyword{models.Keyword{Type: "excluded", For: "title"}}})
		if err.Error() != errors.ErrInvalidKeywordValue {
			t.Errorf("Got %v, expected %v", err, errors.ErrInvalidKeywordValue)
		}
	})

	t.Run("success", func(t *testing.T) {
		err := cfg.validate(models.Observation{Name: "asdasda",
			URL:      "https://www.olx.pl",
			Keywords: []models.Keyword{models.Keyword{Type: "excluded", For: "title", Value: "nil"}}})
		if err != nil {
			t.Errorf("Got %v, expected %v", err, nil)
		}
	})
}

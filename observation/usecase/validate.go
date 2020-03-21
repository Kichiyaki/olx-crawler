package usecase

import (
	"olx-crawler/errors"
	"olx-crawler/models"
	"strings"
)

type config struct {
	Name     bool
	URL      bool
	Keywords bool
}

func newConfig() config {
	return config{
		true, true, true,
	}
}

func (cfg config) validate(o models.Observation) error {
	if cfg.Name && o.Name == "" {
		return errors.Wrap(errors.ErrInvalidObservationName, []error{})
	} else if cfg.URL && (o.URL == "" || !strings.Contains(o.URL, "olx.pl")) {
		return errors.Wrap(errors.ErrInvalidObservationURL, []error{})
	}

	if cfg.Keywords {
		for _, keyword := range o.Keywords {
			if keyword.Type != "one_of" && keyword.Type != "excluded" && keyword.Type != "required" {
				return errors.Wrap(errors.ErrInvalidKeywordType, []error{})
			} else if keyword.Type == "one_of" && keyword.Group == "" {
				return errors.Wrap(errors.ErrInvalidKeywordGroup, []error{})
			} else if keyword.For != "title" && keyword.For != "description" {
				return errors.Wrap(errors.ErrInvalidKeywordFor, []error{})
			} else if keyword.Value == "" {
				return errors.Wrap(errors.ErrInvalidKeywordValue, []error{})
			}
		}
	}

	return nil
}

package usecase

import (
	"olx-crawler/errors"
	"olx-crawler/models"
	"strings"
)

type config struct {
	Name    bool
	URL     bool
	Exclude bool
	OneOf   bool
}

func newConfig() config {
	return config{
		true, true, true, true,
	}
}

func (cfg config) validate(o models.Observation) error {
	if cfg.Name && o.Name == "" {
		return errors.Wrap(errors.ErrInvalidObservationName, []error{})
	} else if cfg.URL && (o.URL == "" || !strings.Contains(o.URL, "olx.pl")) {
		return errors.Wrap(errors.ErrInvalidObservationURL, []error{})
	}

	if cfg.Exclude {
		for _, excl := range o.Exclude {
			if excl.For != "title" && excl.For != "description" {
				return errors.Wrap(errors.ErrInvalidExcludeFor, []error{})
			} else if excl.Value == "" {
				return errors.Wrap(errors.ErrInvalidExcludeValue, []error{})
			}
		}
	}

	if cfg.OneOf {
		for _, oneOf := range o.OneOf {
			if oneOf.For != "title" && oneOf.For != "description" {
				return errors.Wrap(errors.ErrInvalidOneOfFor, []error{})
			} else if oneOf.Value == "" {
				return errors.Wrap(errors.ErrInvalidOneOfValue, []error{})
			}
		}
	}

	return nil
}

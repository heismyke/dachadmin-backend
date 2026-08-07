package application

import (
	"dach-admin/internal/domain"
	"net/mail"
	"strings"
	"time"
)

func Required(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return domain.ValidationError{Field: field, Message: "is required"}
	}
	return nil
}

func Email(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return domain.ValidationError{Field: field, Message: "is required"}
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return domain.ValidationError{Field: field, Message: "must be a valid email"}
	}
	return nil
}

func NonNegative(field string, value *float64) error {
	if value != nil && *value < 0 {
		return domain.ValidationError{Field: field, Message: "must be non-negative"}
	}
	return nil
}

func Positive(field string, value float64) error {
	if value <= 0 {
		return domain.ValidationError{Field: field, Message: "must be positive"}
	}
	return nil
}

func DateRange(from *time.Time, to *time.Time) error {
	if from != nil && to != nil && to.Before(*from) {
		return domain.ValidationError{Field: "valid_to", Message: "must be after valid_from"}
	}
	return nil
}

func In[T ~string](field string, value T, allowed ...T) error {
	for _, item := range allowed {
		if value == item {
			return nil
		}
	}
	return domain.ValidationError{Field: field, Message: "has invalid value"}
}

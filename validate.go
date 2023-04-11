package publication

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	tagJson = "json"
)

const (
	validatorTagTitle = "title"
	validatorTagText  = "text"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
var ErrInvalidTitle = errors.New("wrong title")
var ErrInvalidText = errors.New("wrong text")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s", v.Err)
}

func (v ValidationErrors) Error() string {
	var s string
	for _, err := range v {
		s += fmt.Sprintf("%s", err.Err)
	}
	return s
}

func (v *ValidationErrors) Add(err error) {
	if err == nil {
		return
	}
	*v = append(*v, ValidationError{
		Err: err,
	})

}

func Validate(v any) error {
	var Errors ValidationErrors
	value := reflect.ValueOf(v)
	if value.Type().Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < value.Type().NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)
		validator, ok := field.Tag.Lookup(tagJson)
		if !ok {
			continue
		}
		if !field.IsExported() {
			Errors.Add(ErrValidateForUnexportedFields)
			continue
		}
		switch validator {
		case validatorTagTitle:
			Errors.Add(validateTitle(fieldValue.String()))
		case validatorTagText:
			Errors.Add(validateText(fieldValue.String()))
		}
	}
	if len(Errors) == 0 {
		return nil
	}
	return Errors
}

func validateTitle(title string) error {
	if title == "" {
		return ErrInvalidTitle
	}
	if len(title) >= 100 {
		return ErrInvalidTitle
	}
	return nil
}

func validateText(text string) error {
	if text == "" {
		return ErrInvalidText
	}
	if len(text) >= 500 {
		return ErrInvalidText
	}
	return nil
}

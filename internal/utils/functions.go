package utils

import (
	"errors"
	"net/http"
	"reflect"
)

func HasAtLeastOneField(v interface{}) bool {
	val := reflect.ValueOf(v)

	// Se for ponteiro, pega o valor
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Itera sobre os campos da struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Verifica se o campo é um ponteiro e não é nil
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			return true
		}

		// Para campos não-ponteiro, verifica se não é zero value
		if field.Kind() != reflect.Ptr && !field.IsZero() {
			return true
		}
	}

	return false
}

func GetApiErr(err error) *HTTPError {
	var apiErr *HTTPError

	if errors.As(err, &apiErr) {
		return apiErr
	} else {
		return NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}

}

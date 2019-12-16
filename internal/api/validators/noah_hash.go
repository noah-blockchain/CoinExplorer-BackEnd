package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func NoahTxHash(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return isValidNoahHash(field.String())
}

func isValidNoahHash(hash string) bool {
	return regexp.MustCompile("^Nt([A-Fa-f0-9]{64})$").MatchString(hash)
}

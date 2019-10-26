package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func NoahPublicKey(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return isValidNoahPublicKey(field.String())
}

func isValidNoahPublicKey(publicKey string) bool {
	return regexp.MustCompile("^Np([A-Fa-f0-9]{64})$").MatchString(publicKey)
}

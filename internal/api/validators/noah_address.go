package validators

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"regexp"
)

func NoahAddress(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if fieldType.String() == "[]string" {
		data, _ := field.Interface().([]string)
		for _, address := range data {
			if !isValidNoahAddress(address) {
				return false
			}
		}

		return true
	}

	return isValidNoahAddress(field.String())
}

func isValidNoahAddress(address string) bool {
	return regexp.MustCompile("^NOAHx([A-Fa-f0-9]{40})$").MatchString(address)
}

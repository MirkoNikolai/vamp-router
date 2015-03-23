package haproxy

import (
	valid "github.com/asaskevich/govalidator"
	"regexp"
)

/* Initialize a set of validator used throughout the Haproxy package

Validators are used on the side of the API but also internally to check validity
of generated artifacts like socket paths.
*/
func init() {

	// validation for route names. Should be ascii, alphanumeric but allowing - _ . (dash, underscore, period)
	valid.TagMap["routeName"] = valid.Validator(func(str string) bool {

		pattern := "^[a-zA-Z0-9]{1}[a-zA-Z0-9.\\-_]{3,63}$"
		routeName := regexp.MustCompile(pattern)
		return routeName.MatchString(str)
	})

	// validation for full sockets paths. These cannot be longer than 103 characters.
	valid.TagMap["socketPath"] = valid.Validator(func(str string) bool {

		pattern := "^[a-zA-Z0-9/]{1}[a-zA-Z0-9.\\-_/]{1,102}$"
		socketPath := regexp.MustCompile(pattern)
		return socketPath.MatchString(str)
	})

}

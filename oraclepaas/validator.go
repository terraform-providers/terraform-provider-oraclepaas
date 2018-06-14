package oraclepaas

import (
	"fmt"
	"regexp"
)

func validateAccessRuleName(v interface{}, k string) (ws []string, errors []error) {

	value := v.(string)

	if len(value) > 50 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q name can only be between 1-50 characters. Got: %s", k, value))
	}

	re := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]+$")

	if !re.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must start with a letter and contain only letters, numbers or underscore (_). Got: %s", k, value))
	}
	return
}

/**
  Validates the service name used for the mysql instance.
  Rules of the validation :
  - Must not exceed 50 characters,
  - Must start with a letter
  - Must contain only letters, numbers, or hyphens.
  - Must not end with a hyphen.
  Must not contain any other special characters.
*/
func validateMySQLServiceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 50 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q name can only be between 1-50 characters. Got: %s", k, value))
	}

	re := regexp.MustCompile("^[a-zA-Z]([a-zA-Z0-9-]*[A-Za-z0-9])+$")
	if !re.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must start with a letter and contain only letters, numbers or hyphen (-) and cannot end with a hyphen. Got: %s", k, value))
	}
	return
}

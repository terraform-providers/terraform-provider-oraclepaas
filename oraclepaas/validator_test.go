package oraclepaas

import (
	"testing"
)

func TestValidateAccessRuleName(t *testing.T) {
	validNames := []string{
		"SampleNAme",
		"SampleName1",
		"Sample_Name_1",
		"sample_Name1",
	}

	for _, v := range validNames {
		_, errors := validateAccessRuleName(v, "rule_name")
		if len(errors) != 0 {
			t.Fatalf("%q rule name sshould pass: %q", v, errors)
		}
	}

	invalidNames := []string{
		"1nvalidName_startswithnumber",
		"InvalidName!_hasIllegalCharacters",
		"Invalid-Name_Has Space and -hyphen",
		"SuperLongNameThatExceedsMaxLength_123456789012345678901234567890",
	}

	for _, v := range invalidNames {
		_, errors := validateAccessRuleName(v, "rule_name")
		if len(errors) == 0 {
			t.Fatalf("%q rule name should fail: %q", v, errors)
		}
	}
}

func TestValidateMySQLServiceName(t *testing.T) {
	validNames := []string{
		"SampleNAme",
		"SampleName1",
		"Sample-Name-1",
		"sample-Name1",
	}

	for _, v := range validNames {
		_, errors := validateMySQLServiceName(v, "rule_name")
		if len(errors) != 0 {
			t.Fatalf("%q rule name should pass: %q", v, errors)
		}
	}

	invalidNames := []string{
		"1nvalidName_startswithnumber",
		"InvalidName!_hasIllegalCharacters",
		"Invalid-Name-Has space and -underscore",
		"Invalid-Name-ends-with-hyphen-",
		"SuperLongNameThatExceedsMaxLength_123456789012345678901234567890",
	}

	for _, v := range invalidNames {
		_, errors := validateMySQLServiceName(v, "rule_name")
		if len(errors) == 0 {
			t.Fatalf("%q rule name should fail: %q", v, errors)
		}
	}
}

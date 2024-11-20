package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0)
	for _, e := range v {
		errs = append(errs, fmt.Sprintf("%s ", e.Err))
	}

	return strings.Join(errs, "\n")
}

func Validate(v interface{}) error {
	t := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	if t.Kind() != reflect.Struct {
		return errors.New("program error: Input data should be a struct")
	}

	var ve ValidationErrors
	var err error

	for i := 0; i < t.NumField(); i++ {
		if vv.Field(i).CanInterface() && t.Field(i).Tag.Get("validate") != "" {
			for _, rule := range strings.Split(t.Field(i).Tag.Get("validate"), "|") {
				ve, err = applyOneRule(rule, t.Field(i).Name, vv.Field(i), ve)
				if err != nil {
					return err
				}
			}
		}
	}
	return ve
}

func applyOneRule(rule string, nameField string, dataForValidate reflect.Value,
	ves ValidationErrors,
) (ValidationErrors, error) {
	rl := strings.SplitN(rule, ":", 2)
	if len(rl) != 2 {
		return nil, fmt.Errorf("program error: Field %s -> invalid rule", nameField)
	}

	switch dataForValidate.Kind() {
	case reflect.String:
		validErr, err := validateString(rl[0], rl[1], nameField, dataForValidate.String())
		if !errors.Is(validErr.Err, nil) {
			ves = append(ves, validErr)
			return ves, nil
		}
		return ves, err

	case reflect.Int:
		validErr, err := validateInt(rl[0], rl[1], nameField, dataForValidate.Int())
		if !errors.Is(validErr.Err, nil) {
			ves = append(ves, validErr)
			return ves, nil
		}
		return ves, err

	case reflect.Slice:
		var err error
		for y := 0; y < dataForValidate.Len(); y++ {
			ves, err = applyOneRule(rule, nameField+" element "+strconv.Itoa(y), dataForValidate.Index(y), ves)
			if !errors.Is(err, nil) {
				return ves, err
			}
		}
		return ves, nil

	default:
		return nil, fmt.Errorf("program error: Field %s -> type is not supported", nameField)
	}
}

func validateString(ruleName string, ruleValue string, nameField string,
	dataStrForValidator string,
) (ValidationError, error) {
	switch ruleName {
	case "len":
		rv, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ValidationError{}, err
		}
		if len(dataStrForValidator) != rv {
			return ValidationError{
				nameField,
				fmt.Errorf("validation error: Field %s -> length must be %s", nameField, ruleValue),
			}, nil
		}
	case "regexp":
		re, err := regexp.Compile(ruleValue)
		if err != nil {
			return ValidationError{}, fmt.Errorf("program error: Field %s -> invalid regexp pattern", nameField)
		}
		if !re.MatchString(dataStrForValidator) {
			return ValidationError{
				nameField,
				fmt.Errorf("validation error: Field %s -> regexp does not match %s", nameField, dataStrForValidator),
			}, nil
		}
	case "in":
		values := strings.Split(ruleValue, ",")
		for _, v := range values {
			if dataStrForValidator == v {
				return ValidationError{nameField, nil}, nil
			}
		}
		return ValidationError{
			nameField,
			fmt.Errorf("validation error: Field %s -> must be one of %s", nameField, ruleValue),
		}, nil
	default:
		return ValidationError{}, fmt.Errorf("program error: Field %s -> unknown rule %s", nameField, ruleName)
	}
	return ValidationError{nameField, nil}, nil
}

func validateInt(ruleName string, ruleValue string, nameField string, dataIntForValidator int64,
) (ValidationError, error) {
	switch ruleName {
	case "min":
		rv, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return ValidationError{}, fmt.Errorf("program error: Field %s -> invalid parse int", nameField)
		}
		if dataIntForValidator < rv {
			return ValidationError{
				nameField,
				fmt.Errorf("validation error: Field %s -> must be >= %s", nameField, ruleValue),
			}, nil
		}
	case "max":
		rv, err := strconv.ParseInt(ruleValue, 10, 64)
		if err != nil {
			return ValidationError{}, fmt.Errorf("program error: Field %s -> invalid parse int", nameField)
		}
		if dataIntForValidator > rv {
			return ValidationError{
				nameField,
				fmt.Errorf("validation error: Field %s -> max is %s", nameField, ruleValue),
			}, nil
		}
	case "in":
		values := strings.Split(ruleValue, ",")
		for _, v := range values {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return ValidationError{}, fmt.Errorf("program error: Field %s -> invalid parse int", nameField)
			}
			if dataIntForValidator == vInt {
				return ValidationError{nameField, nil}, nil
			}
		}
		return ValidationError{
			nameField,
			fmt.Errorf("validation error: Field %s -> must be one of %s ", nameField, ruleValue),
		}, nil
	default:
		return ValidationError{}, fmt.Errorf("program error: Field %s -> unknown rule %s ", nameField, ruleName)
	}
	return ValidationError{nameField, nil}, nil
}

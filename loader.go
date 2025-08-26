package goenvloader

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrEmptyEnvVar = func(env string) error {
		return fmt.Errorf("required environment variable %s is empty", env)
	}
	ErrEmptyTag = func(env string) error {
		return fmt.Errorf("tag is empty for environment variable %s", env)
	}
	ErrProcessField = errors.New("failed to process field")
)

func Load(cfg interface{}) error {
	configStruct := reflect.ValueOf(cfg)
	// Check if the provided config is a pointer to a struct
	if configStruct.Kind() != reflect.Ptr || configStruct.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got %T", cfg)
	}

	// Get the actual struct to access it
	configStruct = configStruct.Elem()
	// Get type of struct to access its fields
	t := configStruct.Type()

	// Process each field in the struct
	for i := 0; i < configStruct.NumField(); i++ {
		field := configStruct.Field(i)
		fieldType := t.Field(i)

		err := processField(field, fieldType)
		if err != nil {
			return fmt.Errorf(err.Error(), ErrProcessField.Error())
		}
	}

	return nil
}

func processField(field reflect.Value, fieldType reflect.StructField) error {
	// We can add more tags for added options and requirements
	tag := fieldType.Tag.Get("env")
	requiredTag := fieldType.Tag.Get("required")
	defaultTag := fieldType.Tag.Get("default")
	if tag == "" {
		// By default, env tag is required to process a field
		// If no "env" tag, check and process nested struct, or else skip field with error
		if field.Kind() == reflect.Struct {
			if err := Load(field.Addr().Interface()); err != nil {
				return fmt.Errorf(err.Error(), "failed to load nested struct %s", fieldType.Name)
			}
		} else {
			// In this case, we skip the field since it's not a struct and doesn't have a tag
			return ErrEmptyTag(fieldType.Name)
		}
	} else {
		// Process env tag. This is where we can add more tag options and add requirements
		value := os.Getenv(tag)
		if value == "" {
			// If the required tag is "true" and the value is empty, return an error
			if requiredTag == "true" {
				return ErrEmptyEnvVar(tag)
			}
			// If the required tag is not "true" and value is empty, use default tag if available
			if defaultTag != "" {
				value = defaultTag
			} else {
				return nil // No value to set and no default value, so skip
			}
		}

		// Set the field value
		if err := setFieldValue(field, value); err != nil {
			return fmt.Errorf(err.Error(), "failed to set field %s", fieldType.Name)
		}

		// Validation based on field type. Not sure on this one maybe we can flesh it out according to our needs later on
		if err := validateField(field); err != nil {
			return fmt.Errorf(err.Error(), "validation failed for field %s", fieldType.Name)
		}
	}

	return nil
}

// setFieldValue sets a value to a struct field.
func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return errors.New("cannot set field")
	}

	switch field.Kind() {
	case reflect.Int:
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(v))
	case reflect.String:
		field.SetString(value)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// validateField adds field validation as required.
func validateField(field reflect.Value) error {
	switch field.Kind() {
	case reflect.Int:
		if field.Int() <= 0 {
			return errors.New("integer must be greater than 0")
		}
	case reflect.String:
		if field.String() == "" {
			return errors.New("string cannot be empty")
		}
	default:
		return nil
	}
	return nil
}

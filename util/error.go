package util

import (
	"fmt"
)

func ErrorGenerate(tag string, err error) error {
	errors := fmt.Errorf("%s -> %w", tag, err)
	return errors
}

func ErrorHandler(tag string, err error) {
	errors := fmt.Errorf("%s -> %w", tag, err)
	fmt.Println(errors.Error())
}

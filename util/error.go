package util

import (
	"fmt"
)

func ErrorGenerate(tag string, err interface{}) error {
	errors := fmt.Errorf("%s -> %v", tag, err)
	return errors
}

func ErrorHandler(tag string, err interface{}) {
	errors := fmt.Errorf("%s -> %v", tag, err)
	fmt.Println(errors.Error())
}

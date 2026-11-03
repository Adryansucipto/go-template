package validator

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

type Validatable interface {
	Validate() error
}

type CustomValidator struct{}

func New() *CustomValidator {
	return &CustomValidator{}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if v, ok := i.(Validatable); ok {
		return v.Validate()
	}
	return nil
}

func RequireNonEmpty(fields map[string]string) error {
	var missing []string
	for name, val := range fields {
		if strings.TrimSpace(val) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return echo.NewHTTPError(400, fmt.Sprintf("missing required fields: %s", strings.Join(missing, ", ")))
	}
	return nil
}


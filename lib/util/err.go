package util

import (
	"errors"
	"fmt"
	errs "zinx/lib/enum/err"
)

// NewErrorWithPattern creates a new error with pattern
func NewErrorWithPattern(format errs.Pattern, a ...any) error {
	return errors.New(fmt.Sprintf(string(format), a))
}

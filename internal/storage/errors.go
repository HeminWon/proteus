package store

import (
	"errors"
	"fmt"
)

var ErrInvalidSettingsRoot = errors.New("invalid settings root: expected JSON object")

func WrapSettingsParseError(err error) error {
	if errors.Is(err, ErrInvalidSettingsRoot) {
		return fmt.Errorf("invalid settings root in %s: expected JSON object", SettingsPath())
	}
	return fmt.Errorf("failed to parse %s: %v", SettingsPath(), err)
}

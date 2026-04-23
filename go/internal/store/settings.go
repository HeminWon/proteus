package store

import (
	"encoding/json"
	"os"
)

type JsonObject map[string]any

type SettingsReadResult struct {
	Exists bool
	Data   JsonObject
}

func ReadSettings() (SettingsReadResult, error) {
	path := SettingsPath()
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return SettingsReadResult{Exists: false, Data: JsonObject{}}, nil
		}
		return SettingsReadResult{}, err
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return SettingsReadResult{}, err
	}

	obj, ok := parsed.(map[string]any)
	if !ok {
		return SettingsReadResult{}, ErrInvalidSettingsRoot
	}

	return SettingsReadResult{Exists: true, Data: JsonObject(obj)}, nil
}

func WriteSettings(settings JsonObject) error {
	content, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	return writeFileAtomic(SettingsPath(), content)
}

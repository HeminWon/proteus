package store

import (
	"encoding/json"
	"os"
)

type CacheData struct {
	Active *struct {
		Claude string `json:"claude,omitempty"`
	} `json:"active,omitempty"`
}

func ReadCache() CacheData {
	raw, err := os.ReadFile(CachePath())
	if err != nil {
		return CacheData{}
	}

	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return CacheData{}
	}

	obj, ok := parsed.(map[string]any)
	if !ok {
		return CacheData{}
	}

	activeRaw, ok := obj["active"].(map[string]any)
	if !ok {
		return CacheData{}
	}

	claude, ok := activeRaw["claude"].(string)
	if !ok {
		return CacheData{}
	}

	return CacheData{Active: &struct {
		Claude string `json:"claude,omitempty"`
	}{Claude: claude}}
}

func WriteCache(cache CacheData) error {
	content, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	return writeFileAtomic(CachePath(), content)
}

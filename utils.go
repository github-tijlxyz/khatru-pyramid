package main

import (
	"encoding/json"
)

func getProfileInfoFromJson(jsonStr string) (string, string) {
	fieldOrder := []string{"displayName", "display_name", "username", "name"}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		log.Error().Err(err).Msg("failed to read profile from json")
		return "", ""
	}

	var displayname string = "..."
	var picture string = ""

	for _, fieldName := range fieldOrder {
		if val, ok := data[fieldName]; ok {
			if strVal, ok := val.(string); ok && strVal != "" {
				if fieldName == "picture" {
					picture = strVal
				}
				if fieldName == "name" {
					displayname = strVal
				} else if displayname == "" {
					displayname = strVal
				}
			}
		}
	}

	return displayname, picture
}

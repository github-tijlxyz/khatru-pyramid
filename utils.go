package main

import (
	"context"
	"encoding/json"

	"github.com/nbd-wtf/go-nostr"
)

func isPublicKeyInWhitelist(pubkey string) bool {
	if pubkey == s.RelayPubkey {
		return true
	}

	for i := 0; i < len(whitelist); i++ {
		if whitelist[i].PublicKey == pubkey {
			return true
		}
	}
	return false
}

func deleteFromWhitelistRecursively(ctx context.Context, target string) {
	var updatedWhitelist []WhitelistEntry
	var queue []string

	// Remove from whitelist
	for _, user := range whitelist {
		if user.PublicKey != target {
			updatedWhitelist = append(updatedWhitelist, user)
		}
		if user.InvitedBy == target {
			queue = append(queue, user.PublicKey)
		}
	}
	whitelist = updatedWhitelist

	// Remove all events
	filter := nostr.Filter{
		Authors: []string{target},
	}
	events, _ := db.QueryEvents(ctx, filter)
	for ev := range events {
		err := db.DeleteEvent(ctx, ev)
		if err != nil {
			log.Error().Err(err).Msg("failed to delete event")
		}
	}

	// Recursive
	for _, pk := range queue {
		deleteFromWhitelistRecursively(ctx, pk)
	}
}

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

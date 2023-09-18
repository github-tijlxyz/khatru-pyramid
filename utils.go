package main

import (
	"context"
	"log"

	"github.com/nbd-wtf/go-nostr"
)

func isPkInWhitelist(targetPk string) bool {
    for i := 0; i < len(whitelist); i++ {
        if whitelist[i].Pk == targetPk {
            return true
        }
    }
    return false
}

func deleteFromWhitelistRecursively (ctx context.Context, target string) {
	var updatedWhitelist []User
	var queue []string

	// Remove from whitelist
	for _, user := range whitelist {
		if user.Pk != target {
			updatedWhitelist = append(updatedWhitelist, user)
		}
		if user.InvitedBy == target {
			queue = append(queue, user.Pk);
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
			log.Println("error while deleting event", err)
		}
	}

	// Recursive
	for _, pk := range queue {
		deleteFromWhitelistRecursively(ctx, pk)
	}
}

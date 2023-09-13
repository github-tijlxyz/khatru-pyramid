package main

import (
	"context"

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

func deleteFromWhitelistRecursively (target string) {
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
	go func () {
		var filter nostr.Filter = nostr.Filter{
			Authors: []string{target},
		}
		events, _ := db.QueryEvents(context.TODO(), filter)
		for ev := range events {
			db.DeleteEvent(context.TODO(), ev)
		}
	}()

	// Recursive
	for _, pk := range queue {
		deleteFromWhitelistRecursively(pk)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr"
)

const WHITELIST_FILE = "users.json"

type Whitelist map[string]string // { [user_pubkey]: [invited_by] }

func addToWhitelist(pubkey string, inviter string) error {
	if nostr.IsValidPublicKeyHex(pubkey) && isPublicKeyInWhitelist(inviter) && !isPublicKeyInWhitelist(pubkey) {
		whitelist[pubkey] = inviter
	}
	return saveWhitelist()
}

func isPublicKeyInWhitelist(pubkey string) bool {
	_, ok := whitelist[pubkey]
	return ok
}

func isAncestorOf(pubkey string, target string) bool {
	ancestor := target // we must find out if we are an ancestor of the target, but we can delete ourselves
	for {
		if ancestor == pubkey {
			break
		}

		parent, ok := whitelist[ancestor]
		if !ok {
			// parent is not in whitelist, this means this is a top-level user and can
			// only be deleted by manually editing the users.json file
			return false
		}

		ancestor = parent
	}
	return true
}

func removeFromWhitelist(target string, deleter string) error {
	// check if this user is a descendant of the user who issued the delete command
	if !isAncestorOf(deleter, target) {
		return fmt.Errorf("insufficient permissions to delete this")
	}

	// if we got here that means we have permission to delete the target
	delete(whitelist, target)

	// delete all people who were invited by the target
	removeDescendantsFromWhitelist(target)

	return saveWhitelist()
}

func removeDescendantsFromWhitelist(ancestor string) {
	for pubkey, inviter := range whitelist {
		if inviter == ancestor {
			delete(whitelist, pubkey)
			removeDescendantsFromWhitelist(pubkey)
		}
	}
}

func loadWhitelist() error {
	b, err := os.ReadFile(WHITELIST_FILE)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &whitelist); err != nil {
		return err
	}

	return nil
}

func saveWhitelist() error {
	jsonBytes, err := json.Marshal(whitelist)
	if err != nil {
		return err
	}

	if err := os.WriteFile(WHITELIST_FILE, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

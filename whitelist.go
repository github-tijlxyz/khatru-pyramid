package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr"
)

type Whitelist map[string]string // { [user_pubkey]: [invited_by] }

func addToWhitelist(pubkey string, inviter string) error {
	if !isPublicKeyInWhitelist(inviter) {
		return fmt.Errorf("pubkey %s doesn't have permission to invite", inviter)
	}

	if !nostr.IsValidPublicKey(pubkey) {
		return fmt.Errorf("pubkey invalid: %s", pubkey)
	}

	if isPublicKeyInWhitelist(pubkey) {
		return fmt.Errorf("pubkey already in whitelist: %s", pubkey)
	}

	whitelist[pubkey] = inviter
	return saveWhitelist()
}

func isPublicKeyInWhitelist(pubkey string) bool {
	_, ok := whitelist[pubkey]
	return ok
}

func canInviteMore(pubkey string) bool {
	if pubkey == "" {
		return false
	}

	if pubkey == s.RelayPubkey {
		return true
	}

	if !isPublicKeyInWhitelist(pubkey) {
		return false
	}

	count := 0
	for _, inviter := range whitelist {
		if inviter == pubkey {
			count++
		}
		if count >= s.MaxInvitesPerPerson {
			return false
		}
	}
	return true
}

func isAncestorOf(ancestor string, target string) bool {
	parent, ok := whitelist[target]
	if !ok {
		// parent is not in whitelist, this means this is a top-level user and can
		// only be deleted by manually editing the users.json file
		return false
	}

	if parent == ancestor {
		// if the pubkey is the parent, that means it is an ancestor
		return true
	}

	// otherwise we climb one degree up and test with the parent of the target
	return isAncestorOf(ancestor, parent)
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
	b, err := os.ReadFile(s.UserdataPath)
	if err != nil {
		// If the whitelist file does not exist, with RELAY_PUBKEY
		if errors.Is(err, os.ErrNotExist) {
			whitelist[s.RelayPubkey] = ""
			jsonBytes, err := json.Marshal(&whitelist)
			if err != nil {
				return err
			}
			if err := os.WriteFile(s.UserdataPath, jsonBytes, 0644); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
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

	if err := os.WriteFile(s.UserdataPath, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

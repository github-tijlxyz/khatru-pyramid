package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"golang.org/x/exp/slices"
)

type WhitelistEntry struct {
	InvitedBy string `json:"invited_by"`
	PublicKey string `json:"pk"`
}

var whitelist []WhitelistEntry

func whitelistRejecter(ctx context.Context, evt *nostr.Event) (reject bool, msg string) {
	// check if user in whitelist
	if !isPublicKeyInWhitelist(evt.PubKey) {
		return true, "You are not invited to this relay"
	}

	/*
	 kind 20201
	 invited/whitelisted user invites new user
	*/
	if evt.Kind == 20201 {
	}

	/*
	 kind 20202
	 p tag = user removes user they invited OR admin removes user
	 e tag = admin removes event
	*/
	if evt.Kind == 20202 {
		pTags := evt.Tags.GetAll([]string{"p"})
		for _, tag := range pTags {
			for _, user := range whitelist {
				/*
				 1: User in whitelist
				 2: Cant remove self
				 3: User should have invited user OR be relay admin
				*/
				if user.PublicKey == tag.Value() && evt.PubKey != tag.Value() && (user.InvitedBy == evt.PubKey || evt.PubKey == s.RelayPubkey) {
					log.Info().Str("user", tag.Value()).Msg("deleting user")
					deleteFromWhitelistRecursively(ctx, tag.Value())
				}
			}
		}
		if evt.PubKey == s.RelayPubkey {
			eTags := evt.Tags.GetAll([]string{"e"})
			for _, tag := range eTags {
				filter := nostr.Filter{
					IDs: []string{tag.Value()},
				}
				events, _ := db.QueryEvents(ctx, filter)

				for evt := range events {
					log.Info().Str("event", evt.ID).Msg("deleting event")
					err := db.DeleteEvent(ctx, evt)
					if err != nil {
						log.Warn().Err(err).Msg("failed to delete event")
					}
				}
			}
		}
	}

	return false, ""
}

func addToWhitelist(ctx context.Context, pubkey string, invitedBy string) error {
	if nostr.IsValidPublicKeyHex(pubkey) && !isPublicKeyInWhitelist(pubkey) {
		whitelist = append(whitelist, WhitelistEntry{PublicKey: pubkey, InvitedBy: invitedBy})
	}
	return saveWhitelist()
}

func removeFromWhitelist(ctx context.Context, pubkey string) error {
	idx := slices.IndexFunc(whitelist, func(we WhitelistEntry) bool { return we.PublicKey == pubkey })
	if idx == -1 {
		return nil
	}

	whitelist = append(whitelist[0:idx], whitelist[idx+1:]...)
	return saveWhitelist()
}

func loadWhitelist() error {
	if _, err := os.Stat("whitelist.json"); os.IsNotExist(err) {
		whitelist = []WhitelistEntry{}
		return nil
	} else if err != nil {
		return err
	}

	fileContent, err := os.ReadFile("whitelist.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &whitelist); err != nil {
		return err
	}

	return nil
}

func saveWhitelist() error {
	jsonBytes, err := json.Marshal(whitelist)
	if err != nil {
		return err
	}

	if err := os.WriteFile("whitelist.json", jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/nbd-wtf/go-nostr"
)

type User struct {
	Pk 		  string `json:"pk"`
	InvitedBy string `json:"invited_by"`
}

var whitelist []User

func whitelistRejecter(ctx context.Context, evt *nostr.Event) (reject bool, msg string) {

	// check if user in whitelist
	if !isPkInWhitelist(evt.PubKey) {
		return true, "You are not invited to this relay"
	}
	
	// 20201 = user invites new user
	if (evt.Kind == 20201) {
		pTags := evt.Tags.GetAll([]string{"p"})
		for _, tag := range pTags {
			if !isPkInWhitelist(tag.Value()) {
				if nostr.IsValidPublicKeyHex(tag.Value()) {
					whitelist = append(whitelist, User{Pk: tag.Value(), InvitedBy: evt.PubKey})
				}
			}
		}
	}

	// 20202 = user removes user they invited or admin removes invite
	if (evt.Kind == 20202) {
		pTags := evt.Tags.GetAll([]string{"p"})
		for _, tag := range pTags {
			for _, user := range whitelist {
				if user.Pk == tag.Value() && (user.InvitedBy == evt.PubKey || evt.PubKey == relayMaster) {
					deleteFromWhitelistRecursively(ctx, tag.Value())
				}
			}
		}
	}

	// 20203 = admin deletes event
	if (evt.Kind == 20203 && evt.PubKey == relayMaster) {
		eTags := evt.Tags.GetAll([]string{"e"})
		for _, tag := range eTags { 
			filter := nostr.Filter{
				IDs: []string{tag.Value()},
			}
			events, _ := db.QueryEvents(ctx, filter);

			for ev := range events {
				err := db.DeleteEvent(ctx, ev)
				if err != nil {
					log.Println("error while deleting event", err)
				}
			}
		}
	}

	return false, ""

}

func loadWhitelist () error {
	if _, err := os.Stat("whitelist.json"); os.IsNotExist(err) {
		whitelist = []User{}
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

func saveWhitelist () error {
	jsonBytes, err := json.Marshal(whitelist)
	if err != nil {
		return err
	}
	
	if err := os.WriteFile("whitelist.json", jsonBytes, 0644); err != nil {
		return err
	}
	
	return nil
}

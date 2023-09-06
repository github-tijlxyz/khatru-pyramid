package main

import (
	"context"
	"encoding/json"
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
	
	// add user(s) if needed
	if (evt.Kind == 20201) {
		pTags := evt.Tags.GetAll([]string{"p"})
		for _, tag := range pTags {
			if !isPkInWhitelist(tag.Value()) {
				whitelist = append(whitelist, User{Pk: tag.Value(), InvitedBy: evt.PubKey})
			}
		}
	}

	// remove user(s) if needed
	if (evt.Kind == 20202) {
		pTags := evt.Tags.GetAll([]string{"p"})
		for _, tag := range pTags {
			for _, user := range whitelist {
				if user.Pk == tag.Value() && user.InvitedBy == evt.PubKey {
					deleteFromWhitelistRecursively(tag.Value())
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

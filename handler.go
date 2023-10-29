package main

import (
	"context"
	"net/http"

	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/theplant/htmlgo"
)

// embed ui files

func inviteTreeHandler(w http.ResponseWriter, r *http.Request) {
	content := inviteTreePageHTML(r.Context(), InviteTreePageParams{
		LoggedUser: getLoggedUser(r),
	})
	htmlgo.Fprint(w, baseHTML(content), r.Context())
}

func addToWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)

	pubkey := r.PostFormValue("pubkey")
	if pfx, value, err := nip19.Decode(pubkey); err == nil && pfx == "npub" {
		pubkey = value.(string)
	}

	if err := addToWhitelist(r.Context(), pubkey, loggedUser); err != nil {
		http.Error(w, "failed to add to whitelist: "+err.Error(), 500)
		return
	}
	content := buildInviteTree(r.Context(), s.RelayPubkey, loggedUser)
	htmlgo.Fprint(w, content, r.Context())
}

func removeFromWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	loggedUser := getLoggedUser(r)
	pubkey := r.PostFormValue("pubkey")
	if err := removeFromWhitelist(r.Context(), pubkey, loggedUser); err != nil {
		http.Error(w, "failed to remove from whitelist: "+err.Error(), 500)
		return
	}
	content := buildInviteTree(r.Context(), s.RelayPubkey, loggedUser)
	htmlgo.Fprint(w, content, r.Context())
}

func reportsViewerHandler(w http.ResponseWriter, r *http.Request) {
	// var formattedReportsData template.HTML = ""

	// events, _ := db.QueryEvents(context.Background(), nostr.Filter{
	// 	Kinds: []int{1984},
	// 	Limit: 52,
	// })

	// type Report struct {
	// 	ID         string
	// 	ByUser     string
	// 	AboutUser  string
	// 	AboutEvent string
	// 	Type       string
	// 	Content    string
	// }

	// for ev := range events {
	// 	pTag := ev.Tags.GetFirst([]string{"p"})

	// 	eTag := ev.Tags.GetFirst([]string{"e"})
	// 	if pTag != nil {
	// 		typeReport := eTag.Relay()[6:]
	// 		if typeReport == "" {
	// 			typeReport = pTag.Relay()[6:]
	// 		}
	// 		report := Report{
	// 			ID:         ev.ID,
	// 			ByUser:     ev.PubKey,
	// 			AboutUser:  pTag.Value(),
	// 			AboutEvent: eTag.Value(),
	// 			Type:       typeReport,
	// 			Content:    ev.Content,
	// 		}
	// 		// get AboutEvent content, note1 ect
	// 		formattedReportsData += template.HTML(fmt.Sprintf(`
	// 		<div>
	// 		<p><b>Report %v</b></p>
	// 		<p>By User: <a class="user" href="nostr:%v">%v</a></p>
	// 		<p>About User: <a class="user" href="nostr:%v">%v</a></p>`,
	// 			report.ID,
	// 			getUserInfo(context.Background(), report.ByUser).Npub,
	// 			getUserInfo(context.Background(), report.ByUser).Name,
	// 			getUserInfo(context.Background(), report.AboutUser).Npub,
	// 			getUserInfo(context.Background(), report.AboutUser).Name,
	// 		))
	// 		if report.AboutEvent != "" {
	// 			// fetch event data
	// 			aboutEvents, _ := db.QueryEvents(context.TODO(), nostr.Filter{
	// 				IDs: []string{report.AboutEvent},
	// 			})
	// 			for aboutEvent := range aboutEvents {
	// 				formattedReportsData += template.HTML(fmt.Sprintf(`
	// 				<p>
	// 				About Event: <ul>
	// 				<p>Kind: %v</p>
	// 				<p>Tags: %v</p>
	// 				<p>Content: %v</p>
	// 				</ul>
	// 				</p>`,
	// 					template.HTMLEscaper(aboutEvent.Kind),
	// 					template.HTMLEscaper(aboutEvent.Tags),
	// 					template.HTMLEscaper(aboutEvent.Content),
	// 				))
	// 			}
	// 		}
	// 		formattedReportsData += template.HTML(fmt.Sprintf(`
	// 		<p>Type: %v</p>`,
	// 			report.Type,
	// 		))
	// 		if report.Content != "" {
	// 			formattedReportsData += template.HTML(fmt.Sprintf(`
	// 			<p>Content: %v</p>
	// 			<div>
	// 			<button data-actionarg='[["e", "%v"],["p", "%v"]]' class="removefromrelay">Ban Reported User and Remove Report</button>
	// 			<button data-actionarg='[["e", "%v"]]' class="removefromrelay">Remove This Report</button>
	// 			<button data-actionarg='[["p", "%v"]]' class="removefromrelay">Ban User who wrote report</button>
	// 			</div>
	// 			</div>
	// 			<hr />`,
	// 				template.HTMLEscaper(report.Content),
	// 				template.HTMLEscaper(report.ID),
	// 				template.HTMLEscaper(report.AboutUser),
	// 				template.HTMLEscaper(report.ID),
	// 				template.HTMLEscaper(report.ByUser),
	// 			))
	// 		}
	// 	}
	// }

	// data := map[string]interface{}{
	// 	"Relayname":        s.RelayName,
	// 	"Relaydescription": s.RelayDescription,
	// 	"Pagetitle":        "Reports Viewer",
	// 	"Pagecontent":      formattedReportsData,
	// }

	// tmpl, err := template.ParseFS(dist, "ui/dist/index.html")
	// if err != nil {
	// 	http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// // Execute the template with the provided data and write it to the response
	// err = tmpl.Execute(w, data)
	// if err != nil {
	// 	http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	content := homePageHTML(r.Context(), HomePageParams{
		RelayOwnerInfo: getUserInfo(context.Background(), s.RelayPubkey),
	})
	htmlgo.Fprint(w, baseHTML(content), r.Context())
}

# Khatru Invite

A relay based on [Khatru](https://github.com/fiatjaf/khatru) with a invite hierarchy feature.

### Deploy with docker

1. create and manually add a pubkey to users.json: `touch users.json && echo '{"your nostr hex pubkey":""}' > users.json`
2. deploy with docker: `docker run -p 3334:3334 -v ./users.json:/app/users.json -v ./db:/app/db -e DOMAIN=yourdomain.example.com -e RELAY_NAME="your relay name" -e RELAY_PUBKEY="your nostr hex pubkey" tijlxyz/khatru-pyramid:latest`


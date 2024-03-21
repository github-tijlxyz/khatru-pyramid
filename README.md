# Khatru Pyramid

A relay based on [Khatru](https://github.com/fiatjaf/khatru) with a invite hierarchy feature.

### Deploy with docker

```
$ docker run \
    -p 3334:3334 \
    -v ./users.json:/app/users.json \
    -v ./db:/app/db \
    -e DOMAIN="yourdomain.example.com" \
    -e RELAY_NAME="your relay name" \
    -e RELAY_PUBKEY="your nostr hex pubkey" \
    tijlxyz/khatru-pyramid:latest
```

### Deploy with

 - [YunoHost](https://github.com/YunoHost-Apps/khatru-pyramid_ynh) ([app catalog](https://apps.yunohost.org/catalog) [pending](https://github.com/YunoHost/apps/pull/2077))
 - [Cloudron](https://github.com/github-tijlxyz/khatru-pyramid_cloudron) ([app catalog](https://www.cloudron.io/store/index.html) [pending](https://forum.cloudron.io/topic/11146/khatru-pyramid-a-nostr-relay))

### Manually build

```
$ git clone https://github.com/github-tijlxyz/khatru-pyramid 
$ cd khatru-pyramid
$ go build # or run
$ DOMAIN=example.com RELAY_NAME=test RELAY_PUBKEY=yourpubkey ./khatru-pyramid
```

### Configuration

Format `users.json` as follows:

```json
{ "[user_pubkey_hex]": "[invited_by_pubkey_hex]" }
```

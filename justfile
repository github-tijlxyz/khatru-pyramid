dev:
    ag -l --go | entr -r godotenv go run .

build:
    CC=musl-gcc go build -ldflags='-linkmode external -extldflags "-static"' -o ./khatru-invite

deploy: build
    ssh root@turgot 'systemctl stop pyramid';
    scp khatru-invite turgot:pyramid/khatru-invite
    ssh root@turgot 'systemctl start pyramid'

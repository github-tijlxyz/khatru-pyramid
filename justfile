dev:
    ag -l --go | entr -r godotenv go run .

build:
    CC=musl-gcc go build -ldflags='-linkmode external -extldflags "-static"' -o ./khatru-invite

deploy: build
    ssh root@cantillon 'systemctl stop pyramid';
    scp khatru-invite cantillon:pyramid/khatru-invite
    ssh root@cantillon 'systemctl start pyramid'

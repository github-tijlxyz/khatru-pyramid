dev:
    ag -l --go | entr -r godotenv go run .

build:
    CC=musl-gcc go build -ldflags='-linkmode external -extldflags "-static"' -o ./khatru-pyramid

deploy target: build
    ssh root@{{target}} 'systemctl stop pyramid';
    scp khatru-pyramid {{target}}:pyramid/khatru-invite
    ssh root@{{target}} 'systemctl start pyramid'

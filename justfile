dev:
    fd 'go|templ' | entr -r bash -c 'just templ && godotenv go run .'

build: templ
    CC=musl-gcc go build -ldflags='-linkmode external -extldflags "-static"' -o ./khatru-pyramid

templ:
    templ generate

deploy target: build
    ssh root@{{target}} 'systemctl stop pyramid';
    scp khatru-pyramid {{target}}:pyramid/khatru-invite
    ssh root@{{target}} 'systemctl start pyramid'

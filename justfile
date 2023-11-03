dev:
    godotenv go run .

build:
    CC=musl-gcc go build -ldflags='-s -w -linkmode external -extldflags "-static"' -o ./khatru-invite

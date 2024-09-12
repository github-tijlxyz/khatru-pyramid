FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o khatru-invite .

FROM ubuntu:latest
COPY --from=builder /app/khatru-invite /app/
ENV DATABASE_PATH="/app/db"
ENV USERDATA_PATH="/app/users.json"
CMD ["/app/khatru-invite"]


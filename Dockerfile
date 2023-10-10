FROM node:latest AS ui-builder

WORKDIR /app

COPY . .

WORKDIR /app/ui

RUN yarn install
RUN yarn build


FROM golang:1.20 AS go-builder

WORKDIR /app

COPY . .

COPY --from=ui-builder /app/ui/dist /app/ui/dist

RUN go build -o app


FROM golang:1.20

WORKDIR /app

COPY --from=go-builder /app/app /app/app

COPY --from=ui-builder /app/ui/dist /app/ui/dist


EXPOSE 3334

CMD ["./app"]

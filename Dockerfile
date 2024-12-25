FROM golang:latest AS app

WORKDIR /usr/src/app

COPY . .

# Build latinaapi
RUN go mod download && go mod tidy && go mod verify
RUN go build -o ./app ./cmd/main.go

ENV GIN_MODE=release
EXPOSE 8080

CMD ["./app"]
FROM golang:1.25.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest
RUN go install github.com/mfridman/tparse@latest
RUN go install github.com/oligot/go-mod-upgrade@latest
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.6.2
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

CMD sh -c "air -c .air.${SERVICE}.toml"

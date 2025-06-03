FROM golang:1.23.9

WORKDIR app


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/bin/main ./cmd/app/main.go


CMD ["/app/bin/main"]

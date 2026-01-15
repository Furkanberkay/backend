FROM golang:1.25.3-alpine3.20

WORKDIR /src/app

RUN go install github.com/air-verse/air@latest

COPY . .
RUN go mod tidy

CMD ["air", "-c", ".air.toml"]

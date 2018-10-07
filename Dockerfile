FROM golang:1.11 AS build

WORKDIR /mu2-encode

ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN CGO_ENABLED=0 go build -o mu2-encode main.go

FROM alpine:latest AS RUN

WORKDIR /app/
COPY --from=build /mu2-encode/mu2-encode mu2-encode
RUN apk add --no-cache ca-certificates ffmpeg

CMD ["./mu2-encode"]

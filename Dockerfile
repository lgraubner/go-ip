ARG GO_VERSION=1.17.6

## build container
FROM golang:${GO_VERSION}-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# install git, timezone data
RUN apk update && apk add --no-cache git tzdata

RUN adduser \
    --disabled-password \
    --gecos "" \
    --no-create-home \
    --uid "1002" \
    "nonroot"

WORKDIR /src

# use hack for possible non existent go.sum
# https://stackoverflow.com/a/46801962
COPY ./go.mod ./go.sum* .

RUN go mod download

COPY . .

RUN go test -v

RUN go build -o /bin/app .

## prod container
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /bin/app /bin/app

USER nonroot:nonroot

ENTRYPOINT ["/bin/app"]

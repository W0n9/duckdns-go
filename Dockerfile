FROM golang:1.15-alpine3.12 as build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /tmp/duckdns-go

RUN apk --no-cache add alpine-sdk ca-certificates
COPY . .
RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags '-s -w' -o duckdns-go ./

FROM scratch
LABEL name="duckdns-go"

WORKDIR /root
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /tmp/duckdns-go/duckdns-go duckdns-go

CMD ["./duckdns-go", "-update-ip"]
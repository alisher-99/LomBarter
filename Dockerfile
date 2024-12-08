# Step 1: Prepare
FROM golang:1.20-alpine3.18 as builder

COPY . /app
WORKDIR /app

RUN cd cmd/app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags app -installsuffix app -o /bin/app -mod vendor

# Step 2: Final
FROM alpine:3.18

COPY --from=builder /bin/app /app
COPY --from=builder /app/version /version
COPY --from=builder /app/swagger /swagger
COPY --from=builder /app/config /config

CMD ["/app"]
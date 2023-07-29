FROM golang AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"'

FROM scratch

EXPOSE 8080

COPY --from=builder /app/share-secret /share-secret

ENTRYPOINT ["/share-secret"]

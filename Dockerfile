FROM golang:1.11-alpine as builder

LABEL maintainer="AnotherGitProfile"

RUN apk add --update --no-cache ca-certificates git
WORKDIR /bin/
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /bin/main

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/main /bin/main
COPY --from=builder /bin/config.yml /etc/eos_exporter/config.yml

EXPOSE 9386
ENTRYPOINT [ "/bin/main" ]
CMD [ "-config.file=/etc/eos_exporter/config.yml" ]
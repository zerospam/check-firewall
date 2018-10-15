FROM golang:1.10-alpine

COPY . /go/src/github.com/zerospam/checkfirewall
WORKDIR /go/src/github.com/zerospam/checkfirewall
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o firewall-checker .

FROM scratch

COPY --from=0 /go/src/github.com/zerospam/checkfirewall/firewall-checker /firewall-checker

CMD ["/firewall-checker"]

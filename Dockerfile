FROM golang:1.10.1-stretch as builder

COPY database/database.go /go/src/github.com/Its-Alex/dhcp4-reservation/database/database.go
COPY dhcp/dhcp.go /go/src/github.com/Its-Alex/dhcp4-reservation/dhcp/dhcp.go
COPY cmd/root.go /go/src/github.com/Its-Alex/dhcp4-reservation/cmd/root.go
COPY main.go /go/src/github.com/Its-Alex/dhcp4-reservation
COPY Gopkg.lock /go/src/github.com/Its-Alex/dhcp4-reservation
COPY Gopkg.toml /go/src/github.com/Its-Alex/dhcp4-reservation

WORKDIR /go/src/github.com/Its-Alex/dhcp4-reservation

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    && dep ensure -v -vendor-only \
    # Install gox
    && go get -u github.com/mitchellh/gox \
    && gox -osarch="linux/386" .

FROM alpine:3.7

# Copy executalle from builder
COPY --from=builder go/src/github.com/Its-Alex/dhcp4-reservation/dhcp4-reservation_linux_386 /usr/local/bin/dhcp4

EXPOSE 67

CMD ["/usr/local/bin/dhcp4"]
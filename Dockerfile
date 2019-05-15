FROM golang:1.12

WORKDIR /go/src/github.com/zom-bi/ovpn-certman
ADD . .
RUN \
    go get -tags="dev" -v github.com/zom-bi/ovpn-certman && \
    go get github.com/shurcooL/vfsgen/cmd/vfsgendev && \
    go generate github.com/zom-bi/ovpn-certman/assets && \
    go build -tags="netgo"

FROM scratch
ENV \
    OAUTH2_CLIENT_ID="" \
    OAUTH2_CLIENT_SECRET="" \
    OAUTH2_AUTH_URL="https://gitlab.example.com/oauth/authorize" \
    OAUTH2_TOKEN_URL="https://gitlab.example.com/oauth/token" \
    OAUTH2_REDIRECT_URL="https://vpn.example.com/login/oauth2/redirect" \
    USER_ENDPOINT="https://gitlab.example.com/api/v4/user" \
    VPN_DEV="tun" \
    VPN_HOST="vpn.example.com" \
    VPN_PORT="1194" \
    VPN_PROTO="udp"
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/github.com/zom-bi/ovpn-certman/ovpn-certman /certman
ENTRYPOINT ["/certman"]

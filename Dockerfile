FROM golang:1.9

WORKDIR /go/src/git.klink.asia/paul/certman
ADD . .
RUN \
    go get github.com/shurcooL/vfsgen/cmd/vfsgendev && \
    go generate git.klink.asia/paul/certman/assets && \
    go get -v git.klink.asia/paul/certman && \
    go build -tags netgo

FROM scratch
ENV \
    OAUTH2_CLIENT_ID="" \
    OAUTH2_CLIENT_SECRET="" \
    APP_KEY="" \
    OAUTH2_AUTH_URL="https://gitlab.example.com/oauth/authorize" \
    OAUTH2_TOKEN_URL="https://gitlab.example.com/oauth/token" \
    USER_ENDPOINT="https://gitlab.example.com/api/v4/user" \
    OAUTH2_REDIRECT_URL="https://certman.example.com/login/oauth2/redirect"
COPY --from=0 /go/src/git.klink.asia/paul/certman/certman /
ENTRYPOINT ["/certman"]

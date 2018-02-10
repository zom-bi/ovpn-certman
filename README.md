# Certman
Certman is a simple certificate manager web service for OpenVPN.

## Installation
### Binary
There are prebuilt binary files for this application. They are statically
linked and have no additional dependencies. Supported plattforms are:
 * Windows (XP and up)
 * Linux (2.6.16 and up)
 * Linux ARM (for raspberry pi, 3.0 and up)
Simply download them from the "artifacts" section of this project.
### Docker
A prebuilt docker image (10MB) is available:
```bash
docker pull docker.klink.asia/paul/certman
```
### From Source-Docker
You can easily build your own docker image from source
```bash
docker build -t docker.klink.asia/paul/certman .
```

## Configuration
Certman assumes the root certificates of the VPN CA are located in the same
directory as the binary, If that is not the case you need to copy over the
`ca.crt` and `ca.key` files before you are able to generate certificates
with this tool.

Additionally, the project is configured by the following environment
variables:
 * `OAUTH2_CLIENT_ID` the Client ID, assigned during client registration
 * `OAUTH2_CLIENT_SECRET` the Client secret, assigned during client registration
 * `OAUTH2_AUTH_URL` the URL to the "/authorize" endpoint of the identity provider
 * `OAUTH2_TOKEN_URL` the URL to the "/token" endpoint of the identity provider
 * `OAUTH2_REDIRECT_URL` the redirect URL used by the app, usually the hostname suffixed by "/login/oauth2/redirect"
 * `USER_ENDPOINT` the URL to the Identity provider user endpoint, for gitlab this is "/api/v4/user". The "username" attribute of the returned JSON will used for authentication.
 * `APP_KEY` random ASCII string, 32 characters in length. Used for cookie generation.
 * `APP_LISTEN` port and ip to listen on, e.g. `:8000` or `127.0.0.1:3000`
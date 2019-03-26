# authproxy

Package authproxy provides a reverse proxy that appends a Google ID Token to request headers. It also exposes a sub package that can be used as a server immediately.

## Usage

#### Docker

`docker run -p 9090:9090 -e REVERSE_PROXY_URL=mySecureServer.com -e SERVICE_ACCOUNT_KEY=<base64-encoded-json-key> marwan91/authproxy"`

#### Go 

You can use the server directly as such

```bash
~ go install marwan.io/authproxy/cmd/authproxy
~ export REVERSE_PROXY_URL=mySecureServer.com
~ export SERVICE_ACCOUNT_KEY=<base64-encoded-json-key>
~ authproxy
# listening on port :9090
# you can set the PORT env if you want to user a port other than 9090
```

Or you can use it programtically: 

```bash
~ go install marwan.io/authproxy
```

```golang
package main

import "marwan.io/authproxy"

func main() {
    handler, err := authproxy.GetProxyHandler("mySecureProxy.com", "base64EncodedKey")
    // handle err and use the proxy handler
    // ...
}
```
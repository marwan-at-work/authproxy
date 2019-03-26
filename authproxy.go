// Package authproxy provides a reverse proxy
// that appends an authorization token to request headers.
// It also exposes a sub package that can be used as a
// server immediately with Google ID Tokens.
package authproxy

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/oauth2/google"
)

// GetProxyHandler takes the URL that should be proxied to and a base64 encoded
// Service Account secret key. For every request, it will use the secret key
// to create an ID Token and appends to proxied requests as "Authorization: Bearer <token>"
func GetProxyHandler(proxyURL, serviceAccount string) (http.Handler, error) {
	target, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse cloud run url: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	targetQuery := target.RawQuery
	proxy.Director = func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	if serviceAccount == "METADATA_FLAVOR" {
		proxy.Transport = &metaTransport{proxyURL}
	} else {
		cl, err := client(proxyURL, serviceAccount)
		if err != nil {
			return nil, err
		}
		proxy.Transport = cl.Transport
	}

	return proxy, nil
}

func client(targetAudience, key string) (*http.Client, error) {
	bts, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64 service account: %v", err)
	}
	cfg, err := google.JWTConfigFromJSON(bts)
	if err != nil {
		return nil, fmt.Errorf("could not get jwt of service account: %v", err)
	}
	cfg.PrivateClaims = map[string]interface{}{"target_audience": targetAudience}
	cfg.UseIDToken = true
	return cfg.Client(context.Background()), nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

type metaTransport struct {
	url string
}

func (m *metaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tok, err := m.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)

	return http.DefaultTransport.RoundTrip(req)
}

func (m *metaTransport) getToken() (string, error) {
	url := "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Set("audience", m.url)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not request creds: %v", err)
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read creds: %v", err)
	}

	return strings.TrimSpace(string(bts)), nil
}

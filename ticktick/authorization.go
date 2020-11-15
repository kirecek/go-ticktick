package ticktick

import (
	"context"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

// list of ticktick permission scopes.
const (
	ScopeReadTask  = "tasks:read"
	ScopeReadWrite = "tasks:write"
)

const (
	defaultAuthURL  = "https://ticktick.com/oauth/authorize"
	defaultTokenURL = "https://ticktick.com/oauth/token"
)

// OAuthConfig provides minimal OAuth configuraion in order to generate TickTick access token.
//
// https://developer.ticktick.com/docs/index.html#/openapi?id=authorization
type OAuthConfig struct {
	Scopes       []string
	ClientID     string
	ClientSecret string
}

// NewOAuthClient creates a HTTP client authorized with TickTick API.
func NewOAuthClient(ctx context.Context, config *OAuthConfig) *http.Client {
	cfg := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  defaultAuthURL,
			TokenURL: defaultTokenURL,
		},
	}

	cacheFile := tokenCacheFile(cfg)
	token, err := tokenFromFile(cacheFile)
	if err != nil {
		token = tokenFromWeb(ctx, cfg)
		saveToken(cacheFile, token)
	} else {
		log.Printf("Using cached token %#v from %q", token, cacheFile)
	}

	return cfg.Client(ctx, token)
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)

	// ts := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	ts, err := listenAndServe(":8000", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))
	if err != nil {
		panic(err)
	}
	defer ts.Close()

	// config.RedirectURL = ts.URL
	config.RedirectURL = "http://127.0.0.1:8000"
	authURL := config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	log.Printf("Authorize this app at: %s", authURL)
	open.Run(authURL)
	code := <-ch
	log.Printf("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return token
}

func saveToken(file string, token *oauth2.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}

func tokenCacheFile(config *oauth2.Config) string {
	hash := fnv.New32a()
	hash.Write([]byte(config.ClientID))
	hash.Write([]byte(config.ClientSecret))
	hash.Write([]byte(strings.Join(config.Scopes, " ")))
	fn := fmt.Sprintf("go-ticktick-tok%v", hash.Sum32())
	return filepath.Join(osUserCacheDir(), url.QueryEscape(fn))
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth2.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func osUserCacheDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Caches")
	case "linux", "freebsd":
		return filepath.Join(os.Getenv("HOME"), ".cache")
	}
	log.Printf("TODO: osUserCacheDir on GOOS %q", runtime.GOOS)
	return "."
}

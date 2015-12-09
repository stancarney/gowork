package gowork

import (
	"strings"
	"github.com/gorilla/context"
	"net/http"
	"log"
	"net/http/httputil"
	"github.com/gorilla/sessions"
)

type getSession func(r *http.Request, name string) (*sessions.Session, error)
type getUser func(ctx Context, id string) (user interface{}, err error)
type newRequestContext func(r *http.Request) Context

type DebugHandler struct {
	Debug bool
	Name  string
}

func (s *DebugHandler) New(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			log.Printf("+++++++++++++++++++++++++++++++++++++ %s ++++++++++++++++++++++++++++++++++++++++++", s.Name)
			data, _ := httputil.DumpRequest(r, true)
			log.Printf("REQUEST DUMP (To disable this set debug=false in config.toml): %s\n\n", data)
		}
		handler.ServeHTTP(w, r)
	})
}

type StaticAssetHandler struct {
	AssetUrlPrefix string
	AssetUrls      []string
	AssetFilePath  string
}

func (s *StaticAssetHandler) New(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		if s.IsStaticAsset(url) {
			http.ServeFile(w, r, s.AssetFilePath + strings.TrimPrefix(url, s.AssetUrlPrefix))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (s *StaticAssetHandler) IsStaticAsset(url string) bool {
	for _, v := range s.AssetUrls {
		if strings.HasPrefix(url, s.AssetUrlPrefix + v) {
			return true
		}
	}
	return false
}

// ContextHandler is a simple http.Handler that attaches the configured Impl Context to the Gorilla Context. It effectively hides the Gorilla Context from the rest of the application.
type ContextHandler struct {
	NewRequestContext newRequestContext
}

func (s *ContextHandler) New(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := s.NewRequestContext(r)
		context.Set(r, ReqCtx, ctx)
		handler.ServeHTTP(w, r)
	})
}

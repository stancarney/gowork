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
	Name string
}

func (s *DebugHandler) New(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.Debug {
			log.Printf("++++++++++++++++++++++++++++++++ %s +++++++++++++++++++++++++++++++++++++++++++++++", s.Name)
			data, _ := httputil.DumpRequest(r, true)
			log.Printf("REQUEST DUMP (To disable this set debug=false in config.toml): %s\n\n", data)
		}
		handler.ServeHTTP(w, r)
	})
}

// AuthHandler manages the setting of the userId into the session when the required cookie is present.
// This handler requires ContextHandler to be run first as it depends on the Context being present.
type AuthHandler struct {
	GetSession       getSession
	GetUser          getUser
	StaticAssetPaths []string
	InsecureUrls     []string
	CookieName       string
	IndexFile        string
	Prefix           string
}

func (s *AuthHandler) New(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Do this right away so we have the ip stored on the context.
		ctx := GetContext(r)
		ctx.Put("ip", GetIp(r))

		url := r.URL.String()
		if s.IsStaticAsset(url) {
			handler.ServeHTTP(w, r)
			return
		}

		if s.IsInsecureUrl(url) {
			handler.ServeHTTP(w, r)
			return
		}

		session, err := s.GetSession(r, s.CookieName)
		if err != nil {
			//Session expired
			WriteErrorToJSON(w, 401, "Session Expired, please re-login.")
			return
		}

		//Try to load the user and set it on the context on each request. Regardless if the individual is logged in or not.
		if userid, ok := session.Values["userid"]; ok {
			if user, err := s.GetUser(ctx, userid.(string)); err == nil {
				ctx.Put("user", user)
			}
		}

		//if auth exists, the user is logged in.
		if _, ok := session.Values["auth"]; ok {
			session.Save(r, w) //Gorilla's context is cleared out after the next handler (Gorilla's handlers) run, therefore we have to save the session before it runs in order to update the last access time.
			handler.ServeHTTP(w, r)
			return
		}

		//Calls to /sys are from angular.
		if strings.HasPrefix(url, "/sys") {
			WriteErrorToJSON(w, 401, "Login Required")
			return
		}

		//This would be a direct call to a page (say /client) without a valid session. The index page needs to be returned so angular can request the page (/client) and get the 401 code above.
		http.ServeFile(w, r, s.IndexFile)
	})
}

func (s *AuthHandler) IsStaticAsset(url string) bool {
	for _, v := range s.StaticAssetPaths {
		if strings.HasPrefix(url, s.Prefix + v) {
			return true
		}
	}
	return false
}

func (s *AuthHandler) IsInsecureUrl(url string) bool {
	for _, v := range s.InsecureUrls {
		if strings.HasPrefix(url, v) {
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

package gowork

import (
	"github.com/gorilla/sessions"
	"net/http"
	"github.com/gorilla/securecookie"
	"encoding/base32"
	"strings"
)

//RequestContext is just a wrapper on SimpleRequestContext for the time being.
type RequestContext struct {
	SimpleRequestContext
}

type SessionProvider interface {
	GetSession(ctx Context, id string) (*Session, error)
	CreateSession(ctx Context, session *Session) error
	UpdateSession(ctx Context, session *Session) (err error)
}

type CookieStore struct {
	Codecs          []securecookie.Codec
	Options         *sessions.Options // default configuration
	SessionProvider SessionProvider
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new FilesystemStore is 4096.
func (s *CookieStore) MaxLength(l int) {
	for _, c := range s.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *CookieStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *CookieStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = false
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			//load returns an error when the session is not found. We don't want it returned out of this function.
			err := s.load(GetContext(r), session)
			if err != nil {
				session.ID = ""
				session.IsNew = true
			}
		}
	}

	return session, err
}

// Save adds a single session to the response.
func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.ID == "" { //New Session
		session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		if err := s.save(GetContext(r), session); err != nil {
			return err
		}

		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	} else { //Existing Session
		if err := s.update(GetContext(r), session); err != nil {
			return err
		}
	}
	return nil
}

func (s *CookieStore) ConvertToInterfaceMap(in map[string]string) (out map[interface{}]interface{}) {
	out = make(map[interface{}]interface{}, len(in)) //TODO:Stan I think we can encode the value as a string regardless of type maybe...
	for k, v := range in {
		out[k] = v
	}
	return
}

func (s *CookieStore) ConvertToStringMap(in map[interface{}]interface{}) (out map[string]string) {
	out = make(map[string]string, len(in)) //TODO:Stan I think we can encode the value as a string regardless of type maybe...
	for k, v := range in {
		out[k.(string)] = v.(string)
	}
	return
}

// save writes encoded session.Values to a DB.
func (s *CookieStore) save(ctx Context, session *sessions.Session) error {
	t := CurrentTime()
	cs := Session{
		Id: session.ID,
		Created: t,
		LastAccess: t,
		Values: s.ConvertToStringMap(session.Values),
		UserId: ctx.GetString("userid"), //We store the UserId directly on model.Session to make it easier to report via the UI.
		Version: 0,
	}

	return s.SessionProvider.CreateSession(ctx, &cs)
}

// update writes encoded session.Values to a DB.
func (s *CookieStore) update(ctx Context, session *sessions.Session) (err error) {
	cs, err := s.SessionProvider.GetSession(ctx, session.ID)
	if err != nil {
		return
	}

	cs.LastAccess = CurrentTime()
	cs.Values = s.ConvertToStringMap(session.Values)

	err = s.SessionProvider.UpdateSession(ctx, cs)
	return
}

// load reads and decodes sessions contents from Cassandra into session.Values.
// Returns an error if the session is not found.
func (s *CookieStore) load(ctx Context, session *sessions.Session) error {
	cs, err := s.SessionProvider.GetSession(ctx, session.ID)
	if err != nil {
		return err
	}

	session.Values = s.ConvertToInterfaceMap(cs.Values)
	return nil
}


package gowork

import (
	"github.com/gorilla/context"
	"net/http"
	"sync"
)

const (
	ReqCtx ContextKey = "REQ_CTX"
)

type ContextKey string

//Context wraps the Gorilla context and is used to store things unique to each HTTP Request.
type Context interface {
	Get(key interface{}) interface{}
	Put(key interface{}, value interface{})

	//Utility method to reduce some repeating code
	GetString(key interface{}) string
}

//SimpleContext is the simplest implementation of a Context used by the application.
type SimpleContext struct {
	values map[interface{}]interface{}
	mutex  *sync.Mutex
}

func (ac *SimpleContext) Get(key interface{}) interface{} {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	return ac.values[key]
}

func (ac *SimpleContext) Put(key interface{}, value interface{}) {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	ac.values[key] = value
	return
}

func (ac *SimpleContext) GetString(key interface{}) (str string) {
	r := ac.Get(key)
	if r != nil {
		str = r.(string)
	}
	return
}

type SimpleRequestContext struct {
	Request *http.Request
	SimpleContext
}

func NewSimpleContext() *SimpleContext {
	return &SimpleContext{make(map[interface{}]interface{}), &sync.Mutex{}}
}

func NewSimpleRequestContext(r *http.Request) *SimpleRequestContext {
	return &SimpleRequestContext{r, *NewSimpleContext()}
}

func GetContext(r *http.Request) (ctx Context) {
	//Gorilla's context provides the locking. No need to duplicate it here.
	if c := context.Get(r, ReqCtx); c != nil {
		ctx = c.(Context)
	}
	return
}

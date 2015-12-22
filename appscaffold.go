package gowork

import (
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"encoding/json"
)

type createEvent func(ctx Context, code EventCode, err interface{}, description ...string)

type AppScaffold struct {
	CreateEvent createEvent
}

//GetEntities uses reflection to look up the call function on app.App.Impl and call it.
//This was introduced to minimize the amount of duplicated code in the controllers (i.e. sys).
func (a * AppScaffold) GetEntities(w http.ResponseWriter, r *http.Request, user User, function interface{}, perm Permission) {

	if user == nil {
		WriteErrorToJSON(w, 500, "User session not found")
		return
	}
	
	ctx := GetContext(r)

	if !user.HasPermission(perm) {
		a.CreateEvent(ctx, DENIED, nil, GetFunctionName(function))
		WriteErrorToJSON(w, 403, "Permission denied")
		return
	}

	a.CreateEvent(ctx, READ, nil, GetFunctionName(function))

	argCtx := reflect.ValueOf(ctx)
	argDate := reflect.ValueOf(GetDate(r))
	argLimit := reflect.ValueOf(GetLimit(r))

	f := reflect.ValueOf(function)

	in := make([]reflect.Value, 0, 3)
	in = append(in, argCtx)
	if f.Type().NumIn() == 3 {
		in = append(in, argDate)
	}
	in = append(in, argLimit)

	result := f.Call(in)

	data := result[0].Interface()
	err := result[1].Interface()
	if err != nil {
		WriteErrorToJSON(w, 500, err)
		return
	}

	WriteJSON(w, data)
	return
}

//GetEntity uses reflection to look up the call function and call it along with the id request parameter (path actually).
//It is almost identical to GetEntities with the exception of passing two args to the call function. The function to call and the id.
func (a * AppScaffold) GetEntity(w http.ResponseWriter, r *http.Request, user User, function interface{}, perm Permission) {

	if user == nil {
		WriteErrorToJSON(w, 500, "User session not found")
		return
	}

	ctx := GetContext(r)
	vars := mux.Vars(r)

	if !user.HasPermission(perm) {
		a.CreateEvent(ctx, DENIED, nil, GetFunctionName(function), vars["id"])
		WriteErrorToJSON(w, 403, "Permission denied")
		return
	}

	f := reflect.ValueOf(function)
	arg0 := reflect.ValueOf(ctx)
	arg1 := reflect.ValueOf(vars["id"])
	in := []reflect.Value{arg0, arg1}
	result := f.Call(in)

	data := result[0].Interface()
	err := result[1].Interface()
	if err != nil {
		a.CreateEvent(ctx, ERROR, err, GetFunctionName(function), "Error", err.(error).Error())
		WriteErrorToJSON(w, 404, err)
		return
	}

	a.CreateEvent(ctx, READ, nil, GetFunctionName(function))

	WriteJSON(w, data)
	return
}

//EntityOp fulfills Create, Update, and Delete functions. On a delete the id var must be present (i.e. route defined like: /myentity/{id}) and entity must be nil. The provided function will be called deleting the correct entity.
//With Create and Update entity mush be a pointer to the struct of the correct model so a json.Decode can populate it with information present in the request.
func (a * AppScaffold) EntityOp(w http.ResponseWriter, r *http.Request, user User, entity interface{}, function interface{}, successEventCode EventCode, perm Permission) {

	if user == nil {
		WriteErrorToJSON(w, 500, "User session not found")
		return
	}

	ctx := GetContext(r)

	if !user.HasPermission(perm) {
		a.CreateEvent(ctx, DENIED, nil, GetFunctionName(function))
		WriteErrorToJSON(w, 403, "Permission denied")
		return
	}

	vars := mux.Vars(r)

	if entity != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(entity); err != nil {
			panic(err)
			WriteErrorToJSON(w, 500, err)
			return
		}
	} else {
		entity = vars["id"]
	}

	f := reflect.ValueOf(function)

	arg0 := reflect.ValueOf(ctx)
	arg1 := reflect.ValueOf(entity)
	in := []reflect.Value{arg0, arg1}
	result := f.Call(in)

	if err := result[0].Interface(); err != nil {
		WriteErrorToJSON(w, 500, err)
		return
	}

	a.CreateEvent(ctx, successEventCode, nil, GetFunctionName(function))
	WriteJSON(w, entity)
	return
}
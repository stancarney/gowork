package gowork

import (
	"reflect"
	"sync"
	"fmt"
	"time"
	"log"
)

type MemTable struct {
	Table map[string]interface{}
	mutex *sync.Mutex
	Debug bool
}

// Create makes a new table entry with the provided id. An error is never returned from this function and is intended for use in overriding functions.
func (mt *MemTable) Create(id string, o interface{}) error {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	if mt.Debug {
		log.Printf("%s: %s; %v\n", CurrentFunctionName(4), id, o)
	}

	//dereference underlying entity and store a copy as we don't want it being changed behind our back.
	v := reflect.Indirect(reflect.ValueOf(o))
	mt.Table[id] = v.Interface()

	return nil
}

func (mt *MemTable) Get(id string) (interface{}, error) {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	if mt.Debug {
		log.Printf("%s: %v\n", CurrentFunctionName(4), id)
	}

	r, e := mt.Table[id]

	if !e {
		return nil, NewNotFoundError()
	}

	return r, nil
}

func (mt *MemTable) Update(id string, o interface{}) error {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	if mt.Debug {
		log.Printf("%s: %v; %v\n", CurrentFunctionName(4), id, o)
	}

	if _, e := mt.Table[id]; !e {
		return NewNotFoundError()
	}

	//dereference underlying entity and store a copy as we don't want it being changed behind our back.
	v := reflect.Indirect(reflect.ValueOf(o))
	mt.Table[id] = v.Interface()

	return nil
}

//TODO:Stan date is a project specific item. Should really move it out of here and into the various other projects.
func (mt *MemTable) GetAll(date time.Time, limit int, entity interface{}) (interface{}, error) {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	if mt.Debug {
		log.Printf("%s: time: %s, limit: %d\n", CurrentFunctionName(4), date, limit)
	}

	//Create slice
	ev := reflect.Indirect(reflect.ValueOf(entity))
	slice := reflect.MakeSlice(reflect.SliceOf(ev.Type()), 0, len(mt.Table))

	//Create pointer to slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)

	i := 0
	for _, v := range mt.Table {
		value := reflect.Indirect(reflect.ValueOf(v))
		if !date.IsZero() {

			//Create ptr to value as all stored entities should be values not ptrs. If they weren't we made them values with reflect.Indirect above.
			ptr := reflect.New(reflect.TypeOf(v))
			temp := ptr.Elem()
			temp.Set(value)

			dv := temp.Addr().MethodByName("Date") //TODO:Stan Date still leaks the cassandra partition information in. Look at resolving this when I have nothing better to do.
			if dv.IsValid() {
				fn := dv.Call([]reflect.Value{})
				t := fn[0].Interface().(time.Time)
				if !t.Equal(date) {
					continue
				}
			}
		}

		x.Elem().SetLen(i + 1)
		x.Elem().Index(i).Set(value)
		i++
	}

	return x.Elem().Interface(), nil
}

func (mt *MemTable) All() interface{} {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	return StringMapToSlice(mt.Table)
}

func (mt MemTable) Dump() {
	for k, v := range mt.Table {
		fmt.Printf("\n%s: %v\n", k, v)
	}
}

func NewMemTable() *MemTable {
	return &MemTable{Table: make(map[string]interface{}), mutex: &sync.Mutex{}}
}

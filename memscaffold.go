package gowork

import (
	"log"
	"reflect"
	"sync"
	"errors"
)

type MemTable struct {
	Table memMap
	mutex *sync.Mutex
}

type memMap map[string]interface{}

func (mt *MemTable) Create(id string, o interface{}) {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	mt.Table[id] = o
}

func (mt *MemTable) Get(id string) (interface{}, error) {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	r, e := mt.Table[id]
	if !e {
		return nil, errors.New("not found")
	}

	return r, nil
}

func (mt *MemTable) Update(id string, o interface{}) error {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	old, e := mt.Table[id]
	if !e {
		return errors.New("not found")
	}

	oldValue := reflect.Indirect(reflect.ValueOf(old))
	oldVersion := oldValue.FieldByName("Version")
	if oldVersion.IsValid() {
		i := oldVersion.Interface().(int)

		newValue := reflect.Indirect(reflect.ValueOf(o))
		newVersion := newValue.FieldByName("Version")
		j := newVersion.Interface().(int)

		if i == j {
			newVersion.SetInt(int64(i + 1))
			mt.Table[id] = o
			return nil
		}

		return errors.New("Stale Entity. It has been updated in another session! Please reload and try again.")
	}

	//Record doesn't have the Version field. Update record
	mt.Table[id] = o
	return nil
}

func (mt *MemTable) GetAll(date string, limit int, entity interface{}) interface{} {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()

	//Create slice
	ev := reflect.Indirect(reflect.ValueOf(entity))
	slice := reflect.MakeSlice(reflect.SliceOf(ev.Type()), 0, len(mt.Table))

	//Create pointer to slice
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)

	i := 0
	for _, v := range mt.Table {
		value := reflect.Indirect(reflect.ValueOf(v))
		if date != "" {
			dv := value.FieldByName("Date")
			if dv.IsValid() {
				t := dv.Interface().(string)
				if t != date {
					continue
				}
			}
		}

		x.Elem().SetLen(i + 1)
		x.Elem().Index(i).Set(value)
		i++
	}

	return x.Elem().Interface()
}

func (mt MemTable) Dump() {
	log.Println(mt.Table)
}

func NewMemTable() *MemTable {
	return &MemTable{Table: memMap{}, mutex: &sync.Mutex{}}
}

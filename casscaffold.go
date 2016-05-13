package gowork

import (
	"strings"
	"github.com/gocql/gocql"
	"bytes"
	"reflect"
	"log"
	"fmt"
	"strconv"
	"time"
)

const (
	DEFAULT_CONSISTENCY = gocql.LocalQuorum
)

type unmarshaler func(entityMap map[string]interface{}, entity interface{})

type Cassandra struct {
	Session     *gocql.Session
	Debug       bool
	Unmarshaler unmarshaler
	Partition   string //TODO:Stan it is unlikely that all tables will have the same partition interval (i.e. date, hour, etc..). Look at a way to expand/replace this.
}

func (c *Cassandra) BuildInsertStatement(table string, entity interface{}, overrides map[string]interface{}) (string, []interface{}) {

	var qbuf bytes.Buffer
	var qmbuf bytes.Buffer
	qbuf.WriteString("INSERT INTO ")
	qbuf.WriteString(table)
	qbuf.WriteString(" (")
	val := reflect.Indirect(reflect.ValueOf(entity))
	params := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("datastore")
		if tag == "" {
			tag = strings.ToLower(typeField.Name)
		}

		qbuf.WriteString(tag)
		qmbuf.WriteString("?")
		if i + 1 < val.NumField() {
			qbuf.WriteString(",")
			qmbuf.WriteString(",")
		}

		o, ok := overrides[tag]
		if ok {
			params[i] = o
		} else {
			params[i] = valueField.Interface()
		}
		delete(overrides, tag)

		switch params[i].(type) {
		case MonetaryAmount:
			params[i] = params[i].(MonetaryAmount).Dec //TODO:Stan I should really make a datastore tag for monetaryamount now that Currency is stored in Amount.
		}
	}

	//Write the rest of the overrides.
	for k, v := range overrides {
		if len(params) > 0 {
			qbuf.WriteString(",")
			qmbuf.WriteString(",")
		}
		qbuf.WriteString(k)
		qmbuf.WriteString("?")
		params = append(params, v)
	}

	qbuf.WriteString(") ")
	qbuf.WriteString("VALUES (")
	qbuf.Write(qmbuf.Bytes())
	qbuf.WriteString(")")

	qs := qbuf.String()

	if c.Debug {
		log.Printf("[query statement=%q values=%+v]\n", qs, params)
	}

	return qs, params
}

// BuildUpdateStatement builds a Cassandra Update statement based on the provided entity and overrides.
// Right now it only works for entities with a single id field named id or tagged with "datastore=id".
func (c *Cassandra) BuildUpdateStatement(table string, entity interface{}, overrides map[string]interface{}) (string, []interface{}) {

	var qbuf bytes.Buffer
	qbuf.WriteString("UPDATE ")
	qbuf.WriteString(table)
	qbuf.WriteString(" SET ")

	val := reflect.Indirect(reflect.ValueOf(entity))
	params := make([]interface{}, val.NumField())

	var id interface{}

	pindex := 0 //we need to index the param
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("datastore")
		if tag == "" {
			tag = strings.ToLower(typeField.Name)
		}

		//TODO: Stan look at pulling the hardcoded id out and having it set as a configuration parameter
		if tag == "id" {
			id = valueField.Interface()
			continue
		}

		qbuf.WriteString(tag)
		qbuf.WriteString(" = ?")
		if i + 1 < val.NumField() {
			qbuf.WriteString(", ")
		}

		o, ok := overrides[tag]
		if ok {
			params[pindex] = o
		} else {
			params[pindex] = valueField.Interface()
		}
		delete(overrides, tag)

		switch params[pindex].(type) {
		case MonetaryAmount:
			params[pindex] = params[pindex].(MonetaryAmount).Dec
		}

		pindex++;
	}

	params[pindex] = id
	qbuf.WriteString(" WHERE id = ?")

	o, ok := overrides[c.Partition]
	if ok {
		qbuf.WriteString(" AND ")
		qbuf.WriteString(c.Partition)
		qbuf.WriteString(" = ?")
		//qmbuf.WriteString(",?")
		params = append(params, o)
	}

	qs := qbuf.String()

	if c.Debug {
		log.Printf("[query statement=%q values=%+v]\n", qs, params)
	}

	return qs, params
}

//Builds a CQL INSERT INTO statement with the provided information and executes it.
//Overrides is a map keyed on the datastore tag that allows different values to be specified than the value provided in entity.
func (c *Cassandra) Insert(table string, entity interface{}, overrides map[string]interface{}, consistency gocql.Consistency) error {

	qs, params := c.BuildInsertStatement(table, entity, overrides)

	return c.Session.Bind(qs, func(q *gocql.QueryInfo) ([]interface{}, error) {
		return params, nil
	}).Consistency(consistency).Exec()
}

//TODO: Modify this to work with CAS tables.
//Builds a CQL UPDATE statement with the provided information and executes it.
//Overrides is a map keyed on the datastore tag that allows different values to be specified than the value provided in entity.
func (c *Cassandra) Update(table string, entity interface{}, overrides map[string]interface{}, consistency gocql.Consistency) error {

	qs, params := c.BuildUpdateStatement(table, entity, overrides)

	return c.Session.Bind(qs, func(q *gocql.QueryInfo) ([]interface{}, error) {
		return params, nil
	}).Consistency(consistency).Exec()
}

// GetById is a utility function used to simplify loading an entity from the datastore by Id. Partition Time is optional and is only used by tables that have a partition string partition key.
func (c *Cassandra) GetById(table string, id string, partition time.Time, entity interface{}, consistency gocql.Consistency) (err error) {

	var params []interface{}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table))

	if !partition.IsZero() {
		buf.WriteString(" and " + c.Partition + " = ?")
		params = []interface{}{id, partition}
	} else {
		params = []interface{}{id}
	}

	if c.Debug {
		log.Printf("[query statement=%q values=%+v]\n", buf.String(), params)
	}

	result := make(map[string]interface{})
	if err = c.Session.Bind(buf.String(), func(q *gocql.QueryInfo) ([]interface{}, error) {
		return params, nil
	}).Consistency(consistency).MapScan(result); err != nil {
		if err.Error() == "not found" {
			//gocql uses string messages to differentiate errors.
			return NewNotFoundError()
		}
		return
	}

	c.Unmarshaler(result, entity)

	return
}

func (c *Cassandra) GetAll(table string, limit int, partition time.Time, entity interface{}, consistency gocql.Consistency) (entities interface{}, err error) {

	var params []interface{}

	var buf bytes.Buffer
	buf.WriteString("SELECT * FROM ")
	buf.WriteString(table)

	if !partition.IsZero() {
		buf.WriteString(" WHERE " + c.Partition + " = ? ")
		params = []interface{}{partition}
	}

	if limit > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(limit))
	}

	if c.Debug {
		log.Printf("[query statement=%q values=%+v]\n", buf.String(), params)
	}

	iter := c.Session.Bind(buf.String(), func(q *gocql.QueryInfo) ([]interface{}, error) {
		return params, nil
	}).Consistency(consistency).Iter()

	result, err := iter.SliceMap()
	if err != nil {
		return
	}

	//Build a slice of the required type and size
	val := reflect.Indirect(reflect.ValueOf(entity))
	size := len(result)
	slice := reflect.MakeSlice(reflect.SliceOf(val.Type()), size, size)

	for i, e := range result {
		j := reflect.New(val.Type()) //new pointer to a zero'd entity of the correct type.
		c.Unmarshaler(e, j.Interface())
		slice.Index(i).Set(j.Elem()) //add dereferenced entity to slice.
	}

	entities = slice.Interface()
	err = iter.Close()
	return
}

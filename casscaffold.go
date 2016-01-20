package gowork

import (
	"strings"
	"github.com/gocql/gocql"
	"bytes"
	"reflect"
	"log"
	"fmt"
	"strconv"
)

const (
	DEFAULT_CONSISTENCY = gocql.LocalQuorum
)

type unmarshaler func(entityMap map[string]interface{}, entity interface{})

type Cassandra struct {
	Session     *gocql.Session
	Debug       bool
	Unmarshaler unmarshaler
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

		switch params[i].(type) {
		case MonetaryAmount:
			params[i] = MarshalMonetaryAmount(params[i].(MonetaryAmount))
		}
	}

	qbuf.WriteString(") ")
	qbuf.WriteString("VALUES (")
	qbuf.Write(qmbuf.Bytes())
	qbuf.WriteString(")")

	qs := qbuf.String()

	if c.Debug {
		log.Printf("\nStatement:  %s\nParameters: %s\n\n", qs, params)
	}

	return qs, params
}

//Builds a CQL INSERT INTO statement with the provided information and executes it.
//Overrides is a map keyed on the datastore tag that allows different values to be specified than the value provided in entity.
func (c *Cassandra) Insert(table string, entity interface{}, overrides map[string]interface{}, consistency gocql.Consistency) (err error) {

	qs, params := c.BuildInsertStatement(table, entity, overrides)

	err = c.Session.Bind(qs, func(q *gocql.QueryInfo) ([]interface{}, error) {
		return params, nil
	}).Consistency(consistency).Exec()

	return
}

// GetById is a utility function used to simplify loading an entity from the datastore by Id. Date is optional and is only used by tables that have a date string partition key.
func (c *Cassandra) GetById(table string, id string, date string, entity interface{}, consistency gocql.Consistency) (err error) {

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table))
	if date != "" {
		buf.WriteString(" and date = ?")
	}

	result := make(map[string]interface{})
	if err = c.Session.Bind(buf.String(), func(q *gocql.QueryInfo) ([]interface{}, error) {
		if date != "" {
			return []interface{}{id, date}, nil
		} else {
			return []interface{}{id}, nil
		}
		return nil, nil
	}).Consistency(consistency).MapScan(result); err != nil {
		if err.Error() == "not found" { //gocql uses string messages to differentiate errors.
			return NewNotFoundError()
		}
		return
	}

	c.Unmarshaler(result, entity)

	return
}

func (c *Cassandra) GetAll(table string, limit int, date string, entity interface{}, consistency gocql.Consistency) (entities interface{}, err error) {

	var buf bytes.Buffer
	buf.WriteString("SELECT * FROM ")
	buf.WriteString(table)

	if date != "" {
		buf.WriteString(" WHERE date = ? ")
	}

	if limit > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(limit))
	}

	iter := c.Session.Bind(buf.String(), func(q *gocql.QueryInfo) ([]interface{}, error) {
		if date != "" {
			return []interface{}{date}, nil
		}
		return nil, nil
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


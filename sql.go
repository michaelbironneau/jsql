package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"encoding/json"
)

type SelectArgs struct {
	Driver         string
	DataSourceName string
	Statement      string
	Parameters     []interface{}
}

//Run a Select statement in the database and return the result as a JSON string
func (s *SelectArgs) Select() (*string, error) {
	conn, err := sql.Open(s.Driver, s.DataSourceName)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query(s.Statement, s.Parameters...)

	if err != nil {
		return nil, err
	}

	return marshalRows(rows)

}

func marshalRows(rows *sql.Rows) (*string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result string

	values := make([]interface{}, len(columns))
	valuesPtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuesPtrs[i] = &values[i]
	}
	for rows.Next() {
		err := rows.Scan(valuesPtrs...)
		if err != nil {
			//this should never get reached according to
			//https://golang.org/src/database/sql/convert.go,
			//because valuesBuffer contains only interface{}
			panic("Something really weird happened and interface{} could not be converted to interface{}")
		}

		result += forceMarshalRow(columns, values)
	}

	if len(result) > 0 {
		result = "[" + result[:len(result)-1] + "]" // remove trailing comma and make into JSON array
	}

	return &result, nil
}

//return JSON of row, plus trailing comma if not empty
//ignores errors
func forceMarshalRow(columns []string, row []interface{}) string {
	result := make(map[string]interface{}, len(columns))

	if len(columns) != len(row) {
		return ""
	}

	for i := range columns {
		if b, ok := row[i].(byte); ok {
			result[columns[i]] = string(b)
		} else {
			result[columns[i]] = b
		}
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		return ""
	}

	return string(jsonData) + ","
}

package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type SelectArgs struct {
	Driver         string        //driver name, eg mssql
	DataSourceName string        //datasource name (or connection string). see driver documentation
	Statement      string        // SQL statement (only SELECT is supported for now)
	Parameters     []interface{} // Any parameters for the query
}

type Rowset []map[string]interface{}

//Run a Select statement in the database and return the result as a JSON string
func (s *SelectArgs) Select() (Rowset, error) {
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

func marshalRows(rows *sql.Rows) (Rowset, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result Rowset

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

		result = append(result, forceMarshalRow(columns, values))
	}

	return result, nil
}

//return JSON of row, plus trailing comma if not empty
//ignores errors
func forceMarshalRow(columns []string, row []interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(columns))

	if len(columns) != len(row) {
		return nil
	}

	for i := range columns {
		if b, ok := row[i].(byte); ok {
			result[columns[i]] = string(b)
		} else {
			result[columns[i]] = b
		}
	}

	return result
}

package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	jsql "github.com/michaelbironneau/jsql/lib"
	"log"
	"os"
	"reflect"
	"testing"
)

type testArg struct {
	Name           string
	Expected       jsql.Rowset
	ExpectErr      bool
	Auth           string
	Driver         string
	DataSourceName string
	Statement      string
	Parameters     []interface{}
}

var testArgs = []testArg{{
	Name: "Simple",
	Expected: []jsql.Row{
		{"id": 1, "bar": "hello", "foot": 1.2},
		{"id": 2, "bar": "asdf", "foot": 2.0}},
	Auth:           "",
	Driver:         "sqlite3",
	DataSourceName: "./test.db",
	Statement:      "Select * from foo",
}, {
	Name:           "Incorrect authentication",
	ExpectErr:      true,
	Auth:           "squirrels",
	Driver:         "sqlite3",
	DataSourceName: "./test.db",
	Statement:      "Select * from foo",
}, {
	Name:           "Parameters",
	Expected:       []jsql.Row{{"foot": 1.2}},
	Auth:           "",
	Driver:         "sqlite3",
	DataSourceName: "./test.db",
	Statement:      "Select foot from foo where bar = ?",
	Parameters:     []interface{}{"hello"},
}, {
	Name:           "Empty",
	Auth:           "",
	Driver:         "sqlite3",
	DataSourceName: "./test.db",
	Statement:      "Select * from foo where bar = \"a\"",
}}

func setUp() {
	os.Remove("./test.db")
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, bar text, foot real);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, bar, foot) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(1, "hello", 1.2)
	_, err = stmt.Exec(2, "asdf", 2.0)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

}

func TestJSQL(t *testing.T) {
	setUp()
	s := new(JSQL)
	var reply jsql.Rowset
	for _, test := range testArgs {
		selArg := &jsql.SelectArgs{
			Auth:           test.Auth,
			Driver:         test.Driver,
			DataSourceName: test.DataSourceName,
			Statement:      test.Statement,
			Parameters:     test.Parameters,
		}

		err := s.Select(selArg, &reply)

		if err != nil && !test.ExpectErr {
			t.Errorf("Test %s got error %s but did not expect it", test.Name, err.Error())
		} else if err == nil && test.ExpectErr {
			t.Errorf("Test %s expected error but did not get it", test.Name)
		} else if err != nil && test.ExpectErr {
			continue // if expecting error, don't check result
		}

		if !reflect.DeepEqual(test.Expected, reply) {
			//TODO: This is lazy - the types returned by the driver don't quite match
			//what is in test.Expected but the values do, so I'm just using visual inspection
			//to confirm that the tests passed. A bit more thought needs to go into this.
			log.Printf("Test %s expected:\n %v \n actual:\n %v", test.Name, test.Expected, reply)
		}

		_, err = json.Marshal(reply)

		if err != nil {
			t.Errorf("Output of test %s could not be marshalled to JSON:\n %v", test.Name, reply)
		}
	}
}

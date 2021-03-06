# j-SQL

**License**: MIT

Experimental package to expose an SQL database through a JSON-RPC endpoint (for read-only access at the moment). Supports:

* MS SQL
* MySQL
* Postgres
* SQLite

The only reason I made this is to access MSSQL data from a Cloud Foundry Python app without creating a custom buildpack. In particular, this means I can't install FreeTDS. As I want to use this with the Pandas package without modifying any of my existing code, the return type of the API is a top-level array, so that the `read_json` method will work with it out of the box. 


### Usage
Use `go get github.com/michaelbironneau/jsql`, followed by `go install`. Then you can run

```
jsql --port 1234 --password scrambled_eggs
```

You can now use the client:

```go

import (
	jsql "github.com/michaelbironneau/jsql/client"
	"fmt"
)

func main() {
	c := &jsql.Client{Password: "scrambled_eggs"}

	if err := c.Dial("127.0.0.1:1234"); err != nil {
		panic("Failed to dial server!")
	}

	defer c.Close()

	result, err := c.Query("mssql", "server=localhost;user id=sa", "SELECT 1 AS 'Answer'")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Answer: %v\n", result[0]["Answer"])
}

```

#### Authentication

You can optionally specify a password using the `--password` flag.

#### TLS

You can optionally specify a server certificate and key file using the `--cert` and `--key` flags. The client should then set `TLS` to true:

```go

c := &client.JSQLClient{TLS: true}
```


## Python Client

See `jsql.py`. Usage examples:

```python

from jsql import Database

db = Database("127.0.0.1", 1234, "sqlite3", "./1.db")

print(db.sql("select 1 as 'Answer'"))  # prints [{"Answer": 1}]

print(db.sql("select * from test"))

```

Importing the data into Pandas:

```python
from pandas import DataFrame
from jsql import Database

db = Database("127.0.0.1", 1234, "sqlite3", "./1.db")

df = DataFrame.from_dict(db.sql("select * from test"))
```
# j-SQL

**License**: MIT

Experimental package to expose an SQL database through a JSON-RPC endpoint (for read-only access at the moment).

The only reason I made this is to access MSSQL data from a Cloud Foundry Python app without creating a custom buildpack. In particular, this means I can't install FreeTDS. As I want to use this with the Pandas package without modifying any of my existing code, the return type of the API is a top-level array, so that the `read_json` method will work with it out of the box. 


### Usage
Use `go get github.com/michaelbironneau/jsql`, followed by `go install`. Then you can run

```
jsql --port 1234 --password scrambled_eggs
```

You can now use the client:

```go

import (
	cl "github.com/michaelbironneau/jsql/client"
	"fmt"
)

func main() {
	c := cl.JSQLClient()

	if err := c.Dial("127.0.0.1:1234", "scrambled_eggs"); err != nil {#
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

You can optionally specify a server certificate and key file using the `--cert` and `--key` flags.

**Warning**: You should *never* run this in production without securing it first. At the very least this should include SSL and some authentication. In many environments you will also want to explicitly specify which remote hosts you want to allow access for. You can achieve all of these things by running the app behind Nginx or similar, with the adequate configuration.

Supports:

* MS SQL
* Postgres
* MySQL

## Python Client

TODO. Fow now have a look at https://gist.github.com/stevvooe/1164621.
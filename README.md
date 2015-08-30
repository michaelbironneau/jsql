# j-SQL

**License**: MIT

Experimental package to expose an SQL database through a JSON-RPC endpoint (for read-only access at the moment).

The only reason I made this is to access MSSQL data from a Cloud Foundry Python app without creating a custom buildpack. In particular, this means I can't install FreeTDS. As I want to use this with the Pandas package without modifying any of my existing code, the return type of the API is a top-level array, so that the `read_json` method will work with it out of the box. 


### Usage
Use `go get github.com/michaelbironneau/jsql`, followed by `go install`. Then you can run

```
jsql --port 1234
```

#### Authentication

You can optionally specify a password using the `--password` flag. 

#### TLS

You can optionally specify a server certificate and key file using the `--cert` and `--key` flags.

You now have a j-SQL daemon listening on port 1234:

```go

// client.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

type SelectArgs struct {
	Driver         string        //driver name, eg mssql
	DataSourceName string        //datasource name (or connection string). see driver documentation
	Statement      string        // SQL statement (only SELECT is supported for now)
	Parameters     []interface{} // Any parameters for the query
}

type Reply []map[string]interface{}

func main() {

	client, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := &Args{"mssql", "server=localhost;user id=sa", "SELECT 1 as 'Answer'"}
	var reply Reply
	c := jsonrpc.NewClient(client)
	err = c.Call("JSQL.Select", args, &reply)
	if err != nil {
		log.Fatal("error:", err)
	}

	fmt.Printf("Result: %v\n", reply[0]['Answer'])
}

```


**Warning**: You should *never* run this in production without securing it first. At the very least this should include SSL and some authentication. In many environments you will also want to explicitly specify which remote hosts you want to allow access for. You can achieve all of these things by running the app behind Nginx or similar, with the adequate configuration.

Supports:

* MS SQL
* Postgres
* MySQL

## Python Client

TODO. Fow now have a look at https://gist.github.com/stevvooe/1164621.
# j-SQL

Expose an SQL database through a JSON-RPC endpoint (for read-only access at the moment).

**License**: MIT

**Disclaimer**: You should *never* run this in production without securing it first. At the very least this should include SSL and some authentication. In many environments you will also want to explicitly specify which remote hosts you want to allow access for. You can achieve all of these things by running the app behind Nginx or similar, with the adequate configuration.

Supports:

* MS SQL
* Postgres
* MySQL

The only reason I made this is to access MSSQL data from a Cloud Foundry Python app without creating a custom buildpack. In particular, this means I can't install FreeTDS. As I want to use this with the Pandas package without modifying any of my existing code, the return type of the API is a top-level array, so that the `read_json` method will work with it out of the box. 

## Python Client

TODO. Fow now have a look at https://gist.github.com/stevvooe/1164621.


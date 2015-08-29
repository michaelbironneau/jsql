# j-SQL

Expose an SQL database through a JSON-RPC endpoint (for read-only access at the moment).

The reason I'm doing this is to access an MS SQL database from a Cloud Foundry deployment without writing a custom buildpack. In particular, this means that I can't install FreeTDS or any other driver.




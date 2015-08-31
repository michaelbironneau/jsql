"""
A minimal implementation of a JSON-RPC client
for an SQL database
--------------------------------

@author: michael.bironneau@openenergi.com
License: MIT

Sample usage:

```python

from jsql import Database

db = Database("127.0.0.1", 1234, "sqlite3", "./1.db")

print(db.sql("select 1 as 'Answer'"))  # prints [{"Answer": 1}]

print(db.sql("select * from test"))

```

Example usage (with Pandas):

```python
from pandas import DataFrame
from jsql import Database

db = Database("127.0.0.1", 1234, "sqlite3", "./1.db")

df = DataFrame.from_dict(db.sql("select * from test"))
```

"""


import json
import socket
import ssl


class Database(object):
    """Creates a jSQL client to access a SQL database"""
    def __init__(self, host, port, driver, connection_string, password="", use_ssl=False):
        """
            Initialize the connection to the RPC server.

        * `addr`: Remote address (eg. '127.0.0.1:1234')
        * `driver`: Name of driver (one of 'mssql', 'mysql', 'pg', or 'sqlite3')
        * `connection_string`: Connection string (DataSourceName in the Go code)
        * `password`: Password to access the RPC server (NOT the database user password)
        * `use_ssl`: Whether to use SSL

        """
        self._socket = socket.create_connection((host, port))
        if use_ssl:
            self._socket = ssl.wrap_socket(self._socket)

        self._pass = password
        self._driver = driver
        self._connection_string = connection_string
        self._id = 1

    def sql(self, statement, params=[]):
        """
            Run SQL statement (SELECT only). Return an array of rows: a dict where the keys are
            column names and the corresponding values are the row values.

        `statement`: SQL statement (only SELECT supported for now)
        `params`: List of parameters, if the statement is parametrized
        """
        select_args = {
            "auth": self._pass,
            "driver": self._driver,
            "datasource_name": self._connection_string,
            "statement": statement,
            "params": params
        }
        msg = self._make_request(select_args)
        self._socket.sendall(json.dumps(msg))
        m_id = self._id
        self._id += 1
        resp = self.recv()

        json_resp = json.loads(resp)

        if json_resp['id'] != m_id:
            raise Exception("Response id doesn't match request id")

        if json_resp['error'] is not None:
            raise Exception(json_resp['error'])

        return json_resp['result']

    def _make_request(self, inner_request):
        return dict(id=self._id,
                    params=[inner_request],
                    method="JSQL.Select")

    def recv(self):
        data = self._socket.recv(4096)
        if end_marker(data):
            return data
        while True:
            d = self._socket.recv(4096)
            if len(d) == 0:
                break
            data += d
            if end_marker(data):
                break
        return data

    def __del__(self):
        if self._socket:
            try:
                self._socket.close()
            except:
                pass


def end_marker(data):
    """Go always outputs a line feed at the end of output, and just to be sure
    we check if '}' is the next to last character as expected. This seems somewhat
    brittle but it works better in practice than using short timeouts, since some
    database queries will break that very easily."""
    if ord(data[-1]) == 10 and data[-2] == '}':
        return True

if __name__ == "__main__":
    db = Database('127.0.0.1', 1234, "sqlite3", "./1.db")
    print(db.sql("select * from test"))

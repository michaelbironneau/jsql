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
import time
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
        resp = recv_timeout(self._socket)

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


def recv_timeout(the_socket,timeout=2):
    #make socket non blocking
    the_socket.setblocking(0)
     
    #total data partwise in an array
    total_data=[];
    data='';
     
    #beginning time
    begin=time.time()
    while 1:
        #if you got some data, then break after timeout
        if total_data and time.time()-begin > timeout:
            break
         
        #if you got no data at all, wait a little longer, twice the timeout
        elif time.time()-begin > timeout*2:
            break
         
        #recv something
        try:
            data = the_socket.recv(8192)
            if data:
                total_data.append(data)
                #change the beginning time for measurement
                begin=time.time()
            else:
                #sleep for sometime to indicate a gap
                time.sleep(0.1)
        except:
            pass
     
    #join all parts to make final string
    return ''.join(total_data)

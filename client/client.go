package client

import (
	"crypto/tls"
	"github.com/michaelbironneau/jsql"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type JSQLClient struct {
	Addr       string //server address, eg. '127.0.0.1:1234'
	Password   string //password for server
	TLS        bool   //whether the server uses TLS
	SkipVerify bool   //whether to skip TLS certificate verification
	client     *rpc.Client
	conn       net.Conn
}

func (j *JSQLClient) Dial(address string, password string) error {
	j.Addr = address
	j.Password = password

	var (
		tlsConfig *tls.Config
		err       error
		conn      net.Conn
	)

	if j.TLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: j.SkipVerify,
		}
		conn, err = tls.Dial("tcp", j.Addr, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", j.Addr)
	}

	if err != nil {
		return err
	}
	j.client = jsonrpc.NewClient(conn)
	return nil
}

func (j *JSQLClient) Close() error {
	return j.conn.Close()
}

func (j *JSQLClient) Query(driver string, dataSourceName string, statement string, params ...interface{}) (jsql.Rowset, error) {
	args := &jsql.SelectArgs{
		Auth:           j.Password,
		Driver:         driver,
		DataSourceName: dataSourceName,
		Statement:      statement,
		Parameters:     params,
	}

	var reply jsql.Rowset

	err := j.client.Call("JSQL.Select", args, &reply)

	return reply, err

}

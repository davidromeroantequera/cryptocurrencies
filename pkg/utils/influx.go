package utils

import (
	"github.com/influxdata/influxdb/client/v2"
)

const serverAddress = "http://localhost:8086"

func QueryDB(cmd string, dbName string) (res []client.Result, err error) {
	// This function is not meant for intensive DB operations, as it opens and closes the client
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: serverAddress})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	q := client.Query{
		Command:  cmd,
		Database: dbName,
	}

	if response, err := c.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func CreateDatabase(dbName string) error {
	const (
		Q = "\""
		queryBase = "CREATE DATABASE "
	)

	_, err := QueryDB(queryBase + Q + dbName + Q, dbName)

	return err
}

func DropDatabase(dbName string) error {
	const(
		Q = "\""
		queryBase = "DROP DATABASE "
	)

	_, err := QueryDB(queryBase + Q + dbName + Q, dbName)

	return err
}

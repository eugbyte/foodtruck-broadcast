package dbconfig

import (
	"fmt"
	"strings"
	"testing"

	assrt "github.com/blend/go-sdk/assert"
)

func TestConnString(t *testing.T) {
	assert := assrt.New(t)

	config := Config{
		Host:     "abc",
		Password: "123",
	}

	connstr := ConnString(Postgres, config)
	fmt.Println(connstr)

	assert.True(strings.Contains(connstr, "host=abc"))
	assert.True(strings.Contains(connstr, "password=123"))
	assert.False(strings.Contains(connstr, "dbname="), "if config value is empty, then the connection string should completely omit the 'key=pair' value")
}

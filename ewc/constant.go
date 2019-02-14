package ewc

import (
	"runtime"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	if runtime.GOOS == "android" {
		baseUrl = "http://10.0.2.2:9999"
	} else {
		baseUrl = "http://127.0.0.1:9999"
	}
}

const (
	dbName   = "Eisenwort.sqlite"
	chanSize = 10
	driver   = "sqlite3"
)

var baseUrl = ""
var userId = int64(0)
var dbPath = ""
var userIdHeader = ""
var connectionString = ""

var jiter = jsoniter.ConfigFastest

package ewc

import (
	"runtime"
)

func init() {
	if runtime.GOOS == "android" {
		baseUrl = "http://10.0.2.2:9999"
	} else {
		baseUrl = "http://127.0.0.1:9999"
	}
}

const (
	dbName   = "ewc.sqlite"
	chanSize = 5
	driver   = "sqlite3"
	IdHeader = "X-Auth-Id"
)

var baseUrl = ""
var userID = int64(0)
var dbPath = ""
var userIDHeader = ""
var connectionString = ""
var alghorinthm = ""
var jwtToken = ""

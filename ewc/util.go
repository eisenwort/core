package ewc

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Util struct {
}

func NewUtil() *Util {
	item := new(Util)

	return item
}

func (u *Util) Setup(data *SetupData) {
	if data.DbDriver == "sqlite3" {
		if data.DbPath != "" && !strings.HasSuffix(data.DbPath, "/") {
			data.DbPath += "/"
		}
		if data.ConnectionString == "" {
			connectionString = data.DbPath + dbName
		} else {
			connectionString = data.ConnectionString
		}
	} else {
		connectionString = data.ConnectionString
	}

	driver = data.DbDriver
	pageLimit = data.PageLimit
}

func (u *Util) CloseApp() {
	db := getDb()

	if err := db.Close(); err != nil {
		log.Println("close db error:", err)
	}

	httpClient.CloseIdleConnections()

	// encrypt DB

}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func getClaims(token string) JwtClaims {
	claims := JwtClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		log.Println("get claims error:", err)
	}

	return claims
}

func serialize(item interface{}) string {
	return string(serializeByte(item))
}

func serializeByte(item interface{}) []byte {
	data, err := json.Marshal(item)

	if err != nil {
		log.Println("serialize object error:", err)
		return []byte("{}")
	}

	return data
}

func deserialize(data string, item interface{}) {
	if err := json.Unmarshal([]byte(data), item); err != nil {
		log.Println("deserialize object error:", err)
	}
}

func getBodyString(body io.Reader) string {
	data, err := ioutil.ReadAll(body)

	if err != nil {
		log.Println("get http body error:", err)
		return "{}"
	}

	return string(data)
}

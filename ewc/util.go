package ewc

import (
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Setup(data SetupData) {
	if data.DbDriver == "sqlite3" {
		if !strings.HasSuffix(data.DbPath, "/") {
			data.DbPath += "/"
		}
		connectionString = data.DbPath + dbName
	} else {
		connectionString = data.ConnectionString
	}

	driver = data.DbDriver
	pageLimit = data.PageLimit
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

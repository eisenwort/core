package ewc

import (
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Setup(data SetupData) {
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

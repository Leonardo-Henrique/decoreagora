package repositories

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/config"
	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var strConn string = "%s:%s@tcp(%s:%s)/%s?&parseTime=True&loc=UTC&collation=utf8mb4_0900_ai_ci"

func ConnectToDatabase(dsn string) (*sql.DB, error) {

	if config.C.IsProd == "true" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile("./certs/ca-certificate-prod-db.crt")
		if err != nil {
			return nil, err
		}

		ok := rootCertPool.AppendCertsFromPEM(pem)
		if !ok {
			return nil, errors.New("failed when trying to append db cert from pem")
		}

		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})

		strConn += "&tls=custom"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successufly connected to database")

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(30)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil

}

func MakeDSNString(dsn models.DSN) string {
	return fmt.Sprintf(strConn,
		dsn.User,
		dsn.Pass,
		dsn.Host,
		dsn.Port,
		dsn.DatabaseName,
	)
}

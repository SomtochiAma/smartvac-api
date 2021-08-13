package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"os"
	"strings"
)

var DB *gorm.DB

type cred struct {
	username string
	password string
	host     string
	dbname   string
	sslMode  string
}

func Init() (*gorm.DB, error) {
	cred, err := getCreds()
	if err != nil {
		logrus.Errorf("error getting creds: %s", err)
		return nil, err
	}
	fmt.Println(cred.sslMode)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s",
		cred.host,
		cred.username,
		cred.password,
		cred.dbname, cred.sslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db

	if os.Getenv("DROP_TABLE") == "true" {
		logrus.Info("dropping existing tables")
		DB.Migrator().DropTable(&Reading{})
		DB.Migrator().DropTable(&User{})
		DB.Migrator().DropTable(&Payment{})
	}

	DB.AutoMigrate(&Reading{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Payment{})

	//payments := Payment{
	//	UserID: 1,
	//	Amount: 1000,
	//	Units:  250,
	//	Day:    time.Now(),
	//}
	//if err := db.Create(&payments); err != nil {
	//	logrus.Errorf("unable to create base payment: %s", err)
	//}

	return db, nil
}

func getCreds() (*cred, error) {
	var username, password, host, dbname string
	sslMode := os.Getenv("SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		logrus.Info("getting credentials from DATABASE_URL env_var")
		parsedURL, err := url.Parse(dbURL)
		if err != nil {
			return nil, err
		}

		username = parsedURL.User.Username()
		host = strings.Split(parsedURL.Host, ":")[0]
		dbname = parsedURL.Path
		if strings.Contains(parsedURL.Path, "/") {
			dbname = dbname[1:]
		}
		password, _ = parsedURL.User.Password()
	} else {
		logrus.Info("getting credentials from individual env_var")
		requiredVar := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_NAME", "SSL_MODE"}
		for _, val := range requiredVar {
			if _, exist := os.LookupEnv(val); !exist {
				return nil, errors.New(fmt.Sprintf("env variable %s not set", val))
			}
		}

		username, password, host, dbname = os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	}
	return &cred{
		username: username,
		password: password,
		host:     host,
		dbname:   dbname,
		sslMode:  sslMode,
	}, nil
}

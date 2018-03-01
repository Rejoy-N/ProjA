package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type Aduser struct {
	Auid     string            `valid:"required,uuidv4"`
	Username string            `valid:"required,alpha"`
	Password string            `valid:"required"`
	Fname    string            `valid:"required,alpha"`
	Lname    string            `valid:"required,alpha"`
	Email    string            `valid:"required,email"`
	Errors   map[string]string `valid:"-"`
}

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "rejoy"
	DB_PASSWORD = "rejoy"
	DB_NAME     = "persproj"
)

func EncryptPassword(password string) string {
	pass := []byte(password)
	hashpw, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashpw)
}

func saveAdminuser(a *Aduser) error {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database..")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Exec("Create table if not exists adusers(auid varchar(50) primary key, firstname varchar(50) not null, lastname varchar(50) not null, username varchar(50) not null unique, email varchar(50) not null unique, password varchar(100) not null)")
	fmt.Println("Preparing to Insert data into the database....")
	db.Exec("Insert into adusers (auid, firstname, lastname, username, email, password) values ($1, $2, $3, $4, $5, $6);", a.Auid, a.Fname, a.Lname, a.Username, a.Email, a.Password)
	if err != nil {
		panic(err)
	}
	fmt.Println("Saved Admin user data successfully")
	return err
}

func Auid() string {
	id := uuid.NewV4()
	return id.String()
}

func AuserExists(auser *Aduser) (bool, string) {
	fmt.Println("Verifying credentials....")
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database..")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var au, pw string
	fmt.Println("Verifying user information....")
	fmt.Println(auser.Username)
	q, err := db.Query(" Select auid, password from adusers where  username = '" + auser.Username + "'")
	if err != nil {
		// panic(err)
		fmt.Println("Failed to retrieve records.")
		return false, ""
	}
	for q.Next() {
		q.Scan(&au, &pw)
	}
	fmt.Println(au)
	if au != "" {
		ps := bcrypt.CompareHashAndPassword([]byte(pw), []byte(auser.Password))
		fmt.Println("Credential verification completed...")
		if ps == nil {
			return true, au
		}
		return false, ""
	} else {
		return false, ""
	}
}

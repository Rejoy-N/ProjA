package main

import (
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

func saveFeaCat(feacat string) error {
	fmt.Println("Save Feature Category.. ")
	//	validstr := regexp.MustCompile(`^[a-z]`)
	validstr := regexp.MustCompile("\\w")
	fmt.Println(validstr.MatchString(feacat))
	if validstr.MatchString(feacat) == true {
		fmt.Println("Feature Category Function")
		fcuid := Buid()
		db := connectDb()
		db.Exec("create table if not exists featurecat ( fcuid varchar(50) primary key, fcatname varchar(50) not null unique)")
		fmt.Println("Inserting feature category data...")
		_, err := db.Exec("Insert into featurecat (fcuid, fcatname) values ($1, $2);", fcuid, feacat)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		fmt.Println("Saved feature category data....")
		return err
	}
	return nil
}

func getFeacatData() ([]byte, error) {
	db := connectDb()
	db.Exec("create table if not exists featurecat ( fcuid varchar(50) primary key, fcatname varchar(50) not null unique)")
	rows, err := db.Query("select * from featurecat")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	fmt.Println("Retrieving Feature Category JSON data..")
	feacatjondata, err := getTableData(rows, columns)
	if err != nil {
		panic(err)
	}
	return feacatjondata, nil
}

func delFeaCat() error {
	return nil
}

func updateFeaCat(data map[string]string) error {
	fmt.Println("Updating Feature Category data..")
	fmt.Println(data)
	validstr := regexp.MustCompile("\\w")
	fmt.Println(validstr.MatchString(data["fcatname"]))
	if validstr.MatchString(data["fcatname"]) == true {
		fmt.Println("Update Feature Category Data row..")
		saveorupdateRowData(data, "featurecat", "fcuid", 1)
	}
	return nil

}

func saveFeature(id string, feature string) error {
	fmt.Println("Save Feature .. ")

	//	validstr := regexp.MustCompile(`^[a-z]`)
	validstr := regexp.MustCompile("\\w")
	fmt.Println(validstr.MatchString(feature))
	if validstr.MatchString(id) == true {
		if validstr.MatchString(feature) == true {
			fmt.Println("Feature details")
			fuid := Buid()
			db := connectDb()
			db.Exec("create table if not exists features ( fuid varchar(50) primary key, fcatid varchar(50) not null, featurename varchar(50) not null unique)")
			fmt.Println("Inserting feature category data...")
			_, err := db.Exec("Insert into features (fuid, fcatid, featurename) values ($1, (select featurecat.fcuid from featurecat where featurecat.fcuid=($2)), $3);", fuid, id, feature)
			if err != nil {
				panic(err)
			}
			defer db.Close()
			fmt.Println("Saved feature category data....")
			return err
		}
	}

	return nil
}

func getFeaturesData() ([]byte, error) {
	fmt.Println("Get features data function...")

	db := connectDb()
	db.Exec("create table if not exists features ( fuid varchar(50) primary key, fcatid varchar(50) not null, featurename varchar(50) not null unique)")

	rows, err := db.Query("(select *from features INNER JOIN featurecat ON features.fcatid = featurecat.fcuid)")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	fmt.Println("Retrieving Feature Category JSON data..")
	featurejsondata, err := getTableData(rows, columns)
	if err != nil {
		panic(err)
	}
	return featurejsondata, nil

}

func updateFeature(data map[string]string) error {
	fmt.Println("Updating Feature Category data..")
	fmt.Println(data)
	validstr := regexp.MustCompile("\\w")
	fmt.Println(validstr.MatchString(data["featurename"]))
	if validstr.MatchString(data["featurename"]) == true {
		fmt.Println("Update Feature Data row..")
		saveorupdateRowData(data, "features", "fuid", 1)
	}
	return nil

}

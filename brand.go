package main

import (
	"fmt"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

/* const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "rejoy"
	DB_PASSWORD = "rejoy"
	DB_NAME     = "persproj"
) */

type Brand struct {
	Buid    string            `valid:"required:uuidv4"`
	Bname   string            `valid:"required:alphanum"`
	Bimage  string            `valid:"required:alphanum"`
	Bstatus string            `valid:"required:alphanum"`
	Errors  map[string]string `valid:"-"`
}

func SaveBranddata(b *Brand) error {
	db := connectDb()
	db.Exec("create table if not exists brands (buid varchar(50) primary key, bname varchar(50) not null, bimage varchar(50) not null, bstatus varchar(50) not null)")
	fmt.Println("Inserting brand data...")
	_, err := db.Exec("Insert into brands (buid, bname, bimage, bstatus) values ($1, $2, $3, 'Inactive');", b.Buid, b.Bname, b.Bimage)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Saved brand data....")
	return err
}

func Buid() string {
	uid := uuid.NewV4()
	return uid.String()
}

func getBranddata() ([]byte, error) {
	db := connectDb()
	db.Exec("create table if not exists brands (buid varchar(50) primary key, bname varchar(50) not null, bimage varchar(50) not null, bstatus varchar(50) not null)")
	db.Exec("create table if not exists delbrands (buid varchar(50) primary key)")
	db.Exec("create table if not exists brandstatus (buid varchar(50) primary key)")
	rows, err := db.Query("select * from brands where not exists ((select delbrands.buid from delbrands where brands.buid = delbrands.buid) UNION (select brandstatus.buid from brandstatus where brands.buid = brandstatus.buid))")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	fmt.Println("Retrieving json object of brand data records...")
	brandjsondata, err := getTableData(rows, columns)
	fmt.Println(string(brandjsondata))
	if err != nil {
		return nil, err
	}
	return brandjsondata, nil

}

func getBrand(id string) ([]byte, error) {
	fmt.Println("fetching brand data")
	cols := []string{"*"}
	a := []string{id}
	d, err := getRowData(cols, "brands", "buid", a)
	if err != nil {
		panic(err)
	}
	return d, nil
}

func deleteBrand(ids []string) error {
	fmt.Println("Deleting the selected brands")
	a := ids
	b, err := qVal("insert", a, "delbrands", "buid", "")
	if err != nil {
		return err
	}
	fmt.Println(b)
	return nil
}

func deactivateBrand(ids []string) error {
	fmt.Println("Deactivating the selected brands")
	a := ids
	b, err := qVal("deactivate", a, "brands", "bstatus", "buid")
	if err != nil {
		return err
	}
	fmt.Println(b)
	return nil
}

func activateBrand(ids []string) error {
	fmt.Println("Activating the selected brands")
	a := ids
	b, err := qVal("activate", a, "brands", "bstatus", "buid")
	if err != nil {
		return err
	}
	fmt.Println(b)
	return nil
}

/*
func updatebrandRow(data map[string]string) error {
	fmt.Println("Updating Brand data..")
	fmt.Println(data)
	saveorupdateRowData(data, "brands", "buid", 2)
	return nil
}
*/

func updatebrandRow(data *Brand) error {
	fmt.Println("Updating brand data..")
	fmt.Println(data)
	var udata = make(map[string]string)
	udata["Buid"] = data.Buid
	udata["Bstatus"] = data.Bstatus
	udata["Bimage"] = data.Bimage
	udata["Bname"] = data.Bname
	fmt.Println("Printing udata..")
	fmt.Println(udata)
	saveorupdateRowData(udata, "brands", "buid", 2)
	return nil

}

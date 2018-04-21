package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

func connectDb() (db *sql.DB) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	fmt.Println("Reading database credentials....")
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
	// defer db.Close()
	return db
}

func getTableData(rs *sql.Rows, col []string) ([]byte, error) {
	//initialize a empty slice of map tableData  with key as string and data of interface type
	tableData := make([]map[string]interface{}, 0)
	// initialize a slice of map data values of length equal to the number of columns
	fmt.Println("Serializing JSON for records retrieved....n")
	for rs.Next() {
		// initialize a slice of map data of length equal to the number of columns
		container := make([]interface{}, len(col))
		//initialize a slice of map data of length equal to the number of columns
		dest := make([]interface{}, len(col))
		for i, _ := range container {
			dest[i] = &container[i]
		}
		rs.Scan(dest...)
		r := make(map[string]interface{})
		for i, colname := range col {
			val := dest[i].(*interface{})
			r[colname] = *val
		}
		tableData = append(tableData, r)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return nil, err
	}
	fmt.Println("Creating JSON Object.....")
	return jsonData, nil
}

func getRowData(colparams []string, tbname string, colname string, list []string) ([]byte, error) {
	fmt.Println("fetching row data..")
	col := ""
	for i, _ := range colparams {
		col += colparams[i]
		if i < len(colparams)-1 {
			col += ", "
		}
	}
	l := ""
	for j, _ := range list {
		l += "'" + list[j] + "'"
		if j < len(list)-1 {
			l += ", "
		}
	}
	l += ");"
	q := "select " + col + " from " + tbname + " where " + colname + " in (" + l
	fmt.Println(q)
	db := connectDb()
	//rows, err := db.Query("select " + colparams + " from " + tbname + " where" + colname + " in (" + l)
	rows, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	fmt.Println("Retrieving json object of row records...")
	rowsjsondata, err := getTableData(rows, columns)
	fmt.Println(string(rowsjsondata))
	if err != nil {
		return nil, err
	}
	return rowsjsondata, nil
}

func qVal(op string, ids []string, tbname string, colname1 string, colname2 string) (string, error) {
	q := ""
	s1 := ""
	s2 := ""
	s3 := ""
	switch op {
	case "insert":
		q += "insert into " + tbname + " " + "(" + colname1 + ")" + " values "
		s1 = "('"
		s2 = "')"
		s3 = ";"
	case "delete":
		q += "delete from " + tbname + " " + "(" + colname1 + ")" + " values "
		s1 = "('"
		s2 = "')"
		s3 = ";"
	case "deactivate":
		q += "update " + tbname + " set " + colname1 + " = 'Inactive' where " + colname2 + " in ("
		s1 = "'"
		s2 = "'"
		s3 = ");"
	case "activate":
		q += "update " + tbname + " set " + colname1 + " = 'Active' where " + colname2 + " in ("
		s1 = "'"
		s2 = "'"
		s3 = ");"
	}
	for i, _ := range ids {
		fmt.Println(ids[i])
		// q += "('" + ids[i] + "')"
		q += s1 + ids[i] + s2
		if i < len(ids)-1 {
			q += ", "
		}
	}
	// q += ";"
	q += s3
	fmt.Println(q)
	db := connectDb()
	_, err := db.Exec(q)
	if err != nil {
		panic(err)
	}
	return op + " operation successful", nil
}

func saveorupdateRowData(colparams map[string]string, tbname string, pkey string, unchangedparam int) error {
	fmt.Println("updating row data..")
	fmt.Println(colparams)

	q := "Insert into " + tbname + " "
	col := "("
	l := "("
	paramcount := 0
	plist := ""
	vlist := ""

	for key, val := range colparams {
		col += key
		paramcount++
		s := strconv.Itoa(paramcount)
		l += "$" + s

		// vlist += `"` + val + `"`
		vlist += "'" + val + "'"

		if paramcount > unchangedparam {
			plist += key + " = EXCLUDED." + key
		}

		if paramcount < len(colparams) {
			col += ", "
			l += ", "
			if plist != "" {
				plist += ", "
			}
			if vlist != "" {
				vlist += ", "
			}
		} else {
			col += ") "
			l += ") "
			plist += " ;"
			if vlist != "" {
				vlist += ""
			}
		}
	}

	// q += col + "values " + l + "ON CONFLICT " + "(" + pkey + ") DO UPDATE SET " + plist
	q += col + "values " + "(" + vlist + ") " + "ON CONFLICT " + "(" + pkey + ") DO UPDATE SET " + plist

	fmt.Println(q)

	db := connectDb()
	/*  _, err := db.Exec("Insert into brands(Buid, Bstatus, Bimage, Bname) values ($1, $2, $3, $4) on Conflict (buid) do UPDATE set Bimage = EXCLUDED.Bimage, Bname = EXCLUDED.Bname ;", "one", "two", "three", "four")
	 */
	_, err := db.Exec(q)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/*
func readimageFile(r *http.Request) string {
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
	fmt.Println("Reading file info...")
	file, header, err := r.FormFile("brandlogo")
	if err != nil {

		// http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	var fname = header.Filename
	fmt.Println("file read successfully...")
}
*/

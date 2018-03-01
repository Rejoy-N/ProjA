package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
)

func hello(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "Hello")
	render(w, "hello", "Hello")
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, "home", "This is the home page")
}

/*
func adminpage(w http.ResponseWriter, r *http.Request) {
	msg := GetMsg(w, r, "message")
	var auser = &Aduser{}
	auser.Errors = make(map[string]string)
	if msg != "" {
		auser.Errors["message"] = msg
		render(w, "adminpage", auser)
	} else {
		auser := &Aduser{}
		render(w, "adminpage", auser)
	}
}
*/

func adminpage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("Creating Get request....")
		auser := &Aduser{}
		auser.Errors = make(map[string]string)
		auser.Errors["username"] = GetMsg(w, r, "username")
		auser.Errors["password"] = GetMsg(w, r, "password")
		auser.Errors["email"] = GetMsg(w, r, "email")
		fmt.Println(auser.Errors["username"])
		fmt.Println(auser.Errors["password"])
		fmt.Println(auser.Errors["email"])
		render(w, "adminpage", auser)
	case "POST":
		fmt.Println("Creating Post request...")

		username := r.FormValue("uname")
		password := r.FormValue("password")
		email := r.FormValue("email")

		result := true
		if email == "" {
			if username == "" {
				SetMsg(w, "username", "Username field cannot be blank")
				result = false
			}
			if password == "" {
				SetMsg(w, "password", "Password field cannot be blank")
				result = false
			}
		}
		if email != "" {
			SetMsg(w, "email", "An email has been sent to your registered email ID to reset your password")
			result = false
		}
		redirect := "/"
		if result == true {
			fmt.Println(result)
			auser := &Aduser{Username: username, Password: password}
			if b, auid := AuserExists(auser); b == true && auid != "" {
				fmt.Println("Creating admin user session....")
				setsession(&Aduser{Auid: auid}, w)
				redirect = "/admin"
			} else {
				fmt.Println("Error")
				SetMsg(w, "username", "Username does not exist")
				redirect = "/"
			}

		}
		http.Redirect(w, r, redirect, 302)
	}
}

func adminlogin(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("uname")
	pass := r.FormValue("password")
	u := &Aduser{Username: name, Password: pass}
	redirect := "/"
	if name != "" && pass != "" {
		fmt.Println("login successful...")
		fmt.Println(u)
	}
	http.Redirect(w, r, redirect, 302)
	// render(w, "adminpage", "this is the admin login")
}

func adminregister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("Creating Get request....")
		auser := &Aduser{}
		auser.Errors = make(map[string]string)
		auser.Errors["lname"] = GetMsg(w, r, "lname")
		auser.Errors["fname"] = GetMsg(w, r, "fname")
		auser.Errors["email"] = GetMsg(w, r, "email")
		auser.Errors["username"] = GetMsg(w, r, "username")
		auser.Errors["password"] = GetMsg(w, r, "password")
		fmt.Println(auser.Errors["lname"])
		fmt.Println(auser.Errors["fname"])
		fmt.Println(auser.Errors["email"])
		fmt.Println(auser.Errors["username"])
		fmt.Println(auser.Errors["password"])
		render(w, "adminregister", auser)
	case "POST":
		fmt.Println("Creating Post request...")
		auser := &Aduser{
			Auid:     Auid(),
			Fname:    r.FormValue("fname"),
			Lname:    r.FormValue("lname"),
			Username: r.FormValue("uname"),
			Password: r.FormValue("password"),
			Email:    r.FormValue("email"),
		}

		fmt.Println(r.FormValue("fname"))
		fmt.Println(r.FormValue("lname"))
		fmt.Println(r.FormValue("uname"))
		fmt.Println(r.FormValue("password"))
		fmt.Println(r.FormValue("email"))

		fmt.Println(auser)

		result, err := govalidator.ValidateStruct(auser)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Lname"); re == true {
				SetMsg(w, "lname", "Please enter a valid last name")
			}
			if re := strings.Contains(e, "Email"); re == true {
				SetMsg(w, "email", "Please enter a email address")
			}
			if re := strings.Contains(e, "Fname"); re == true {
				SetMsg(w, "fname", "Please enter a valid first name")
			}
			if re := strings.Contains(e, "Username"); re == true {
				SetMsg(w, "username", "Please enter a valid username")
			}
			if re := strings.Contains(e, "Password"); re == true {
				SetMsg(w, "password", "Please enter a valid Password")
			}
		}
		if result == true {
			fmt.Println(result)
			auser.Password = EncryptPassword(auser.Password)
			saveAdminuser(auser)
			http.Redirect(w, r, "/adminpage", 302)
			return
		}
		fmt.Println("Not True")
		http.Redirect(w, r, "/adminregister", 302)
	}
}

func admin(w http.ResponseWriter, r *http.Request) {
	render(w, "admin", "This is the admin page")
}

func brand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("Processing GET request for brand data....")
		brd, err := getBranddata()
		if err != nil {
			fmt.Println(err)
		}
		var brdata ([]map[string]interface{})
		err = json.Unmarshal(brd, &brdata)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(brdata)
		render(w, "brand", brdata)
	case "POST":
		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		fmt.Println("Reading file info...")
		file, header, err := r.FormFile("brandlogo")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		var fname = header.Filename
		fmt.Println(fname)
		fmt.Println("file read successfully...")

		fmt.Println(r.FormValue("brandname"))

		b := &Brand{
			Buid:    Buid(),
			Bname:   r.FormValue("brandname"),
			Bimage:  fname,
			Bstatus: "Inactive",
		}

		fmt.Println(b)
		result, err := govalidator.ValidateStruct(b)
		fmt.Println(err)
		fmt.Println(result)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Bname"); re == true {
				SetMsg(w, "pname", "Please enter a valid brand name")
			}
			if re := strings.Contains(e, "Bimage"); re == true {
				SetMsg(w, "pname", "Error encountered in upload of brand image")
			}
		}
		if result == true {
			f, err := os.Create("./assets/brands/" + fname)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}
			defer f.Close()
			fmt.Println("Saving file to location....")
			if _, err := io.Copy(f, file); err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}
			fmt.Println("Saved file image.....")
			SaveBranddata(b)
		}
		http.Redirect(w, r, "/brand", 302)
	}
}

func deletebrand(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete...deactivate...activate brands...")
	a := r.FormValue("selectedbrandids")
	fmt.Println("Printing the brand ids array")
	fmt.Println(a)
	selectedids := make([]string, 0)
	err := json.Unmarshal([]byte(a), &selectedids)
	if err != nil {
		panic(err)
	}
	fmt.Println(selectedids)
	action := r.FormValue("action")
	fmt.Println(action)
	fmt.Println("here describe entire function for brands deletion")
	switch action {
	case "delete":
		fmt.Println("Selected case is delete")
		deleteBrand(selectedids)
	case "deactivate":
		fmt.Println("Selected case is deactivate")
		deactivateBrand(selectedids)
	case "activate":
		fmt.Println("Selected case is activate")
		activateBrand(selectedids)
	}
	http.Redirect(w, r, "/brand", 302)
}

func updatebrand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("Update brands")
		a := r.FormValue("update")
		fmt.Println(a)
		data, _ := getBrand(a)
		var brrowdata ([]map[string]interface{})
		err := json.Unmarshal(data, &brrowdata)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(brrowdata)
		render(w, "updatebrand", brrowdata)

	case "POST":
		fmt.Println("Updating Brand Data - Timestamp Check, Active Products checks, DB Updation of archieved records")
		fmt.Println("fetch column names or store array... for now using temp fix")

		img := r.FormValue("imgname")
		fmt.Println(img)

		r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)
		fmt.Println("Reading file info...")
		file, header, err := r.FormFile("brandlogo")
		if err == nil {
			img = header.Filename
		}

		fmt.Println("file read successfully...")
		/*
			var udata = make(map[string]string)
			udata["Buid"] = r.FormValue("brandid")
			udata["Bstatus"] = r.FormValue("brandstatus")
			udata["Bimage"] = img
			udata["Bname"] = r.FormValue("brandname")
		*/

		udata := &Brand{
			Buid:    r.FormValue("brandid"),
			Bname:   r.FormValue("brandname"),
			Bimage:  img,
			Bstatus: r.FormValue("brandstatus"),
		}

		fmt.Println(udata)
		result, err := govalidator.ValidateStruct(udata)
		fmt.Println(err)
		fmt.Println(result)
		if err != nil {
			e := err.Error()
			if re := strings.Contains(e, "Bname"); re == true {
				SetMsg(w, "pname", "Please enter a valid brand name")
			}
			if re := strings.Contains(e, "Bimage"); re == true {
				SetMsg(w, "pname", "Error encountered in upload of brand image")
			}
		}
		if result == true {
			f, err := os.Create("./assets/brands/" + img)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}
			defer f.Close()
			fmt.Println("Saving file to location....")
			if _, err := io.Copy(f, file); err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}
			fmt.Println("Saved file image.....")
			updatebrandRow(udata)
		}
		http.Redirect(w, r, "/brand", 302)
	}
}

func features(w http.ResponseWriter, r *http.Request) {
	render(w, "features", "This is the Manage Features page")
}

func product(w http.ResponseWriter, r *http.Request) {
	render(w, "product", "This is the Manage Product page")
}

func render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseGlob("view/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	govalidator.SetFieldsRequiredByDefault(true)

	serverMuxA := http.NewServeMux()
	serverMuxA.HandleFunc("/", adminpage)
	serverMuxA.HandleFunc("/adminregister", adminregister)
	serverMuxA.HandleFunc("/admin", admin)
	serverMuxA.HandleFunc("/brand", brand)
	serverMuxA.HandleFunc("/deletebrand", deletebrand)
	serverMuxA.HandleFunc("/updatebrand", updatebrand)
	serverMuxA.HandleFunc("/features", features)
	serverMuxA.HandleFunc("/product", product)

	serverMuxB := http.NewServeMux()
	serverMuxB.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	serverMuxB.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	serverMuxB.HandleFunc("/", hello)
	serverMuxB.HandleFunc("/home", home)

	go func() {
		serverMuxA.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
		serverMuxA.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
		err := http.ListenAndServe(":8093", serverMuxA)
		log.Fatal(err)
	}()
	err := http.ListenAndServe(":8094", serverMuxB)
	log.Fatal(err)
}

package main

import (
	"6Thaer/model"
	"6Thaer/repository"
	"io"
	"log"
	"net"
	"strings"

	"encoding/json"
	"fmt"

	"io/ioutil"

	"net/http"
	"os"
)

func main() {
	//open sql data
	jsonFile, err := os.Open("cnf/sql.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	data, err := ioutil.ReadFile("cnf/sql.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var obj model.Configuration
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := repository.OpenConnection(obj)
	defer db.Close()
	http.HandleFunc("/creat-user", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "خطا در خواندن درخواست", http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "خطا در خواندن درخواست", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		//login//
		var data map[string]string
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "خطا در پارس کردن JSON", http.StatusInternalServerError)
			return
		}

		username, ok := data["username"]
		if !ok {
			http.Error(w, "Username not found", http.StatusBadRequest)
			return
		}
		password, ok := data["password"]
		if !ok {
			http.Error(w, "Password not found", http.StatusBadRequest)
			return
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM tbl_users WHERE username = $1", username).Scan(&count)
		if err != nil {
			http.Error(w, "Error checking username", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "This username is already taken", http.StatusBadRequest)
			return
		}
		if strings.Contains(username, " ") {
			http.Error(w, "Username cannot contain spaces", http.StatusBadRequest)
			return
		}
		_, err = db.Exec("INSERT INTO tbl_users (username, password) VALUES ($1, $2)", username, password)
		if err != nil {
			http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "User created successfully!")

	})

	log.Println("Starting server...")
	//http.ListenAndServe(":8085", nil)
	l, err := net.Listen("tcp", "localhost:8085")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(l, nil))
	log.Println("Sending request...")
	res, err := http.Get("http://localhost:8085/creat-user")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Reading response...")
	if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}
	bodyString := string(body)
	log.Print(bodyString)
}

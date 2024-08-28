package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type User struct {
	ID          string `json:"ID"`
	UserName    string `json:"UserName"`
	PhoneNumber string `json:"PhoneNumber"`
	Address     string `json:"Address"`
}

func (u User) User_Data() string {
	return fmt.Sprintf("ID: %s, UserName: %s, Phone Number: %s, Address: %s", u.ID, u.UserName, u.PhoneNumber, u.Address)
}

func HandleGetUser(users map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/users/"):]
		user, ok := users[id]
		if !ok {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		userData := user.User_Data()
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(userData))
	}
}

func HandleGetUsers(users map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData := ""
		for key, value := range users {
			userData += "UserName of[" + key + "]: " + value.UserName + "\n"
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(userData))
	}
}

func HandleAddUser(users map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var newUser User
		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			err := json.NewDecoder(r.Body).Decode(&newUser)
			if err != nil {
				http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
				return
			}
			users[newUser.ID] = newUser
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "User created successfully")
		} else {
			fmt.Fprintf(w, "please use json content format when adding a new User")
		}
	}
}
func main() {
	users := make(map[string]User)
	file, err := os.Open("userfile.json")
	if err != nil {
		fmt.Println("File opening error:", err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("File reading error:", err)
		return
	}
	var usersData []User
	err = json.Unmarshal(content, &usersData)
	if err != nil {
		fmt.Println("JSON decoding error:", err)
		return
	}

	for _, user := range usersData {
		users[user.ID] = user
	}
	http.HandleFunc("/users", HandleGetUsers(users))
	http.HandleFunc("/users/", HandleGetUser(users))
	http.HandleFunc("/users/adduser", HandleAddUser(users))
	fmt.Println("started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

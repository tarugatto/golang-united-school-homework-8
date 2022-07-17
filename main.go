package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	fileName := args["fileName"]
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}

		addUser(fileName, item, writer)
	case "list":
		getUsers(fileName, writer)
	case "findById":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}

		findById(fileName, id, writer)
	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}

		removeUser(fileName, id, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	id := flag.String("id", "", "user ID")
	item := flag.String("item", "", "json")
	operation := flag.String("operation", "", "method")
	fileName := flag.String("fileName", "", "file name")

	flag.Parse()

	return Arguments{
		"id":        *id,
		"item":      *item,
		"operation": *operation,
		"fileName":  *fileName,
	}
}

func getUsers(fileName string, w io.Writer) {
	f, _ := ioutil.ReadFile(fileName)

	w.Write(f)
}

func writeToFile(f *os.File, b *[]byte) {
	f.Seek(0, io.SeekStart)
	f.Truncate(0)
	f.Write(*b)
}

func addUser(fileName, item string, w io.Writer) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	b, _ := ioutil.ReadAll(f)

	users := make([]User, 0, 1)
	if len(b) > 0 {
		_ = json.Unmarshal(b, &users)
	}

	var newUser User
	json.Unmarshal([]byte(item), &newUser)
	for _, user := range users {
		if newUser.Id == user.Id {
			w.Write([]byte("Item with id " + user.Id + " already exists"))
			return
		}
	}

	users = append(users, newUser)
	jsn, _ := json.Marshal(&users)

	writeToFile(f, &jsn)
}

func removeUser(fileName, id string, w io.Writer) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	var users []User
	if err := json.Unmarshal(b, &users); err != nil {
		return
	}

	newUsers := make([]User, 0, len(users))
	for _, user := range users {
		if user.Id != id {
			newUsers = append(newUsers, user)
		}
	}

	if len(users) == len(newUsers) {
		w.Write([]byte("Item with id " + id + " not found"))
		return
	}

	jsn, _ := json.Marshal(newUsers)
	writeToFile(f, &jsn)
}

func findById(fileName, id string, w io.Writer) {
	var users []User
	b, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(b, &users)

	for _, user := range users {
		if user.Id == id {
			b, _ = json.Marshal(user)
			w.Write(b)
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt *time.Time `sql:"index"`
	Fname string
	Lname string
	Age   int
}

var Db *gorm.DB

func main() {
	var err error
	Db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=testpostgres password=mysecretpassword sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	Db.LogMode(true)
	defer Db.Close()
	Db.AutoMigrate(&Person{})
	pingerr := Db.DB().Ping()
	if pingerr != nil {
		fmt.Println("error database")
		fmt.Println(pingerr)
	}
	// migrateerror := Db.AutoMigrate(Person{}).Error
	// if migrateerror != nil {
	//  fmt.Println("migrate database")
	//  fmt.Println(migrateerror)
	// }
	fmt.Println("hello world")
	addnumbers(4, 8)
	repeat()
	var p1 Person
	p1.Fname = "john"
	p1.Lname = "doe"
	p1.Age = 24
	p1 = doubleage(p1)
	fmt.Printf("firstname %s, lastname %s, age is %d", p1.Fname, p1.Lname, p1.Age)
	fmt.Println()
	mybytes, _ := json.Marshal(p1)
	fmt.Println(string(mybytes))
	fmt.Println(p1.fullname())
	http.HandleFunc("/", hello)
	http.HandleFunc("/person2", jsonoutput)
	http.HandleFunc("/add", addnumber)
	http.HandleFunc("/api/savepersons", trial)
	http.HandleFunc("/api/listperson", listpersons)
	http.HandleFunc("/api/deletepersons", deleteperson)
	http.ListenAndServe(":8085", nil)
}
func addnumbers(a int, b int) {
	fmt.Println(a + b)
}
func repeat() {
	for i := 1; i < 10; i++ {
		fmt.Println(i)
		if i == 2 {
			fmt.Println("liverpool")
		}
	}
}
func doubleage(p Person) Person {
	p.Age = p.Age * 2
	return p
}
func (p Person) fullname() string {
	return p.Fname + " " + p.Lname
}
func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello People!"))
}
func jsonoutput(w http.ResponseWriter, r *http.Request) {
	p2 := Person{
		Fname: "kiro",
		Lname: "muhindo",
		Age:   70,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p2)
}
func addnumber(w http.ResponseWriter, r *http.Request) {
	// var err error
	// var a, b []string
	a, ok := r.URL.Query()["a"]
	b, ok := r.URL.Query()["b"]
	if ok == false {
		fmt.Println("problem with query")
		http.Error(w, "No parameters", http.StatusBadRequest)
		return
	}
	a1, err := strconv.ParseInt(a[0], 10, 64)
	b1, err := strconv.ParseInt(b[0], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "total is %d", a1+b1)
}
func trial(w http.ResponseWriter, r *http.Request) {
	p := Person{}
	json.NewDecoder(r.Body).Decode(&p)
	Db.Create(&p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
func listpersons(w http.ResponseWriter, r *http.Request) {
	pmany := []Person{}
	Db.Find(&pmany)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(pmany)
}
func deleteperson(w http.ResponseWriter, r *http.Request) {
	id, _ := r.URL.Query()["id"]
	theid, _ := strconv.ParseUint(id[0], 10, 64)
	p := Person{}
	p.ID = uint(theid)
	Db.Delete(&p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

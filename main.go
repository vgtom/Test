package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	Num1 string `json:"num1"`
	Num2 string `json:"num2"`
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE number (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"num1" TEXT,
		"num2" TEXT,
		"result" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("number table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertStudent(db *sql.DB, code string, name string, program string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO number(num1, num2, result) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM number ORDER BY num1")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Number: ", code, "+ ", name, "=", program)
	}
}

func add(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	intVar1, err := strconv.Atoi(book.Num1)
	if err != nil {

	}
	intVar2, err := strconv.Atoi(book.Num2)

	res := intVar1 + intVar2

	t := strconv.Itoa(res)

	c.JSON(http.StatusCreated, res)

	file, err := os.Create("sqlite.db")
	if err != nil {
		fmt.Println("Error")
	}
	file.Close()

	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		fmt.Println("Error")
	}
	defer db.Close()

	createTable(db)

	insertStudent(db, book.Num1, book.Num2, t)

	displayStudents(db)
}

func main() {
	r := gin.New()

	r.GET("/add", func(c *gin.Context) {
		c.JSON(http.StatusOK, "0")
	})

	r.POST("/add", add)

	r.Run()

}

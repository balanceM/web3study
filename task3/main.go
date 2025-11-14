package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func main() {
	dsn := "root:liu123@tcp(127.0.0.1:3306)/goprac?charset=utf8mb4&parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Database connect failed: %v\n", err.Error())
		return
	}
	//
	emps := []Employee{}
	err = db.Select(&emps, "Select * from employees where department = '技术部'")
	if err != nil {
		fmt.Printf("Get employees failed: %v\n", err.Error())
		return
	}
	fmt.Println(emps)
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/spf13/viper"
)

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func randomNumber(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(49 + rand.Intn(8)) //1=49 and 9 = 49+8=57
	}
	return string(bytes)
}

func createTenant(name string, db *sql.DB, tablesNum int, recordsNum int) error {
	// creates schema
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", name))

	for i := 1; i <= tablesNum; i++ {
		tableName := fmt.Sprintf("table%d", i)
		query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
			key serial primary key,
			caller VARCHAR(15) not null,
			callee VARCHAR(15) not null,
			duration integer not null,
			status integer not null,
			ts DATE not null,
			record text);`, name, tableName)

		_, err = db.Exec(query)
		if err != nil {
			break
		}

		fmt.Printf("\tTable %s created\n", tableName)

		// from https://godoc.org/github.com/lib/pq
		txn, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := txn.Prepare(pq.CopyInSchema(name, tableName, "caller", "callee", "duration", "status", "ts", "record"))
		if err != nil {
			log.Fatal(err)
		}

		for j := 0; j < recordsNum; j++ {
			caller := randomNumber(10)
			callee := randomNumber(10)
			duration := rand.Intn(300)
			status := rand.Intn(10)
			ts := time.Now()
			record := randomString(32)

			_, err = stmt.Exec(caller, callee, duration, status, ts, record)
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\t\t%d records added to %s\n", recordsNum, tableName)

		err = stmt.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = txn.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}
	return err
}

// CreateTenants creates num tenants
func CreateTenants(num int) {
	connStr := viper.GetString("database_url")
	fmt.Printf("Conn string is %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to open DB: ", err)
	}

	defer db.Close()

	tablesNum := viper.GetInt("tables")
	recordsNum := viper.GetInt("records")

	for i := 1; i <= num; i++ {
		tenantName := fmt.Sprintf("client%d", i)
		if err := createTenant(tenantName, db, tablesNum, recordsNum); err != nil {
			fmt.Printf("Failed to create client %d: %s\n", i, err)
			return
		}
		fmt.Printf("Tenant %s created\n", tenantName)
	}
}

// TestDB create test query
func TestDB() {
	connStr := "postgres://postgres:postgres@localhost:5432/gobase?sslmode=disable" // viper.GetString("database_url")
	fmt.Printf("Conn string is %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to open DB: ", err)
	}

	defer db.Close()

	fmt.Println("DB opened")
	rows, err := db.Query(`SELECT id, updated_at from profiles`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Rows got")
	var (
		id      int
		updated time.Time
	)

	for rows.Next() {
		err := rows.Scan(&id, &updated)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, updated)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

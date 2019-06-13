package storage

import (
	"database/sql"
	"fmt"
	"os"
	// postgres
	_ "github.com/lib/pq"
)

var db *sql.DB

//Connect initializes the database
func Connect(){

	databaseHost := os.Getenv("DOCTOR_DB_HOST")
	databasePort := os.Getenv("DOCTOR_DB_PORT")
	databaseUser := os.Getenv("DOCTOR_DB_USER")
	databasePass := os.Getenv("DOCTOR_DB_PASS")
	databaseName := os.Getenv("DOCTOR_DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		databaseHost, databasePort, databaseUser, databasePass, databaseName)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

}

// Ping check connection to DB
func Ping() {
	err := db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

// InsertAlert insert alert into DB
func InsertAlert(alertname, startsAT, address, status string) {
	sqlStatement := `INSERT INTO incidents (alertname, starts_at, address, status)
						VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(sqlStatement, alertname, startsAT, address, status)
	if err != nil {
		panic(err)
	}

}

// InsertAlertUnique insert alert into DB unique
func InsertAlertUnique(alertname, startsAT, address, status string) {
	id := GetAlertID(alertname, startsAT, address)
	if id != -1 {
		sqlStatement := `UPDATE incidents SET status = $2 WHERE id = $1;`
		_, err := db.Exec(sqlStatement, id, status)
		if err != nil {
			panic(err)
		}

	} else {
		InsertAlert(alertname, startsAT, address, status)
	}

}

// GetAlertID returns the alert id
func GetAlertID(alertname, startsAT, address string) int {
	var id int

	err := db.QueryRow(`SELECT id FROM incidents WHERE alertname = $1 and starts_at = $2 and address = $3`,
		alertname, startsAT, address).Scan(&id)

	if err == sql.ErrNoRows {
		//log.Fatal("No Results Found")
		id = -1
	}

	if err != nil {
		//log.Fatal(err)
	}
	return id
}

// GetScript returns the script
func GetScript(alertType string) string {
	var script string
	err := db.QueryRow(`SELECT script FROM mapping WHERE alert_type = $1`,
		alertType).Scan(&script)

	if err == sql.ErrNoRows {
		//log.Fatal("No Results Found")
		return "no script"
	}

	if err != nil {
		//log.Fatal(err)
	}
	return script
}

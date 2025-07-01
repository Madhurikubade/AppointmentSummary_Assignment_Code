package database

import (
	"AppointmentSummary_Assignment_Code/config"
	"AppointmentSummary_Assignment_Code/model"
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	Db  *sql.DB
	err error
)

// ConnectMysqlDB establishes a connection to the MySQL database.
func ConnectMysqlDB() {
	cfg := mysql.Config{
		User:                 config.AppConfig.MysqlConf.Username,
		Passwd:               config.AppConfig.MysqlConf.Password,
		Net:                  config.AppConfig.MysqlConf.Net,
		Addr:                 config.AppConfig.MysqlConf.Address,
		DBName:               config.AppConfig.MysqlConf.DatabaseName,
		AllowNativePasswords: config.AppConfig.MysqlConf.AllowNativePasswords,
		ParseTime:            config.AppConfig.MysqlConf.ParseTime,
	}

	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Error connecting to MySQL DB:", err)
	}

	Db.SetMaxOpenConns(config.AppConfig.MysqlConf.MaxConnections)
	Db.SetConnMaxIdleTime(time.Duration(config.AppConfig.MysqlConf.MaxIdleTimeMinutes) * time.Minute)

	createTables()
	loadCSVData()

	if err = Db.Ping(); err != nil {
		log.Fatal("Error pinging MySQL DB:", err)
	}

	log.Println("MySQL Database Connected")
}

// CloseMysql closes the MySQL connection.
func CloseMysql() {
	if Db != nil {
		if err := Db.Close(); err != nil {
			log.Println("Error closing MySQL connection:", err)
		} else {
			log.Println("Closed MySQL connection successfully!")
		}
	}
}

// createTables creates necessary tables if they don't exist.
func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS center (
			CenterID INT PRIMARY KEY,
			CenterName VARCHAR(255)
		);`,
		`CREATE TABLE IF NOT EXISTS patient (
			PatientID INT PRIMARY KEY,
			Salutation VARCHAR(10),
			Name VARCHAR(255),
			Mobile VARCHAR(20)
		);`,
		`CREATE TABLE IF NOT EXISTS doctor_staff (
			DoctorStaffID INT PRIMARY KEY,
			Name VARCHAR(255),
			Mobile VARCHAR(20)
		);`,
		`CREATE TABLE IF NOT EXISTS appointment (
			AppointmentID INT PRIMARY KEY,
			CenterID INT,
			DoctorStaffID INT,
			PatientID INT,
			StartTime DATETIME,
			EndTime DATETIME,
			Status VARCHAR(10),
			Category VARCHAR(100),
			FOREIGN KEY (CenterID) REFERENCES center(CenterID),
			FOREIGN KEY (DoctorStaffID) REFERENCES doctor_staff(DoctorStaffID),
			FOREIGN KEY (PatientID) REFERENCES patient(PatientID)
		);`,
		`CREATE TABLE IF NOT EXISTS doctor_summary_messages (
			ID INT AUTO_INCREMENT PRIMARY KEY,
			DoctorStaffID INT,
			DoctorMobile VARCHAR(20),
			Date DATE,
			Message TEXT,
			FOREIGN KEY (DoctorStaffID) REFERENCES doctor_staff(DoctorStaffID)
		);`,
		`CREATE TABLE IF NOT EXISTS center_summary_messages (
			ID INT AUTO_INCREMENT PRIMARY KEY,
			CenterID INT,
			Date DATE,
			Message TEXT,
			FOREIGN KEY (CenterID) REFERENCES center(CenterID)
		);`,
	}

	for _, q := range queries {
		if _, err := Db.Exec(q); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}

// loadCSVData loads CSV files into their respective tables.
func loadCSVData() {
	loadCSV("data/Center.csv", "INSERT IGNORE INTO center (CenterID, CenterName) VALUES (?, ?)")
	loadCSV("data/Patient.csv", "INSERT IGNORE INTO patient (PatientID, Salutation, Name, Mobile) VALUES (?, ?, ?, ?)")
	loadCSV("data/DoctorStaff.csv", "INSERT IGNORE INTO doctor_staff (DoctorStaffID, Name, Mobile) VALUES (?, ?, ?)")
	loadCSV("data/Appointment.csv", "INSERT IGNORE INTO appointment (AppointmentID, CenterID, DoctorStaffID, PatientID, StartTime, EndTime, Status, Category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
}

// loadCSV reads a CSV file and inserts its data into the DB.
func loadCSV(path string, query string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open CSV file %s: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV data from %s: %v", path, err)
	}

	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to prepare insert statement for %s: %v", path, err)
	}
	defer stmt.Close()

	for _, record := range records[1:] { // skip header
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = strings.TrimSpace(v)
		}
		if _, err := stmt.Exec(args...); err != nil {
			log.Printf("Error inserting into %s: %v", path, err)
		}
	}
}

// ReadDataForDate returns all valid (non-cancelled) appointments for a specific date.
func ReadDataForDate(date string) ([]any, error) {
	query := `
	SELECT 
		c.CenterName, d.Name, d.Mobile, a.Category, a.StartTime, a.EndTime, a.DoctorStaffID, a.CenterID
	FROM appointment a
	JOIN center c ON a.CenterID = c.CenterID
	JOIN doctor_staff d ON a.DoctorStaffID = d.DoctorStaffID
	WHERE DATE(a.StartTime) = ? AND a.Status = 'S';
	`

	rows, err := Db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []any
	for rows.Next() {
		var appt model.AppointmentData
		if err := rows.Scan(&appt.CenterName, &appt.DoctorName, &appt.DoctorMobile, &appt.Category, &appt.StartTime, &appt.EndTime, &appt.DoctorID, &appt.CenterID); err != nil {
			return nil, err
		}
		results = append(results, appt)
	}

	return results, nil
}

// InsertDoctorSummary inserts a doctor summary message into the DB.
func InsertDoctorSummary(msg model.DoctorMessage) error {
	_, err := Db.Exec("INSERT INTO doctor_summary_messages (DoctorStaffID, DoctorMobile, Date, Message) VALUES (?, ?, ?, ?)",
		msg.DoctorID, msg.DoctorMobile, msg.Date, msg.Message)
	return err
}

// InsertCenterSummary inserts a center summary message into the DB.
func InsertCenterSummary(msg model.CenterMessage) error {
	_, err := Db.Exec("INSERT INTO center_summary_messages (CenterID, Date, Message) VALUES (?, ?, ?)",
		msg.CenterID, msg.Date, msg.Message)
	return err
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Response struct {
	Hospitals []*Hospital
}

type Hospital struct {
	Id       int
	Name     string
	Staffs   []*Staff
	Patients []*Patient
	Adresses []*Address
}
type HospitalResponse struct {
	Id   int64
	Name string
}

type Address struct {
	Id         int
	HospitalId int64
	Region     string
	Street     string
}

type Staff struct {
	Id          int
	HospitalID  int64
	FullName    string
	PhoneNumber string
}

type Patient struct {
	Id          int
	HospitalId  int64
	FullName    string
	PatientInfo string
	PhoneNumber string
}

const (
	PostgresHost     = "localhost"
	PostgresPort     = 5432
	PostgresUser     = "ismoiljon12"
	PostgresPassword = "12"
	PostgresDatabase = "hospital_migration"
)

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDatabase,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Failed to connect database", err)
		return
	}
	defer db.Close()
	CreateInfo(db)
	// UpdateInfo(hospitalID,db)

	// DeleteInfo(5, db)
	// GettAllInfo(6, db)
}

func CreateInfo(db *sql.DB) int {
	info := &Hospital{
		Name: "Medical",
		Staffs: []*Staff{
			{
				FullName:    "Bobir Davlatov",
				PhoneNumber: "+998908907867",
			},
			{
				FullName:    "Komil Sattorov",
				PhoneNumber: "+998988765432",
			},
		},
		Patients: []*Patient{
			{
				FullName:    "Lobarxon",
				PhoneNumber: "+998994432345",
				PatientInfo: "Bu bemor hozirgi holati o'rta",
			},
			{
				FullName:    "Gavhar Ismoilova",
				PhoneNumber: "+998995455656",
				PatientInfo: "Bu bemor hozirgi holati yaxshi",
			},
		},
		Adresses: []*Address{
			{
				Region: "Chilonzor",
				Street: "Qatortol 68",
			},
			{
				Region: "Yunusobod",
				Street: "Amir Temur 78",
			},
		},
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Print("Failed to begin:", err)
		return 0
	}

	query1 := `
	INSERT INTO
		hospital(name)
	VALUES
		($1)
	RETURNING
		id, name
	`
	var response HospitalResponse
	err = tx.QueryRow(query1, info.Name).Scan(&response.Id, &response.Name)
	if err != nil {
		tx.Rollback()
		log.Println("Error while inserting hospital info", err)
	}
	fmt.Println(response)

	for _, staff := range info.Staffs {
		query2 := `
		INSERT INTO
			staff(hospital_id, full_name, phone_numbers)
		VALUES
			($1, $2, $3)
		`
		res, err := tx.Exec(query2, response.Id, staff.FullName, staff.PhoneNumber)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error while inserting staff info: ", err)
		}
		if count, _ := res.RowsAffected(); count == 0 {
			fmt.Println("Error while1 inserting staff info: ", err)
		}
	}

	for _, patient := range info.Patients {
		query3 := `
		INSERT INTO
			patients(hospital_id, full_name, patient_info, phone_number)
		VALUES
			($1, $2, $3, $4)
		`
		res, err := tx.Exec(query3, response.Id, patient.FullName, patient.PatientInfo, patient.PhoneNumber)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error while inserting patients info: ", err)
		}
		if count, _ := res.RowsAffected(); count == 0 {
			fmt.Println("Error while inserting patients info: ", err)
		}
	}

	for _, address := range info.Adresses {
		query4 := `
		INSERT INTO
			addresses(hospital_id, region, street)
		VALUES
			($1, $2, $3)
		`
		res, err := tx.Exec(query4, response.Id, address.Region, address.Street)
		if err != nil {
			tx.Rollback()
			log.Println("Error while create address info: ", err)
		}
		if count, _ := res.RowsAffected(); count == 0 {
			fmt.Println("Error while inserting patients info: ", err)
		}

	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to commit:", err)
		return 0
	}
	return int(response.Id)

}

func UpdateInfo(id int, db *sql.DB) {
	info := Hospital{
		Name: "Akfa Medline",
		Staffs: []*Staff{
			{
				FullName:    "Ismoiljon Abdurahomonov",
				PhoneNumber: "+99893 555 45 45",
			},
		},
		Patients: []*Patient{
			{
				FullName:    "John Doe",
				PatientInfo: "Uretrit kassaligi",
				PhoneNumber: "+99895 458 65 56",
			},
		},
		Adresses: []*Address{
			{
				Region: "Tashkent,Chilanzar",
				Street: "Qatortol",
			},
		},
	}

	_, err := db.Exec("UPDATE hospital SET name=$1 WHERE id=$2", info.Name, id)
	if err != nil {
		fmt.Println("Failed to update hospital: ", err)
		return
	}

	for _, patient := range info.Patients {

		_, err = db.Exec("UPDATE patients SET full_name=$1, patient_info=$2, phone_number=$3 WHERE hospital_id=$4", patient.FullName, patient.PatientInfo, patient.PhoneNumber, id)
		if err != nil {
			fmt.Println("Failed to update patients: ", err)
			return
		}
	}
	for _, addres := range info.Adresses {
		_, err = db.Exec("UPDATE addresses SET region=$1, street=$2 WHERE hospital_id=$3", addres.Region, addres.Street, id)
		if err != nil {
			fmt.Println("Failed to update addresses:", err)
			return
		}
	}

	for _, staff := range info.Staffs {
		_, err = db.Exec("UPDATE staff SET full_name=$1, phone_number=$2 WHERE hospital_id=$3", staff.FullName, staff.PhoneNumber, id)
		if err != nil {
			fmt.Println("Failed  to update staff:", err)
			return
		}
	}
}

func DeleteInfo(hospital_id int, db *sql.DB) {
	_, err := db.Exec("DELETE FROM hospital WHERE id=$1", hospital_id)
	if err != nil {
		fmt.Println("Failed to delete hospital:", err)
		return
	}

	_, err = db.Exec("DELETE FROM staff WHERE hospital_id=$1", hospital_id)
	if err != nil {
		fmt.Println("Failed to delete hospital:", err)
		return
	}

	_, err = db.Exec("DELETE FROM addresses WHERE hospital_id=$1", hospital_id)
	if err != nil {
		fmt.Println("Failed to delete hospital:", err)
		return
	}

	_, err = db.Exec("DELETE FROM patients WHERE hospital_id=$1", hospital_id)
	if err != nil {
		fmt.Println("Failed to delete patients:", err)
		return
	}

}

func GettAllInfo(id int, db *sql.DB) {
	resp := Response{}

	hRows, err := db.Query("SELECT id, name FROM hospital")
	if err != nil {
		fmt.Println("Failed to select hospital")
	}
	for hRows.Next() {
		var hospital Hospital
		err := hRows.Scan(
			&hospital.Id,
			&hospital.Name,
		)
		if err != nil {
			fmt.Println("Failed to scan hospital:", err)
			return
		}
		resp.Hospitals = append(resp.Hospitals, &hospital)

		sRows, err := db.Query("SELECT id,hospital_id,full_name,phone_numbers FROM staff WHERE hospital_id=$1", id)
		if err != nil {
			fmt.Println("Failed to select staff:", err)
		}

		for sRows.Next() {
			var staff Staff
			err := sRows.Scan(
				&staff.Id,
				&staff.HospitalID,
				&staff.FullName,
				&staff.PhoneNumber,
			)
			if err != nil {
				fmt.Println("Failed to scan staff:", err)
				return
			}
			hospital.Staffs = append(hospital.Staffs, &staff)
		}

		pRows, err := db.Query("SELECT id,hospital_id,full_name,patient_info,phone_number FROM patients WHERE hospital_id=$1", id)
		if err != nil {
			fmt.Println("Failed to select patients:", err)
			return
		}

		for pRows.Next() {
			var patient Patient
			err := pRows.Scan(
				&patient.Id,
				&patient.HospitalId,
				&patient.FullName,
				&patient.PatientInfo,
				&patient.PhoneNumber,
			)
			if err != nil {
				fmt.Println("Failed to scan patient:", err)
				return
			}
			hospital.Patients = append(hospital.Patients, &patient)
		}

		aRows, err := db.Query("SELECT id, hospital_id,region,street FROM addresses WHERE hospital_id=$1", id)

		for aRows.Next() {
			var address Address
			err := aRows.Scan(
				&address.Id,
				&address.HospitalId,
				&address.Region,
				&address.Street,
			)
			if err != nil {
				fmt.Println("Failed to scan address:", err)
				return
			}
			hospital.Adresses = append(hospital.Adresses, &address)
		}

	}

	for _, hospital := range resp.Hospitals {
		fmt.Println(hospital)
		for _, staff := range hospital.Staffs {
			fmt.Println("Staff")
			fmt.Println(staff)
		}
		for _, patient := range hospital.Patients {
			fmt.Println(patient)
		}
		for _, address := range hospital.Adresses {
			fmt.Println(address)
		}
	}

}

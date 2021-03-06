package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseHandler struct {
	DriverName string
	User       string
	Password   string
	Database   string
}

var database *sql.DB

var stmtIns *sql.Stmt
var stmtOut *sql.Stmt
var stmtDel *sql.Stmt

func (dbReader DatabaseHandler) Begin() (err error) {
	dbOpenString := dbReader.User + ":" + dbReader.Password + "@tcp(sql7.freemysqlhosting.net)/" + dbReader.Database
	database, err = sql.Open(dbReader.DriverName, dbOpenString)
	return err
}

func (dbReader DatabaseHandler) Write(plantName string, wateredSoilMoisture int, positionX int, positionY int) (err error) {
	if stmtIns, err = database.Prepare("INSERT INTO plants_data (name, watered_soil_moisture, pos_x, pos_y) VALUES (?, ?, ?, ?)"); err != nil {
		return err
	}

	if _, err = stmtIns.Exec(plantName, wateredSoilMoisture, positionX, positionY); err != nil {
		return err
	}

	if err = stmtIns.Close(); err != nil {
		return err
	}

	return nil
}

func (dbReader DatabaseHandler) Update(plantName string, wateredSoilMoisture int) (err error) {
	if stmtIns, err = database.Prepare("UPDATE plants_data SET watered_soil_moisture = ? WHERE name = ?"); err != nil {
		return err
	}

	if _, err = stmtIns.Exec(wateredSoilMoisture, plantName); err != nil {
		return err
	}

	if err = stmtIns.Close(); err != nil {
		return err
	}

	return nil
}

func (dbReader DatabaseHandler) GetWateredSoilMoistureFromName(plantName string) (wateredSoilMoisture int, err error) {
	if stmtOut, err = database.Prepare("SELECT watered_soil_moisture FROM plants_data WHERE name = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(plantName).Scan(&wateredSoilMoisture); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err
	}

	return wateredSoilMoisture, err
}

func (dbReader DatabaseHandler) GetWateredSoilMoistureFromId(plantId int) (wateredSoilMoisture int, err error) {
	if stmtOut, err = database.Prepare("SELECT watered_soil_moisture FROM plants_data WHERE plant_id = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(plantId).Scan(&wateredSoilMoisture); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err
	}

	return wateredSoilMoisture, err
}

func (dbReader DatabaseHandler) GetPlantName(plant_id int) (plantName string, err error) {
	if stmtOut, err = database.Prepare("SELECT name FROM plants_data WHERE plant_id = ?"); err != nil {
		return "", err
	}

	if err = stmtOut.QueryRow(plant_id).Scan(&plantName); err != nil {
		return "", err
	}

	if err = stmtOut.Close(); err != nil {
		return "", err
	}

	return plantName, err
}

func (dbReader DatabaseHandler) GetAllPlantsNames() (plants []string, err error) {
	if stmtOut, err = database.Prepare("SELECT name FROM plants_data"); err != nil {
		return nil, err
	}

	rows, err := stmtOut.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		plants = append(plants, name)
	}

	return plants, nil
}

func (dbReader DatabaseHandler) GetAllPlantsSoilMoisture() (plants []string, err error) {
	if stmtOut, err = database.Prepare("SELECT watered_soil_moisture FROM plants_data"); err != nil {
		return nil, err
	}

	rows, err := stmtOut.Query()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var wateredSoilMoisture string

		err = rows.Scan(&wateredSoilMoisture)
		if err != nil {
			return nil, err
		}

		plants = append(plants, wateredSoilMoisture)
	}

	return plants, nil
}

func (dbHandler DatabaseHandler) GetPlantX(plantId int) (plantX int, err error) {
	if stmtOut, err = database.Prepare("SELECT pos_x FROM plants_data WHERE plant_id = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(plantId).Scan(&plantX); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err
	}

	return plantX, nil
}

func (dbHandler DatabaseHandler) GetPlantY(plantId int) (plantY int, err error) {
	if stmtOut, err = database.Prepare("SELECT pos_y FROM plants_data WHERE plant_id = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(plantId).Scan(&plantY); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err
	}

	return plantY, nil
}

func (dbReader DatabaseHandler) DeletePlant(plantName string) (err error) {
	if stmtDel, err = database.Prepare("DELETE FROM plants_data WHERE name = ?"); err != nil {
		return err
	}

	if _, err = stmtDel.Exec(plantName); err != nil {
		return err
	}

	if err = stmtDel.Close(); err != nil {
		return err
	}

	return nil
}

func (dbReader DatabaseHandler) GetNumberOfPlants() (numberOfPlants int, err error) {
	rows, err := database.Query("SELECT COUNT(*) FROM plants_data")

	if err != nil {
		return -1, err
	}

	for rows.Next() {
		err = rows.Scan(&numberOfPlants)
	}

	if err != nil {
		return -1, err
	}

	err = rows.Close()

	if err != nil {
		return -1, err
	}

	return numberOfPlants, nil
}

func (dbReader DatabaseHandler) GetPlantID(wateredSoilMoisture int) (plantID int, err error) {
	if stmtOut, err = database.Prepare("SELECT plant_id FROM plants_data WHERE watered_soil_moisture = ?"); err != nil {
		return -1, err
	}

	if err = stmtOut.QueryRow(wateredSoilMoisture).Scan(&plantID); err != nil {
		return -1, err
	}

	if err = stmtOut.Close(); err != nil {
		return -1, err
	}

	return plantID, err
}

package main

import (
	"aMazeGardenAI/db"
	"aMazeGardenAI/serverUtils"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Json received from aMazeGarden data logger template.
type MeteoData struct {
	AirTemperature  string `json:"airTemperature"`
	AirHumidity     string `json:"airHumidity"`
	SoilMoisture    string `json:"soilMoisture"`
	SoilTemperature string `json:"soilTemperature"`
}

// Database handler object.
var dbHandler = db.DatabaseHandler{
	DriverName: "mysql",
	User:       "aMazeGardenAI",
	Password:   "Tvlwxfcg+1q",
	Database:   "plants",
}

var useForecast = true

var maxSoilTemp = 25

var meteoData MeteoData

var addPlantFormLoader = serverUtils.FileLoader{Path: "templates/addPlantForm.html"}
var editPlantFormLoader = serverUtils.FileLoader{Path: "templates/editPlantForm.html"}
var removePlantFormLoader = serverUtils.FileLoader{Path: "templates/removePlantForm.html"}

func main() {

	// Handle function for receiving data from data logger.
	http.HandleFunc("/", handle)
	http.HandleFunc("/addPlant", addPlantFormLoader.LoadFile)
	http.HandleFunc("/removePlant", deletePlantFormPage)
	http.HandleFunc("/editPlant", editPlantFormPage)
	http.HandleFunc("/addPlantDb", addPlantToDb)
	http.HandleFunc("/removePlantDb", removePlantFormLoader.LoadFile)
	http.HandleFunc("/editPlantDB", editPlantFormLoader.LoadFile)
	http.HandleFunc("/dataLoggerData", handleMeteoDataRequest)

	fmt.Println("Http server listening...")

	// Connecting todatabase.
	err := dbHandler.Begin()

	// Error handling.
	if err != nil {
		panic(err)
	}
	// Listen on port 8080.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func wateringAI(meteoData MeteoData) {
	fmt.Println("Got data: ", meteoData)

	var numberOfPlants int
	var err error

	// Get number of plants in database.
	if numberOfPlants, err = dbHandler.GetNumberOfPlants(); err != nil {
		fmt.Println("Error getting plants number.")
		panic(err)
	}

	// Scan for all plants.
	for i := 1; i < numberOfPlants; i++ {
		// Convert string to int.
		soilMoisture, err := strconv.Atoi(meteoData.SoilMoisture)

		if err != nil {
			fmt.Println("Error getting soil moisture.")
			panic(err)
		}

		soilTemperature, err := strconv.Atoi(meteoData.SoilTemperature)

		if err != nil {
			fmt.Println("Error getting soil temperature.")
			panic(err)
		}

		// Get watered soil moisture of plant from the database.
		wateredSoilMoisture, err := dbHandler.GetWateredSoilMoistureFromId(i)

		// Error handling.
		if err != nil {
			fmt.Println("Error getting watered soil moisture.")
			panic(err)
		}

		// If soil is dry.
		if soilMoisture < wateredSoilMoisture {
			if soilTemperature < maxSoilTemp {
				fmt.Println("Watering needed")
			}
		}
	}
}

func handle(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		panic(err)
	}

	if err := t.Execute(writer, meteoData); err != nil {
		panic(err)
	}
}

func editPlantFormPage(writer http.ResponseWriter, request *http.Request) {
	plants, err := dbHandler.GetAllPlants()

	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/editPlantForm.html")

	if err != nil {
		panic(err)
	}

	options := ""

	for i := 0; i < len(plants); i++ {
		options += "<option value=" + plants[i] + ">" + plants[i] + "</option>"
	}

	err = t.Execute(writer, template.HTML(options))

	if err != nil {
		panic(err)
	}
}

func deletePlantFormPage(writer http.ResponseWriter, request *http.Request) {
	plants, err := dbHandler.GetAllPlants()

	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/removePlantForm.html")

	if err != nil {
		panic(err)
	}

	options := ""

	for i := 0; i < len(plants); i++ {
		options += "<option value=" + plants[i] + ">" + plants[i] + "</option>"
	}

	err = t.Execute(writer, template.HTML(options))

	if err != nil {
		panic(err)
	}
}

func addPlantToDb(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			panic(err)
		}

		plantName := request.FormValue("plantName")
		plantWateredSoilMoisture, err := strconv.Atoi(request.FormValue("wateredSoilMoisture"))

		if err != nil {
			panic(err)
		}

		err = dbHandler.Write(plantName, plantWateredSoilMoisture)

		if err != nil {
			panic(err)
		}

		http.Redirect(writer, request, "/", 303)
	}
}

func deletePlantDb(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			panic(err)
		}

		plantName := request.FormValue("plantName")
		plantWateredSoilMoisture, err := strconv.Atoi(request.FormValue("wateredSoilMoisture"))

		if err != nil {
			panic(err)
		}

		err = dbHandler.Write(plantName, plantWateredSoilMoisture)

		if err != nil {
			panic(err)
		}

		http.Redirect(writer, request, "/", 303)
	}
}

/*
	Handling data logger request.
*/
func handleMeteoDataRequest(writer http.ResponseWriter, request *http.Request) {
	// If request method is POST.
	if request.Method == "POST" {
		// Read request.
		jsn, err := ioutil.ReadAll(request.Body)

		// Parsing request to string.
		requestMessage := string(jsn)

		// Error handling.
		if err != nil {
			panic(err)
		}

		// Decoding json.
		if err := json.Unmarshal([]byte(requestMessage), &meteoData); err != nil {
			panic(err)
		}

		// Calling AI fuction.
		wateringAI(meteoData)
	}
}

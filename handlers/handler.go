package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/MaximLanBowl/Searcher.git/models"
	"github.com/gin-gonic/gin"
)

func LoadDrivers(filePath string) ([]models.Driver, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var drivers []models.Driver
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		driverID, _ := strconv.Atoi(record[0])
		latitude, _ := strconv.ParseFloat(record[1], 64)
		longitude, _ := strconv.ParseFloat(record[2], 64)
		drivers = append(drivers, models.Driver{DriverID: driverID, Latitude: latitude, Longitude: longitude})
	}
	return drivers, nil
}

func GetETA(clientLat, clientLon, driverLat, driverLon float64) (float64, error) {
	url := "https://maps.starline.ru/api/routing/route"
	reqBody := map[string]interface{}{
		"locations": []map[string]interface{}{
			{"lat": clientLat, "lon": clientLon, "type": "break"},
			{"lat": driverLat, "lon": driverLon, "type": "break"},
		},
		"costing": "auto",
		"alternates": 1,
	}
	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var routeResponse models.RouteResponse
    err = json.NewDecoder(resp.Body).Decode(&routeResponse)
    if err != nil {
        return 0, err
    }

	if len(routeResponse.Routes) > 0 {
		return routeResponse.Routes[0].Duration, nil 
	}
	return 0, nil
}

var drivers []models.Driver
var once sync.Once

func LoadDriversOnce(){
	once.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		filePath := filepath.Join(dir, "driver_positions.csv")
		drivers, err = LoadDrivers(filePath)
		if err != nil {
			log.Fatalf("Ошибка загрузки водителей: %v", err) 
		}
	})
}

func SearchDriver(c *gin.Context) {
    lat, err := strconv.ParseFloat(c.Query("lat"), 64)
    if err != nil || lat == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
        return
    }
    lon, err := strconv.ParseFloat(c.Query("lon"), 64)
    if err != nil || lon == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
        return
    }

    LoadDriversOnce()

    var closestDriver models.Driver
    var minETA float64 = -1

    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, driver := range drivers {
        wg.Add(1)
        go func(driver models.Driver) {
            defer wg.Done()
            eta, err := GetETA(lat, lon, driver.Latitude, driver.Longitude)
            if err != nil {
                log.Println(err)
                return
            }
            mu.Lock()
            if minETA == -1 || eta < minETA {
                minETA = eta
                closestDriver = driver
            }
            mu.Unlock()
        }(driver)
    }
    wg.Wait()

    if minETA == -1 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to find closest driver"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"driver_id": closestDriver.DriverID})
}

// SearchDriver godoc
// @Summary Search for the closest driver
// @Description get closest driver based on client's location
// @ID search-driver
// @Accept  json
// @Produce  json
// @Param lat query float64 true "Client Latitude"
// @Param lon query float64 true "Client Longitude"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /driverSearch [get]
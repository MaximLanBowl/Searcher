package handlers

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

func LoadDrivers(filePath string) ([]Driver, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var drivers []Driver
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
		drivers = append(drivers, Driver{driverID})
	}
	
}
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"encoding/csv"
	"os"
	"io"
	"time"
	"strconv"
	"strings"
)

func main() {

	filename, err := filepath.Abs("./res/PA_Interview_Questions_Data.csv")
	time_string := "15:04:05" // Golang odd time format...

	if err != nil {
        log.Fatal(err)
    }

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

    r := csv.NewReader(file)

	// Get header (we are sure that first line must be the header).
	header := make([]string, 0)

	record, err := r.Read()

	if err != nil {
		log.Fatal(err)
	}

	// Mapping each column id to its string value.
	for i := range record {
		header = append(header, record[i])
	}

	// Time series in the CSV is assumed as ordered.
	var timeseries = make([]time.Time, 0)
	var stored_records = make([][]float64, 0)
	var average = make([]float64, len(record)-1)
	row_length := len(record)-1

	for {
			record, err := r.Read()

			if err != nil {
				log.Fatal(err)
			}

			if err == io.EOF {
				break
			}
			
			// Read the time value.
			current_time, err := time.Parse(time_string, record[0]) // Get time object from string.
			if err != nil {
				log.Fatal(err)
			}
			timeseries = append(timeseries, current_time)

			// Add to stored records.
			new_row := make([]float64, row_length)
			for i, v := range record[1:] {
				new_row[i] = convertToFloat(v)
			}
			stored_records = append(stored_records, new_row)
			
			// Free memory and delete old values (>= 20 sec old)
			n := 0
			for k := range timeseries {
				if current_time.Sub(timeseries[k]).Seconds() >= 20 {
					n++
				}
			}
			timeseries = timeseries[n:]
			stored_records = stored_records[n:]

			// Reset and compute average values.
			for k := range average {
				average[k] = 0
			}
			
			for r := 0; r < len(stored_records); r++ { // For each row.
				for c := 0; c < row_length; c++ {	// For each column.
					average[c] = average[c] + stored_records[r][c]
				}
			}
			
			items := float64(len(stored_records))
			for k := range average {
				average[k] = average[k] / items
			}
			
			// Do whatever you want with the rolling average.
			fmt.Printf("Rolling average for the 20 seconds ending at %s (%.f items found):\n", timeseries[len(timeseries)-1].Format(time_string), items)
			for k := range header[1:] {
				fmt.Println("\t", header[k+1], "=", average[k])
			}
    }
}


func convertToFloat (s string) float64 {
	s = strings.Trim(s, "\"") // Strip quotes, if present.
	s = strings.ReplaceAll(s, ",", "") // Remove thousand separator.
	f, err := strconv.ParseFloat(s, 64) // Convert to float64.
	if err != nil {
		log.Fatal(err)
	}
	return f
}
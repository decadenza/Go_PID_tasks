package main

import (
	"fmt"
	"log"
	"time"
	"path/filepath"
	"io"
	"os"
	"encoding/csv"
	"strconv"
	"strings"
)


type PIDController struct {
	p          		float64	// Proportional gain.
	i         	 	float64	// Integral gain.
	d          		float64	// Derivative gain.
	setpoint  	 	float64	// Desired process value.
	integral_sum  	float64	// Integral sum over time.
	old_err   		float64	// Previous process value needed for derivative.
	oldTime 		time.Time 	// Time of last update, needed for derivative.
	out_min   		float64	// Output min
	out_max    		float64	// Output max
}


func main(){
	var controller PIDController
	controller.setpoint = 6 	// Target pressure GUESS!!!
	controller.out_min = 0 		// Min aperture.
	controller.out_max = 1 		// Max aperture.
	controller.p = 0.8
	controller.i = 0.01
	controller.d = 0.0001

	// Checking for coherence.
	if controller.out_max < controller.out_min {
		log.Fatal("Incorrect out_min and out_max values.")
	}

	header, csv := prepareCSV()

	// Main controlling loop.
	for {

		// Read time and pressure value from csv.
		// This is a mock-up! Real values must be read online!
		record, err := csv.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		G1 := convertToFloat(record[header["G1 (bar)"]])
		t, err := time.Parse("15:04:05", record[header["Timestamp"]])
		if err != nil {
			log.Fatal(err)
		}

		// Calculate controlling value P1.
		// P1++ (opening) => Pressure--
		// P1-- (closing) => Pressure++
		
		// Calculate controlling variable.
		P1 := PID_Step(G1, t, controller)

		// Send P1 into the system.
		// ...
		// Mock up: next output does not depend from my P1 value!!!


		// Mock up of aperture delay (has effect on live data only...).
		//time.Sleep(800 * time.Millisecond)
		
		fmt.Printf("At %s: aperture P1=%f (original %s), pressure G1=%f (error %f)\n", t.Format("15:04:05"), P1, record[header["P1"]], G1, controller.setpoint - G1 )
	}

}

func PID_Step(value float64, current_time time.Time, ctrl PIDController) float64 {

	err := ctrl.setpoint - value
	dt := current_time.Sub(ctrl.oldTime).Seconds()

	ctrl.integral_sum = ctrl.integral_sum + ctrl.i * err * dt

	// Anti-windup... OPTIONAL
	if ctrl.integral_sum > ctrl.out_max {
		ctrl.integral_sum = ctrl.out_max
	} else if ctrl.integral_sum < ctrl.out_min {
		ctrl.integral_sum = ctrl.out_min
	}

	output := ctrl.p * err + ctrl.integral_sum + ctrl.d * (err - ctrl.old_err)/dt
	
	// Update controller values for next derivative.
	ctrl.old_err = err
	ctrl.oldTime = current_time
	
	// Clip within output limits.
	if output > ctrl.out_max {
		output = ctrl.out_max
	} else if output < ctrl.out_min {
		output = ctrl.out_min
	}
	
	return output
}


func prepareCSV() (map[string]int, *csv.Reader){
	// Read time and pressure values from CSV.
	filename, err := filepath.Abs("./res/PA_Interview_Questions_Data.csv")
	if err != nil {
        log.Fatal(err)
    }
	
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

    r := csv.NewReader(file)

	// Header
	// Map string -> index int.
	header := make(map[string]int)
	record, err := r.Read()
	
	if err != nil {
		log.Fatal(err)
	}

	for i,v := range record {
		header[v] = i
	}

	return header, r
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
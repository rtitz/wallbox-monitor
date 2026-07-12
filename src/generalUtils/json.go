package generalUtils

import (
	"encoding/json"
	"fmt"
	"wallbox-monitor/variables"
	"wallbox-monitor/wallboxApi"
)

// Pre-allocate a buffer for JSON marshaling to reduce allocations
var jsonBuffer = make([]byte, 0, 512)

func CreateJsonStringFromWallbox(Date, Time, TimeZone string, status *wallboxApi.Status) []byte {
	data := variables.JsonData{
		Date:     Date,
		Time:     Time,
		TimeZone: TimeZone,
		Car:      status.Car,
		Ust:      status.Ust,
		Amp:      status.Amp,
		Nrg:      status.Nrg,
		Wh:       status.Wh,
		Dws:      status.Dws,
		Psm:      status.Psm,
		Frc:      status.Frc,
		Tma:      status.Tma,
		Alw:      bool(status.Alw),
	}

	jsonBuffer = jsonBuffer[:0]
	jsonString, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling wallbox JSON:", err)
		return []byte("{}")
	}
	return jsonString
}

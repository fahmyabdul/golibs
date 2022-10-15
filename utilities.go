package golibs

import (
	"math"
	"time"
)

// UtilsParseTime :
func UtilsParseTime(inputTime, layout, outLayout string) (string, error) {
	formattedTime, err := time.Parse(layout, inputTime)
	if err != nil {
		return inputTime, err
	}

	return formattedTime.Format(outLayout), nil
}

// InTimeSpan :
func InTimeSpan(format, startDate, endDate, checkDate string) bool {
	start, _ := time.Parse(format, startDate)
	end, _ := time.Parse(format, endDate)
	check, _ := time.Parse(format, checkDate)

	if check.Equal(start) || check.Equal(end) {
		return true
	}

	return check.After(start) && check.Before(end)
}

// DistanceLatLongInKm : Return Distance Between Current LatLong and Destination LatLong in Kilometers
func DistanceLatLongInKm(currentLatitude, currentLongitude, destinationLatitude, destinationLongitude float64) float64 {
	return math.Sqrt(
		math.Pow(111*(destinationLatitude-currentLatitude), 2) +
			math.Pow(111*(currentLongitude-destinationLongitude)*math.Cos(currentLatitude/57.3), 2),
	)
}

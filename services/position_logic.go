package services

import (
	"github.com/gerhardgruber/fame/models"
	"github.com/kellydunn/golang-geo"
)

// CalculateDistance calculates the distance between two Sophy Points
func CalculateDistance(pos1 *models.Position, pos2 *models.Position) float64 {
	p1 := geo.NewPoint(pos1.Longitude, pos1.Latitude)
	p2 := geo.NewPoint(pos2.Longitude, pos2.Latitude)

	return p1.GreatCircleDistance(p2)
}

// CalculateDistanceToDestination - Calculates the distance to the goal after receiving a Position from the mobilephone
func CalculateDistanceToDestination(models.Position, models.MobilePhone) float64 {

	// TODO: Calculate Distance between Mobile Device and Destination
	return 0.0
}

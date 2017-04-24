package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/toefel18/go-deliver/backend/delivery"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover(), middleware.CORS())

	e.GET("/trips", getTrips)
	e.GET("/trips/:id", getTrip)
	e.PATCH("/trips/:id", updateTrip)
	e.GET("/trips/:id/pieces/:pieceId", getPiece)
	e.PATCH("/trips/:id/pieces/:pieceId", updatePiece)
	e.GET("/trips/:id/stops/:stopNumber", getStop)
	e.Logger.Fatal(e.Start(atPort()))
}

func atPort() string {
	if port, set := os.LookupEnv("PORT"); set {
		return ":" + port
	}
	return ":8080"
}

func getTrips(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, []string{trip1.TripNumber}, " ")
}

func getTrip(c echo.Context) error {
	id := c.Param("id")
	if id != trip1.TripNumber {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip Not Found"})
	}
	return c.JSONPretty(http.StatusOK, trip1, "  ")
}

type PatchRequest struct {
	Operation string `json:"operation"`
	Signee    string `json:"signee,omitempty"`
}

func updateTrip(c echo.Context) error {
	id := c.Param("id")
	if id != trip1.TripNumber {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip Not Found"})
	}
	request := new(PatchRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Operation == "start delivery" {
		if err := trip1.StartDelivery(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "trip is in delivery"})
	} else if request.Operation == "finish" {
		if err := trip1.Finish(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "trip finished"})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"message": "unknown operation"})
}

func getPiece(c echo.Context) error {
	tripId := c.Param("id")
	if tripId != trip1.TripNumber {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip Not Found"})
	}
	piece := trip1.FindPiece(c.Param("pieceId"))
	if piece == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Piece Not Found"})
	}
	return c.JSONPretty(http.StatusOK, piece, "  ")
}

func updatePiece(c echo.Context) error {
	id := c.Param("id")
	if id != trip1.TripNumber {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip Not Found"})
	}
	piece := trip1.FindPiece(c.Param("pieceId"))
	if piece == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Piece Not Found"})
	}
	request := new(PatchRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Operation == "delivered" {
		piece.Delivered(request.Signee)
		return c.JSON(http.StatusOK, map[string]string{"message": "piece " + piece.Id + " delivered, signed by " + piece.Signee})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"message": "unknown operation"})
}

func getStop(c echo.Context) error {
	tripId := c.Param("id")
	if tripId != trip1.TripNumber {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Trip Not Found"})
	}

	stopNumber, err := strconv.Atoi(c.Param("stopNumber"))
	if err != nil {
		c.Error(err)
		return c.JSON(http.StatusNotFound, map[string]string{"message": "stop index must be a number"})
	}

	stop := trip1.FindStop(stopNumber)
	if stop == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Stop not found"})
	}
	return c.JSONPretty(http.StatusOK, stop, "  ")
}

/*  TEST DATA BELOW  */

var biltstraat7 = &delivery.Address{
	Kixcode:     "NL3572AA000007X",
	Street:      "Biltstraat",
	HouseNumber: "7",
	PostalCode:  "3572AA",
	City:        "Utrecht",
	Country:     "Nederland",
	Longitude:   5.126964,
	Latitude:    52.094944,
}

var biltstraat451 = &delivery.Address{
	Kixcode:     "NL3572AX000451X",
	Street:      "Biltstraat",
	HouseNumber: "451",
	PostalCode:  "3572AX",
	City:        "Utrecht",
	Country:     "Nederland",
	Longitude:   5.136621,
	Latitude:    52.095425,
}

var runnenbrug21 = &delivery.Address{
	Kixcode:     "NL3981AZ000021X",
	Street:      "Runnenburg",
	HouseNumber: "21",
	PostalCode:  "3981AZ",
	City:        "Bunnik",
	Country:     "Nederland",
	Longitude:   5.194597,
	Latitude:    52.063972,
}

var pieceRadio = &delivery.Piece{
	Id:             "JVGL566684224",
	ReceiverName:   "Connie Plessen",
	ShipmentNumber: 163549755,
	Status:         delivery.Sorted,
}

var pieceAntenna = &delivery.Piece{
	Id:             "JVGL566684225",
	ReceiverName:   "Connie Plessen",
	ShipmentNumber: 163549755,
	Status:         delivery.Sorted,
}

var pieceVoerbak = &delivery.Piece{
	Id:             "JVGL79984753",
	ReceiverName:   "Sjef Speciaal",
	ShipmentNumber: 763249799,
	Status:         delivery.Sorted,
}

var pieceKussen = &delivery.Piece{
	Id:             "JVGL89861267",
	ReceiverName:   "Ronald Reigan",
	ShipmentNumber: 497983327,
	Status:         delivery.Sorted,
}

var stopBiltstraat7 = &delivery.Stop{
	Address: biltstraat7,
	Pieces:  []*delivery.Piece{pieceRadio, pieceAntenna},
}

var stopBiltstraat451 = &delivery.Stop{
	Address: biltstraat451,
	Pieces:  []*delivery.Piece{pieceVoerbak},
}

var stopRunnenburg = &delivery.Stop{
	Address: runnenbrug21,
	Pieces:  []*delivery.Piece{pieceKussen},
}

var trip1 = &delivery.Trip{
	TripNumber: "744567",
	Stops:      []*delivery.Stop{stopBiltstraat7, stopBiltstraat451, stopRunnenburg},
	Status:     delivery.Ready,
}

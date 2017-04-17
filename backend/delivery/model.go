package delivery

import "errors"

var (
	Ready           = "ready"
	Sorted          = "sorted"
	InDelivery      = "in delivery"
	Delivered       = "delivered"
	Finished        = "finished"
	ReturnedToDepot = "returned to depot"
)

type Trip struct {
	TripNumber string  `json:"tripNumber"`
	Stops      []*Stop `json:"stops"`
	Status     string  `json:"status"`
}

type Stop struct {
	Address *Address `json:"address"`
	Pieces  []*Piece `json:"pieces"`
}

type Piece struct {
	Id             string `json:"id"`
	ReceiverName   string `json:"receiverName"`
	ShipmentNumber int64  `json:"shipmentNumber"`
	Status         string `json:"status"`
	Signee         string `json:"signee,omitempty"`
}

type Address struct {
	Kixcode     string  `json:"kixCode"`
	Street      string  `json:"street"`
	HouseNumber string  `json:"houseNumber"`
	PostalCode  string  `json:"postalCode"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func (trip *Trip) StartDelivery() error {
	if trip.Status != Ready {
		return errors.New("cannot start delivery of trip when it has state " + trip.Status + ", state must be " + Ready)
	}
	for _, stop := range trip.Stops {
		for _, piece := range stop.Pieces {
			piece.Status = InDelivery
		}
	}
	trip.Status = InDelivery
	return nil
}

func (trip *Trip) Finish() error {
	if trip.Status != InDelivery {
		return errors.New("cannot finish trip when it has state " + trip.Status + ", state must be " + InDelivery)
	}
	for _, stop := range trip.Stops {
		for _, piece := range stop.Pieces {
			if piece.Status == InDelivery {
				piece.Status = ReturnedToDepot
			}
		}
	}
	trip.Status = Finished
	return nil
}

func (trip *Trip) FindPiece(pieceId string) *Piece {
	for _, stop := range trip.Stops {
		for _, piece := range stop.Pieces {
			if piece.Id == pieceId {
				return piece
			}
		}
	}
	return nil
}

func (trip *Trip) FindStop(stopNumber int) *Stop {
	if stopNumber < len(trip.Stops) && stopNumber >= 0 {
		return trip.Stops[stopNumber]
	}
	return nil
}

func (piece *Piece) Delivered(signee string) {
	piece.Status = Delivered
	piece.Signee = signee
}

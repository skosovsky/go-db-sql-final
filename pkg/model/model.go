package model

type ParcelStatus string

type Parcel struct {
	ID        int
	ClientID  int
	Status    ParcelStatus
	Address   string
	CreatedAt string
}

const (
	ParcelStatusRegistered ParcelStatus = "registered"
	ParcelStatusSent       ParcelStatus = "sent"
	ParcelStatusDelivered  ParcelStatus = "delivered"
)

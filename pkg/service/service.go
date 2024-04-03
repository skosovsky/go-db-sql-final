package service

import (
	"fmt"
	"time"

	"github.com/skosovsky/go-db-sql-final.git/pkg/model"
	"github.com/skosovsky/go-db-sql-final.git/pkg/store"
)

type ParcelService struct {
	store store.ParcelStore
}

func NewParcelService(store store.ParcelStore) ParcelService {
	return ParcelService{store: store}
}

func (s ParcelService) Register(client int, address string) (model.Parcel, error) {
	parcel := model.Parcel{
		ID:        0,
		ClientID:  client,
		Status:    model.ParcelStatusRegistered,
		Address:   address,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.store.Add(parcel)
	if err != nil {
		return model.Parcel{}, fmt.Errorf("add parcel error: %w", err)
	}

	parcel.ID = id

	fmt.Printf("Новая посылка №%d на адрес %s от клиента с идентификатором %d зарегистрирована %s\n", //nolint:forbidigo // it's text app
		parcel.ID, parcel.Address, parcel.ClientID, parcel.CreatedAt)

	return parcel, nil
}

func (s ParcelService) PrintClientParcels(client int) error {
	parcels, err := s.store.GetByClient(client)
	if err != nil {
		return fmt.Errorf("get by client error: %w", err)
	}

	fmt.Printf("Посылки клиента №%d:\n", client) //nolint:forbidigo // it's text app
	for _, parcel := range parcels {
		fmt.Printf("Посылка №%d на адрес %s от клиента с идентификатором %d зарегистрирована %s, статус %s\n", //nolint:forbidigo // it's text app
			parcel.ID, parcel.Address, parcel.ClientID, parcel.CreatedAt, parcel.Status)
	}
	fmt.Println() //nolint:forbidigo // it's text app

	return nil
}

func (s ParcelService) NextStatus(number int) error {
	parcel, err := s.store.Get(number)
	if err != nil {
		return fmt.Errorf("get parcel error: %w", err)
	}

	var nextStatus model.ParcelStatus
	switch parcel.Status {
	case model.ParcelStatusRegistered:
		nextStatus = model.ParcelStatusSent
	case model.ParcelStatusSent:
		nextStatus = model.ParcelStatusDelivered
	case model.ParcelStatusDelivered:
		return nil
	}

	fmt.Printf("У посылки №%d новый статус: %s\n", number, nextStatus) //nolint:forbidigo // it's text app
	err = s.store.SetStatus(number, nextStatus)
	if err != nil {
		return fmt.Errorf("set status error: %w", err)
	}

	return nil
}

func (s ParcelService) ChangeAddress(number int, address string) error {
	err := s.store.SetAddress(number, address)
	if err != nil {
		return fmt.Errorf("change address error: %w", err)
	}

	return nil
}

func (s ParcelService) Delete(number int) error {
	err := s.store.Delete(number)
	if err != nil {
		return fmt.Errorf("delete parcel error: %w", err)
	}

	return nil
}

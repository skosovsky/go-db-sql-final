package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/skosovsky/go-db-sql-final.git/pkg/model"
)

type ParcelStore struct {
	*sql.DB
}

func NewParcelStore(pathFile string) (ParcelStore, error) {
	db, err := sql.Open("sqlite", pathFile)
	if err != nil {
		err = fmt.Errorf("db open error: %w", err)

		return ParcelStore{nil}, err
	}

	return ParcelStore{db}, nil
}

func (p ParcelStore) CloseStore() {
	err := p.Close()
	if err != nil {
		log.Println("db close error")
		return
	}
}

func (p ParcelStore) Add(parcel model.Parcel) (int, error) {
	res, err := p.Exec("INSERT INTO parcel (client_id, status, address, created_at) VALUES (:client_id, :status, :address, :created_at)",
		sql.Named("client_id", parcel.ClientID),
		sql.Named("status", parcel.Status),
		sql.Named("address", parcel.Address),
		sql.Named("created_at", parcel.CreatedAt))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("get last insert id error %w: ", err)
		return 0, err
	}

	return int(id), nil
}

func (p ParcelStore) Get(id int) (model.Parcel, error) {
	var parcel model.Parcel

	row := p.QueryRow("SELECT id, client_id, status, address, created_at FROM parcel WHERE id = :id",
		sql.Named("id", id))

	err := row.Scan(&parcel.ID, &parcel.ClientID, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		err = fmt.Errorf("row scan error: %w", err)
		return model.Parcel{}, err
	}

	return parcel, nil
}

func (p ParcelStore) GetByClient(clientID int) ([]model.Parcel, error) {
	rows, err := p.Query("SELECT id, client_id,status,address,created_at FROM parcel WHERE client_id = :client_id", //nolint:sqlclosecheck // it's strange, because row is closed
		sql.Named("client_id", clientID))
	if err != nil {
		err = fmt.Errorf("db query error: %w", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			err = fmt.Errorf("db rows close error: %w", err)
			log.Println(err)
		}
	}(rows)

	var parcels []model.Parcel
	for rows.Next() {
		var parcel model.Parcel

		err = rows.Scan(&parcel.ID, &parcel.ClientID, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			err = fmt.Errorf("rows scan error: %w", err)
			return nil, err
		}

		parcels = append(parcels, parcel)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows next error: %w", err)
		return nil, err
	}

	return parcels, nil
}

func (p ParcelStore) SetStatus(id int, status model.ParcelStatus) error {
	_, err := p.Exec("UPDATE parcel SET status = :status WHERE id = :id",
		sql.Named("id", id),
		sql.Named("status", status))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return err
	}

	return nil
}

func (p ParcelStore) SetAddress(id int, address string) error {
	_, err := p.Exec("UPDATE parcel SET address = :address WHERE id = :id AND status == :status",
		sql.Named("id", id),
		sql.Named("status", model.ParcelStatusRegistered),
		sql.Named("address", address))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return err
	}

	return nil
}

func (p ParcelStore) Delete(id int) error {
	_, err := p.Exec("DELETE FROM parcel WHERE id = :id AND status == :status",
		sql.Named("id", id),
		sql.Named("status", model.ParcelStatusRegistered))
	if err != nil {
		err = fmt.Errorf("db exec error: %w", err)
		return err
	}

	return nil
}

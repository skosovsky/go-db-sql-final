package service_test

import (
	"database/sql"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/skosovsky/go-db-sql-final.git/pkg/model"
	"github.com/skosovsky/go-db-sql-final.git/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// getTestParcel возвращает тестовую посылку.
func getTestParcel() model.Parcel {
	return model.Parcel{
		ID:        0,
		ClientID:  1000,
		Status:    model.ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки.
func TestAddGetDelete(t *testing.T) {
	db, err := store.NewParcelStore()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *store.ParcelStore) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	parcel := getTestParcel()

	parcel.ID, err = db.Add(parcel)
	require.True(t,
		assert.NoError(t, err),
		assert.NotEmpty(t, parcel.ID))

	parcelDataVerification, err := db.Get(parcel.ID)
	require.NoError(t, err)

	assert.EqualValues(t, parcelDataVerification, parcel)

	err = db.Delete(parcel.ID)
	require.NoError(t, err)

	_, err = db.Get(parcel.ID)
	require.Error(t, err)

	expectedErr := sql.ErrNoRows.Error()
	actualErr := err
	assert.ErrorContains(t, actualErr, expectedErr)
}

// TestSetAddress проверяет обновление адреса.
func TestSetAddress(t *testing.T) {
	db, err := store.NewParcelStore()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *store.ParcelStore) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	parcel := getTestParcel()

	parcel.ID, err = db.Add(parcel)
	require.True(t,
		assert.NoError(t, err),
		assert.NotEmpty(t, parcel.ID))

	newAddress := "new test address"
	err = db.SetAddress(parcel.ID, newAddress)
	require.NoError(t, err)

	parcelDataVerification, err := db.Get(parcel.ID)
	require.NoError(t, err)

	expectedAddress := parcelDataVerification.Address
	actualAddress := newAddress

	assert.EqualValues(t, expectedAddress, actualAddress)
}

// TestSetStatus проверяет обновление статуса.
func TestSetStatus(t *testing.T) {
	db, err := store.NewParcelStore()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *store.ParcelStore) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	parcel := getTestParcel()

	parcel.ID, err = db.Add(parcel)
	require.True(t,
		assert.NoError(t, err),
		assert.NotEmpty(t, parcel.ID))

	err = db.SetStatus(parcel.ID, model.ParcelStatusSent)
	require.NoError(t, err)

	parcelDataVerification, err := db.Get(parcel.ID)
	require.NoError(t, err)

	expectedStatus := parcelDataVerification.Status
	actualStatus := model.ParcelStatusSent

	assert.EqualValues(t, expectedStatus, actualStatus)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента.
func TestGetByClient(t *testing.T) {
	var (
		randSource = rand.NewSource(time.Now().UnixNano())
		randRange  = rand.New(randSource)
	)

	db, err := store.NewParcelStore()
	if err != nil {
		log.Println(err)
		return
	}

	defer func(db *store.ParcelStore) {
		err = db.Close()
		if err != nil {
			log.Println("db close error")
		}
	}(&db)

	parcels := []model.Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]model.Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].ClientID = client
	parcels[1].ClientID = client
	parcels[2].ClientID = client

	for idx := range parcels {
		var id int
		id, err = db.Add(parcels[idx])
		parcels[idx].ID = id

		require.True(t,
			assert.NoError(t, err),
			assert.NotEmpty(t, parcels[idx].ID))

		parcels[idx].ID = id

		parcelMap[id] = parcels[idx]
	}

	storedParcels, err := db.GetByClient(client)
	require.NoError(t, err)

	expectedCountParcels := len(storedParcels)
	actualCountParcels := len(parcels)

	assert.EqualValues(t, expectedCountParcels, actualCountParcels)
	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.ID]
		actualParcel := parcel
		require.True(t, ok)
		assert.EqualValues(t, expectedParcel, actualParcel)
	}
}

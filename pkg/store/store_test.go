package store_test

import (
	"database/sql"
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
	db, err := store.NewParcelStore("../../data/tracker.db")
	require.NoError(t, err)
	defer db.CloseStore()

	parcel := getTestParcel()

	parcel.ID, err = db.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.ID)

	parcelDataVerification, err := db.Get(parcel.ID)
	require.NoError(t, err)

	expectedParcel := parcel
	actualParcel := parcelDataVerification
	assert.EqualValues(t, expectedParcel, actualParcel)

	err = db.Delete(parcel.ID)
	require.NoError(t, err)

	_, err = db.Get(parcel.ID)
	expectedErr := sql.ErrNoRows.Error()
	actualErr := err
	require.ErrorContains(t, actualErr, expectedErr)
}

// TestSetAddress проверяет обновление адреса.
func TestSetAddress(t *testing.T) {
	db, err := store.NewParcelStore("../../data/tracker.db")
	require.NoError(t, err)
	defer db.CloseStore()

	parcel := getTestParcel()

	parcel.ID, err = db.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.ID)

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
	db, err := store.NewParcelStore("../../data/tracker.db")
	require.NoError(t, err)
	defer db.CloseStore()

	parcel := getTestParcel()

	parcel.ID, err = db.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.ID)

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

	db, err := store.NewParcelStore("../../data/tracker.db")
	require.NoError(t, err)
	defer db.CloseStore()

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

		require.NoError(t, err)
		require.NotEmpty(t, parcels[idx].ID)

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

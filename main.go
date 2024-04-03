package main //nolint:cyclop // it's example code

import (
	"log"

	"github.com/skosovsky/go-db-sql-final.git/pkg/service"
	"github.com/skosovsky/go-db-sql-final.git/pkg/store"
	_ "modernc.org/sqlite"
)

func main() {
	// Подключение к БД
	db, err := store.NewParcelStore("data/tracker.db")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.CloseStore()

	app := service.NewParcelService(db)

	// Регистрация посылки
	client := 1
	address := "Псков, д. Пушкина, ул. Колотушкина, д. 5"
	parcel, err := app.Register(client, address)
	if err != nil {
		log.Panicln(err)
		return
	}

	// Изменение адреса
	newAddress := "Саратов, д. Верхние Зори, ул. Козлова, д. 25"
	err = app.ChangeAddress(parcel.ID, newAddress)
	if err != nil {
		log.Println(err)
		return
	}

	// Изменение статуса
	err = app.NextStatus(parcel.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Вывод посылок клиента
	err = app.PrintClientParcels(client)
	if err != nil {
		log.Println(err)
		return
	}

	// Попытка удаления отправленной посылки
	err = app.Delete(parcel.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Вывод посылок клиента
	// предыдущая посылка не должна удалиться, т.к. её статус НЕ «зарегистрирована»
	err = app.PrintClientParcels(client)
	if err != nil {
		log.Println(err)
		return
	}

	// Регистрация новой посылки
	parcel, err = app.Register(client, address)
	if err != nil {
		log.Println(err)
		return
	}

	// Удаление новой посылки
	err = app.Delete(parcel.ID)
	if err != nil {
		log.Println(err)
		return
	}

	// Вывод посылок клиента
	// здесь не должно быть последней посылки, т.к. она должна была успешно удалиться
	err = app.PrintClientParcels(client)
	if err != nil {
		log.Println(err)
		return
	}
}

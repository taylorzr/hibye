package main

import(
	"log"

	"github.com/taylorzr/hibye/compare"
	"github.com/taylorzr/hibye/hipchat"
	"github.com/taylorzr/hibye/notify"
	"github.com/taylorzr/hibye/storage"
)

const old_users_path = "users.csv"
// const old_users_path = "test_old_users.csv"

func check() {
	oldUsers, err := storage.Read(old_users_path)
	if err != nil {
		log.Fatal(err)
	}

	newUsers, err := hipchat.GetAllUsers()
	// newUsers, err := storage.Read("test_new_users.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer storage.Write(old_users_path, newUsers)

	firedUsers := compare.FindFired(oldUsers, newUsers)
	hiredUsers := compare.FindHired(oldUsers, newUsers)

	err = notify.Notify(firedUsers, hiredUsers)
	if err != nil {
		log.Fatal(err)
	}
}

func initialize() {
	newUsers, err := hipchat.GetAllUsers()

	if err != nil {
		log.Fatal(err)
	}

	err = storage.Write(old_users_path, newUsers)

	if err != nil {
		log.Fatal(err)
	}

	err = hipchat.SendMessage("HiBye initialized, will report on next run", hipchat.Yellow)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if storage.Exists(old_users_path) {
		check()
	} else {
		initialize()
	}
}

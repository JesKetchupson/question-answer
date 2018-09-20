package database

import (
	"holy-war-web/api/helpers"
	. "holy-war-web/api/models"
	"strconv"
)

var db, err = helpers.GetDb()

func Migrate() {
	if (!db.HasTable(User{})) {
		db.AutoMigrate(User{})
	}
	if (!db.HasTable(Question{})) {
		db.AutoMigrate(Question{})
	}
	if (!db.HasTable(Object{})) {
		db.AutoMigrate(Object{})

	}
	if (!db.HasTable(Question{})) {
		db.AutoMigrate(User{})

	}
	if (!db.HasTable(Category{})) {
		db.AutoMigrate(Category{})

	}
	if (!db.HasTable(Answer{})) {
		db.AutoMigrate(Answer{})

	}

	println("migrated")
}

func Seed() {
	if (db.HasTable(User{})) {
		for i := 1; i < 11; i++ {
			seed := User{
				Email:    "email" + strconv.Itoa(i),
				Password: "pass" + strconv.Itoa(i),
			}
			db.Create(&seed)

		}
		seed2 := Object{
			CategoryID: 1,
			Name:       "asd",
			Image:      "asd",
		}
		db.Create(&seed2)
		seed3 := Question{
			FirstObjectID:  uint(1),
			SecondObjectID: uint(1),
		}
		db.Create(&seed3)
	}

}

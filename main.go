package main

import (
	"net/http"
		"strconv"
				"os"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
		"github.com/gorilla/handlers"
	. "awesomeProject/api/models"
	"awesomeProject/api/helpers"
	"awesomeProject/api/router"
	)

//имя фамилия, отображаемое имя, аффка
//отдельная таблица связанная с юзером показывающая что и где он выбрал
//Модель вопросов Question{A,B}
var db, err  = helpers.GetDb()
func main() {
	//close DB after Main function's end
	defer db.Close()
	migrate()
	println("Server started on port 8080")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router.Router()))
}

func migrate() {
	db.DropTable(User{})
	if (!db.HasTable(User{})) {
		println("migrated")
		db.AutoMigrate(User{})
		for i := 1; i < 11; i++ {
			seed := User{
				Email:    "email" + strconv.Itoa(i),
				Password: "pass" + strconv.Itoa(i),
			}
			db.Create(&seed)
		}
	}
}


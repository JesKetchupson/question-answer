package main

import (
	"net/http"
	"os"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/gorilla/handlers"
	"awesomeProject/api/helpers"
	"awesomeProject/api/router"
	"awesomeProject/api/database"
)

//имя фамилия, отображаемое имя, аффка
//отдельная таблица связанная с юзером показывающая что и где он выбрал
//Модель вопросов Question{A,B}
var db, err = helpers.GetDb()

func main() {
	//close DB after Main function's end

	database.Migrate()
	database.Seed()
	println("Server started on port 8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router.Router()))
	db.Close()
}

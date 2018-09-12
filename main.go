package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"awesomeProject/api/helpers"
	"net/http"
	"github.com/gorilla/handlers"
	"os"
	"awesomeProject/api/router"
)

//имя фамилия, отображаемое имя, аффка
//отдельная таблица связанная с юзером показывающая что и где он выбрал
//Модель вопросов Question{A,B}
var db, err = helpers.GetDb()

func main() {
	//close DB after Main function's end
	if err != nil {
		panic(err)
	}

	helpers.InitEnvVal("env")

	//database.Migrate()
	//database.Seed()
	//println("Server started on port 8080")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router.Router()))

	defer db.Close()
}

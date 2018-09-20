package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"holy-war-web/api/helpers"
	"holy-war-web/api/router"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

//имя фамилия, отображаемое имя, аффка
//отдельная таблица связанная с юзером показывающая что и где он выбрал
//Модель вопросов Question{A,B}
var db, err = helpers.GetDb()

func init() {
	helpers.InitEnvVal("env")
	//database.Migrate()
	//database.Seed()
}

func main() {
	//close DB after Main function's end
	if err != nil {
		panic(err)
	}

	println("Server started on port 8080")

	http.ListenAndServe(":"+os.Getenv("Port"), handlers.LoggingHandler(os.Stdout, router.Router()))

	defer db.Close()
}

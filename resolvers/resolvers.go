package resolvers
//Unused package
//
import (
	"awesomeProject/api/helpers"
	. "awesomeProject/api/models"
	"github.com/graphql-go/graphql"
)

var db, err = helpers.GetDb()

func Resolver(p graphql.ResolveParams) (interface{}, error) {
	var req User
	db.Find(&req)
	return req, err
}

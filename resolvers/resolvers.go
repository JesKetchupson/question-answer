package resolvers
//Unused package
//
import (
	"holy-war-web/api/helpers"
	. "holy-war-web/api/models"
	"github.com/graphql-go/graphql"
)

var db, err = helpers.GetDb()

func Resolver(p graphql.ResolveParams) (interface{}, error) {
	var req User
	db.Find(&req)
	return req, err
}

package controllers

import (
	"awesomeProject/api/helpers"
	. "awesomeProject/api/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
	"github.com/gorilla/websocket"
	)

var db, _ = helpers.GetDb()

func Registration(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	db.NewRecord(req)
	db.Create(&req)
	req.GenerateAccess(req.Email, req.Password)

	req.GenerateRefresh(req.Email, req.Password)

	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)

	isTrue := db.Where("email=?", req.Email).Where("password=?", req.Password)
	if !isTrue.RecordNotFound() {
		req.GenerateAccess(req.Email, req.Password)

		req.GenerateRefresh(req.Email, req.Password)

		err := json.NewEncoder(w).Encode(req)
		if err != nil {
			w.Write([]byte("Something bad happens"))
		}
	}
}
func Read(w http.ResponseWriter, r *http.Request) {
	var req []User
	db.Find(&req)

	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Create(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	db.NewRecord(req)
	db.Create(&req)
	json.NewEncoder(w).Encode(req)
}
func Update(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	db.Model(User{}).Where("id=?", req.ID).Update(req)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}
func Del(w http.ResponseWriter, r *http.Request) {
	req := helpers.GetDecodedJson(r)
	go helpers.DeleteAfter(&req, 3)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.Write([]byte("Something bad happens"))
	}
}

var questionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Question",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"first_object_id": &graphql.Field{
				Type: graphql.Int,
			},
			"second_object_id": &graphql.Field{
				Type: graphql.Int,
			},
			"first_object": &graphql.Field{
				Type: object,
			},
			"second_object": &graphql.Field{
				Type: object,
			},
			"comment": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var object = graphql.NewObject(graphql.ObjectConfig{
	Name: "Object",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"category": &graphql.Field{
			Type: category,
		},
		"category_id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"image": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var category = graphql.NewObject(graphql.ObjectConfig{
	Name: "Category",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var questions []Question

//var db, err =helpers.GetDb()
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:8080/graphql?query={questions(id:1){first_object,second_object}}
			*/
			"questions": &graphql.Field{
				Type:        questionType,
				Description: "Get question by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"]
					if ok {
						ad := db.Find(&Question{ID: uint(id.(int))}).Value.(*Question)

						a2 := db.Find(&Object{ID: ad.SecondObjectID}).Value.(*Object)
						b2 := db.Find(&Category{ID: a2.CategoryID}).Value.(*Category)
						a2.Category = *b2
						ad.FirstObject = *a2

						so := db.Find(&Object{ID: ad.SecondObjectID}).Value.(*Object)
						cat := db.Find(&Category{ID: so.CategoryID}).Value.(*Category)
						so.Category = *cat
						ad.SecondObject = *so
						return ad, nil
					}
					return nil, errors.New("not found")
				},
			},
			/* Get (read) product list
			   http://localhost:8080/graphql?query={list{id}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(questionType),
				Description: "Get questions list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return questions, nil
				},
			},
			"getRandomQuestion": &graphql.Field{
				Type:        questionType,
				Description: "Get random question for user",
				Args: graphql.FieldConfigArgument{
					"user_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					UserID, _ := p.Args["user_id"]
					var Answers []Answer
					db.Where("user_id=?", UserID).Find(&Answers)
					var ListQuestions []Question
					db.Find(&ListQuestions)
					var i = 0
					for _, QuestionA := range ListQuestions {
						for _, Answer := range Answers {
							if QuestionA.ID != Answer.QuestionID {
								i++
							}
							if i == len(Answers) {
								return QuestionA, nil
							}
						}
						i = 0
					}
					return nil, errors.New("no questions for you")
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new product item
		http://localhost:8080/product?query=mutation+_{create(name:"asd"){id,name,info,price}}
		*/
		"createCategory": &graphql.Field{
			Type:        category,
			Description: "Create new Question",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name, nameOk := params.Args["name"].(string)
				NewCategory := Category{}
				if nameOk {
					NewCategory.Name = name
					db.Create(&NewCategory)
				}
				return NewCategory, nil
			},
		},

		"createQuestion": &graphql.Field{
			Type:        questionType,
			Description: "Create new Question",
			Args: graphql.FieldConfigArgument{
				"first_object_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"second_object_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"comment": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				fo, foOk := params.Args["first_object_id"].(int)
				so, soOk := params.Args["second_object_id"].(int)
				comment, commentOk := params.Args["comment"].(string)
				NewQuestion := Question{}
				if foOk {
					NewQuestion.FirstObjectID = uint(fo)
				}
				if soOk {
					NewQuestion.SecondObjectID = uint(so)
				}
				if commentOk {
					NewQuestion.Comment = comment
				}
				db.Create(&NewQuestion)
				return NewQuestion, nil
			},
		},


		"createObject": &graphql.Field{
			Type:        object,
			Description: "Create new Object",
			Args: graphql.FieldConfigArgument{
				"category_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name, nameOk := params.Args["name"].(string)
				catID, catOk := params.Args["category_id"]

				image, imageOk := params.Args["image"].(string)
				NewObject := Object{}
				if nameOk {
					NewObject.Name = name
				}
				if catOk {
					NewObject.CategoryID = uint(catID.(int))
				}
				if imageOk {
					NewObject.Image = image
				}
				db.Create(&NewObject)
				return NewObject, nil
			},
		},


		"addAnswer": &graphql.Field{
			Type:        object,
			Description: "Add answer on User",
			Args: graphql.FieldConfigArgument{
				"user_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"question_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"object_id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				user, userOk := params.Args["user_id"]
				questionID, questionOk := params.Args["question_id"]
				objectID, objectOk := params.Args["object_id"]

				NewAnswer := Answer{}
				if userOk {
					NewAnswer.UserID = uint(user.(int))
				}
				if questionOk {
					NewAnswer.QuestionID = uint(questionID.(int))
				}
				if objectOk {
					NewAnswer.ObjectID = uint(objectID.(int))
				}
				db.Create(&NewAnswer)
				return NewAnswer, nil
			},
		},


		/* Update product by id
		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
		*/
		//"updateObject": &graphql.Field{
		//	Type:        questionType,
		//	Description: "Update product by id",
		//	Args: graphql.FieldConfigArgument{
		//
		//		"first_object": &graphql.ArgumentConfig{
		//			Type: graphql.Int,
		//		},
		//		"second_object": &graphql.ArgumentConfig{
		//			Type: graphql.String,
		//		},
		//		"comment": &graphql.ArgumentConfig{
		//			Type: graphql.Float,
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		id, _ := params.Args["id"].(int)
		//		name, nameOk := params.Args["name"].(string)
		//		info, infoOk := params.Args["info"].(string)
		//		price, priceOk := params.Args["comment"].(float64)
		//
		//		return product, nil
		//	},
		//},

		/* Delete product by id
		   http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
		*/
		//"delete": &graphql.Field{
		//	Type:        questionType,
		//	Description: "Delete product by id",
		//	Args: graphql.FieldConfigArgument{
		//		"id": &graphql.ArgumentConfig{
		//			Type: graphql.NewNonNull(graphql.Int),
		//		},
		//	},
		//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		//		id, _ := params.Args["id"].(int)
		//		product := Product{}
		//		for i, p := range products {
		//			if int64(id) == p.ID {
		//				product = products[i]
		//				// Remove from product list
		//				products = append(products[:i], products[i+1:]...)
		//			}
		//		}
		//
		//		return product, nil
		//	},
		//},
	},
})

func GraphQl(w http.ResponseWriter, r *http.Request) {
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
	result := executeQuery(r.URL.Query().Get("query"), schema)
	json.NewEncoder(w).Encode(result)
}

var upgrader = websocket.Upgrader{}

// Not best temporary solution
var ConChan []*websocket.Conn

func GraphQlWs(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
	check := true
	for _, z := range ConChan {
		if z == conn {
			check = false
		}
	}
	if check {
		ConChan = append(ConChan, conn)
	}

	go func(conn *websocket.Conn) {
		for {
			Mt, p, err := conn.ReadMessage()
			result := executeQuery(string(p), schema)
			a, _ := json.Marshal(result)

			if err != nil {
				conn.Close()
				return
			}

			for _, b := range ConChan {
				b.WriteMessage(Mt, a)
			}

			if err != nil {
				conn.Close()
				return
			}
		}
	}(conn)
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

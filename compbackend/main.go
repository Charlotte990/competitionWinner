package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

// Entry contains the rules for an Entry type object
type Entry struct {
	ID       int64  `json:"id" graphql:"id"`
	UserName string `json:"username" graphql:"username"`
}

// EntryList is a slice of lots of entries
var EntryList []Entry

// initialize some data into the variable above
func init() {
	entry1 := Entry{ID: 1, UserName: "Charlotte"}
	entry2 := Entry{ID: 2, UserName: "Jerry"}
	entry3 := Entry{ID: 3, UserName: "Harry"}
	entry4 := Entry{ID: 3, UserName: "Ben"}
	EntryList = append(EntryList, entry1, entry2, entry3, entry4)

	rand.Seed(time.Now().UnixNano())
}

// define the rules for an entry in graphql
var entryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Entry",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// define the rules of a query in graphql
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		// returns all entries
		"entryList": &graphql.Field{
			Type:        graphql.NewList(entryType),
			Description: "List of entries",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return EntryList, nil
			},
		},
		// returns 1 random entry
		"randomSelection": &graphql.Field{
			Type:        entryType,
			Description: "Random selection",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Find random winner
				rand.Seed(time.Now().Unix())
				message := EntryList[rand.Intn(len(EntryList))]
				return message, nil
			},
		},
	},
})

// schema is required to pass into the query when it is executed. Uses the rootQuery as template/rule/validator
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

// generates the queries defined in the rootQuery
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

type entriesStruct struct {
	Query string `json:"query"`
}

func entries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
	w.Header().Set("Access-Control-Request-Headers", "Origin, Content-Type, X-Auth-Token, Access-Control-Allow-Methods")
	if r.Method == http.MethodOptions {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var entries entriesStruct
	err = json.Unmarshal(body, &entries)
	if err != nil {
		panic(err)
	}
	fmt.Println(entries.Query)
	result := executeQuery(entries.Query, schema)
	fmt.Println(result)
	json.NewEncoder(w).Encode(result)
}

func main() {
	// create handler for specific url endpoint - this displays the json encoded results of the query sent in the URL
	// i.e. http://localhost:8080/allentries?query={entryList{id,username}}
	http.HandleFunc("/entries", entries)

	// run the server
	fmt.Println("Now server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}

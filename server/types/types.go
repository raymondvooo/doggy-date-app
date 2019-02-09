package types

import (
	"github.com/graph-gophers/graphql-go"
)

type User struct {
	ID    graphql.ID
	Name  string
	Email string
	Dogs  []graphql.ID
}

type Dog struct {
	ID    graphql.ID
	Name  string
	Age   int32
	Breed string
	Owner graphql.ID
}

type Date struct {
	ID          graphql.ID
	Date        string
	Description string
	Dogs        []graphql.ID
	Location    string
	User        graphql.ID
}

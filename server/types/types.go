package types

import (
	"github.com/graph-gophers/graphql-go"
)

type User struct {
	ID              graphql.ID
	Name            string
	Email           string
	Dogs            []graphql.ID
	ProfileImageURL string
	JoinDate        graphql.Time
}

type Dog struct {
	ID              graphql.ID
	Name            string
	Age             int32
	Breed           string
	Owner           graphql.ID
	ProfileImageURL string
}

type Date struct {
	ID          graphql.ID
	Date        graphql.Time
	Description string
	Dogs        []graphql.ID
	Location    string
	User        graphql.ID
}

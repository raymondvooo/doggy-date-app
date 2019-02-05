package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
)

// Root holds a pointer to a graphql object
type Root struct {
	Query *graphql.Object
}

func createQuerySchema(parameter string, gqlType *graphql.Object, resolverFunction graphql.FieldResolveFn) *graphql.Field {
	hello := &graphql.Field{
		Type:    graphql.NewList(gqlType),
		Args:    createConfigArgument(parameter),
		Resolve: resolverFunction,
	}
	return hello
}

func createConfigArgument(parameter string) graphql.FieldConfigArgument {
	configArgument := make(map[string]*graphql.ArgumentConfig)
	configArgument[parameter] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}
	return configArgument
}

// NewRoot returns base query type. This is where we add all the base queries
func NewRoot(db *postgres.Db) *Root {
	// Create a resolver holding our databse. Resolver can be found in resolvers.go
	resolver := Resolver{db: db}

	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called name

	root := Root{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					"users": createQuerySchema("name", User, resolver.UserResolver),
					"dogs":  createQuerySchema("name", Dog, resolver.DogResolver),
				},
			},
		),
	}
	return &root
}

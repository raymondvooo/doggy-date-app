package gql

import (
	// "fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
	"github.com/raymondvooo/doggy-date-app/server/types"
)

// Resolver has a reference database
type Resolver struct {
	Db *postgres.Db
}

// User graphql query
func (r *Resolver) User(args struct{ ID graphql.ID }) (*UserResolver, error) {
	user, err := r.Db.GetUsersById(args.ID)
	if err != nil {
		return &UserResolver{nil, r.Db}, err
	}
	data := &UserResolver{&user, r.Db}
	return data, nil
}

// Dog graphql query
func (r *Resolver) Dog(args struct{ ID graphql.ID }) (*DogResolver, error) {
	dog, user, err := r.Db.GetDogById(args.ID)
	if err != nil {
		return &DogResolver{&types.Dog{}, r.Db, &user}, err
	}
	data := &DogResolver{&dog, r.Db, &user}
	return data, nil
}

// CreateUser graphql mutation
func (r *Resolver) CreateUser(args *struct {
	ID       graphql.ID
	Name     string
	Email    string
	DogID    graphql.ID
	DogName  string
	DogAge   int32
	DogBreed string
}) (*UserResolver, error) {
	var err error
	if uidExists, err := r.Db.CheckIDExists("users", args.ID); !uidExists && err != nil { // check existing userID
		if didExists, err := r.Db.CheckIDExists("dogs", args.DogID); !didExists && err != nil { //check existing dogID
			if dogID, err := r.Db.InsertDog(args.DogID, args.DogName, args.DogAge, args.DogBreed, args.ID); err == nil { // add user data into db
				if u, err := r.Db.InsertUser(args.ID, args.Name, args.Email, dogID); err == nil { // add dog data into db
					return &UserResolver{&u, r.Db}, nil
				}
				return &UserResolver{&types.User{}, r.Db}, err
			}
			r.Db.DeleteDog(args.DogID)
			return &UserResolver{&types.User{}, r.Db}, err
		}
		return &UserResolver{&types.User{}, r.Db}, err
	}
	return &UserResolver{&types.User{}, r.Db}, err
}

// UserResolver structure to resolve a user object type to graphql
type UserResolver struct {
	u  *types.User
	db *postgres.Db
}

// DogResolver structure to resolve a dog object type to graphql
type DogResolver struct {
	d  *types.Dog
	db *postgres.Db
	o  *types.User
}

// ID function required by graphql to return user's ID
func (r *UserResolver) ID() graphql.ID {
	return r.u.ID
}

// Name function required by graphql to return user's name
func (r *UserResolver) Name() *string {
	return &r.u.Name
}

// Email function required by graphql to return user's email
func (r *UserResolver) Email() *string {
	return &r.u.Email
}

// Dogs function required by graphql to return user's Dog array object
func (r *UserResolver) Dogs() *[]*DogResolver {
	if dogs, err := r.db.GetDogsByArray(r.u.Dogs); err == nil {
		var dr []*DogResolver
		for _, d := range dogs {
			dr = append(dr, &DogResolver{d: &d, db: r.db, o: r.u})
		}
		return &dr
	}
	return &[]*DogResolver{{&types.Dog{}, r.db, &types.User{}}}
}

// ID function required by graphql to return dogs's ID
func (r *DogResolver) ID() graphql.ID {
	return r.d.ID
}

// Name function required by graphql to return dogs's name
func (r *DogResolver) Name() *string {
	return &r.d.Name
}

// Age function required by graphql to return user's Age
func (r *DogResolver) Age() *int32 {
	return &r.d.Age
}

// Breed function required by graphql to return dogs's Age
func (r *DogResolver) Breed() *string {
	return &r.d.Breed
}

// Owner function required by graphql to return dogs's User object
func (r *DogResolver) Owner() *UserResolver {
	if u := r.o; u != nil {
		return &UserResolver{u: u, db: r.db}
	}
	return nil
}

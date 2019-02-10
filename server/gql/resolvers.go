package gql

import (
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
	"github.com/raymondvooo/doggy-date-app/server/types"
)

// Resolver has a reference database
type Resolver struct {
	Db *postgres.Db
}

// UserResolver structure to resolve a User object type to graphql
type UserResolver struct {
	u  *types.User
	Db *postgres.Db
}

// DogResolver structure to resolve a Dog object type to graphql
type DogResolver struct {
	d  *types.Dog
	Db *postgres.Db
	o  *types.User
}

// DoggyDateResolver structure to resolve a DoggyDate object type to graphql
type DoggyDateResolver struct {
	date *types.Date
	Db   *postgres.Db
	u    *types.User
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
	if emailExists, err := r.Db.CheckEmailExists(args.Email); emailExists && err == nil { // check existing email
		return &UserResolver{&types.User{}, r.Db}, fmt.Errorf("Email %s is already taken exists in database", args.Email)
	}
	if uidExists, err := r.Db.CheckIDExists("users", args.ID); uidExists && err == nil { // check existing userID
		return &UserResolver{&types.User{}, r.Db}, fmt.Errorf("user ID: %s already exists in database", args.ID)
	}
	if didExists, err := r.Db.CheckIDExists("dogs", args.DogID); didExists && err == nil { //check existing dogID
		return &UserResolver{&types.User{}, r.Db}, fmt.Errorf("dog ID: %s already exists in database", args.ID)
	}
	if dogID, err := r.Db.InsertDog(args.DogID, args.DogName, args.DogAge, args.DogBreed, args.ID); err == nil { // add user data into db
		if u, err := r.Db.InsertUser(args.ID, args.Name, args.Email, dogID); err == nil { // add dog data into db
			return &UserResolver{&u, r.Db}, nil
		} else {
			r.Db.DeleteDog(args.DogID) //error adding user to database, so delete dog associated
			return &UserResolver{&types.User{}, r.Db}, err
		}
	}
	return &UserResolver{&types.User{}, r.Db}, err

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
	if dogs, err := r.Db.GetDogsByArray(r.u.Dogs); err == nil {
		var dr []*DogResolver
		for _, d := range dogs {
			dr = append(dr, &DogResolver{d: &d, Db: r.Db, o: r.u})
		}
		return &dr
	}
	return &[]*DogResolver{{&types.Dog{}, r.Db, &types.User{}}}
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
		return &UserResolver{u: u, Db: r.Db}
	}
	return nil
}

// PlanDate graphql mutation
func (r *Resolver) PlanDate(args *struct {
	ID          graphql.ID
	Date        string
	Description string
	Dogs        []graphql.ID
	Location    string
	User        graphql.ID
}) (*DoggyDateResolver, error) {
	if did, err := r.Db.CheckIDExists("doggy_dates", args.ID); !did && err != nil { // false, err
		if date, err := r.Db.InsertDoggyDate(args.ID, args.Date, args.Description, args.Dogs, args.Location, args.User); err == nil {
			user, err := r.Db.GetUsersById(args.User)
			if err != nil {
				return &DoggyDateResolver{&date, r.Db, &types.User{}}, err
			}
			return &DoggyDateResolver{&date, r.Db, &user}, err
		}
		return &DoggyDateResolver{&types.Date{}, r.Db, &types.User{}}, err
	}
	return &DoggyDateResolver{&types.Date{}, r.Db, &types.User{}}, fmt.Errorf("Doggy Date ID: %s already exists in database", args.ID)
}

// ID function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) ID() graphql.ID {
	return r.date.ID
}

// Date function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Date() *string {
	return &r.date.Date
}

// Description function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Description() *string {
	return &r.date.Description
}

// Dogs function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Dogs() *[]*DogResolver {
	if dogs, err := r.Db.GetDogsByArray(r.date.Dogs); err == nil {
		var dr []*DogResolver
		for _, d := range dogs {
			dr = append(dr, &DogResolver{d: &d, Db: r.Db, o: r.u})
		}
		return &dr
	}
	return &[]*DogResolver{{&types.Dog{}, r.Db, &types.User{}}}
}

// Location function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Location() *string {
	return &r.date.Location
}

// User function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) User() *UserResolver {
	return &UserResolver{r.u, r.Db}
}

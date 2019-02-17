package gql

import (
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
	"github.com/raymondvooo/doggy-date-app/server/types"
	uuid "github.com/satori/go.uuid"
	"log"
)

// Resolver has a reference database
type Resolver struct {
	Db *postgres.Db
}

// UserResolver structure to resolve a User object type to graphql
type UserResolver struct {
	u  *types.User
	d  *[]types.Dog
	Db *postgres.Db
}

// DogResolver structure to resolve a Dog object type to graphql
type DogResolver struct {
	d    *types.Dog
	dogs *[]types.Dog
	Db   *postgres.Db
	o    *types.User
}

// DoggyDateResolver structure to resolve a DoggyDate object type to graphql
type DoggyDateResolver struct {
	date   *types.Date
	Db     *postgres.Db
	u      *types.User
	dogMap *map[graphql.ID]types.Dog
}

// User graphql query
func (r *Resolver) User(args struct{ ID graphql.ID }) (*UserResolver, error) {
	// Check if valid UUID
	uid, err := uuid.FromString(string(args.ID))
	if err != nil {
		log.Println(err)
		return &UserResolver{nil, nil, r.Db}, err
	}
	user, dogs, err := r.Db.GetUserByID(uid)
	if err != nil {
		log.Println(err)
		return &UserResolver{nil, nil, r.Db}, err
	}
	data := &UserResolver{&user, &dogs, r.Db}
	log.Println("Resolve: user graphql query")
	return data, nil
}

// Dog graphql query
func (r *Resolver) Dog(args struct{ ID graphql.ID }) (*DogResolver, error) {
	// Check if valid UUID
	did, err := uuid.FromString(string(args.ID))
	if err != nil {
		log.Println(err)
		return &DogResolver{&types.Dog{}, &[]types.Dog{}, r.Db, &types.User{}}, err
	}
	dogs, user, err := r.Db.GetDogByID(did)
	if err != nil {
		log.Println(err)
		return &DogResolver{&types.Dog{}, &[]types.Dog{}, r.Db, &user}, err
	}
	data := &DogResolver{&dogs[0], &dogs, r.Db, &user}
	log.Println("Resolve: dog graphql query")
	return data, nil
}

// LoginUser graphql query
func (r *Resolver) LoginUser(args struct{ Email string }) (*UserResolver, error) {
	user, dogs, err := r.Db.GetUserByEmail(args.Email)
	if err != nil {
		log.Println(err)
		return &UserResolver{nil, nil, r.Db}, err
	}
	data := &UserResolver{&user, &dogs, r.Db}
	return data, nil
}

// CreateUser graphql mutation
func (r *Resolver) CreateUser(args *struct {
	Name                string
	Email               string
	UserProfileImageURL string
	DogName             string
	DogAge              int32
	DogBreed            string
	DogProfileImageURL  string
}) (*UserResolver, error) {
	var err error
	emailExists, err := r.Db.CheckEmailExists(args.Email)
	if !emailExists && err != nil {
		log.Println("Pass: unused email")
		user, dog, err := r.Db.InsertUserDog(args.Name, args.Email, args.UserProfileImageURL, args.DogName, args.DogAge, args.DogBreed, args.DogProfileImageURL)
		if err != nil {
			log.Println(err)
			return &UserResolver{&types.User{}, &[]types.Dog{}, r.Db}, err
		}
		log.Println("Resolve: createUser graphql mutation")
		return &UserResolver{&user, &[]types.Dog{dog}, r.Db}, nil
	}
	log.Printf("Error: Email %s already exists", args.Email)
	return &UserResolver{&types.User{}, &[]types.Dog{}, r.Db}, fmt.Errorf("Error: Email %s already exists", args.Email)
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
	var dogs []*DogResolver
	for i := 0; i < len(*r.d); i++ {
		d := &(*r.d)[i]
		dogs = append(dogs, &DogResolver{d, r.d, r.Db, r.u})
	}
	return &dogs
}

// JoinDate function required by graphql to return user's email
func (r *UserResolver) JoinDate() *graphql.Time {
	return &r.u.JoinDate
}

// ProfileImageURL function required by graphql to return user's email
func (r *UserResolver) ProfileImageURL() *string {
	return &r.u.ProfileImageURL
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
		return &UserResolver{u, r.dogs, r.Db}
	}
	return nil
}

// ProfileImageURL function required by graphql to return user's email
func (r *DogResolver) ProfileImageURL() *string {
	return &r.d.ProfileImageURL
}

// GetDoggyDates function required by graphql query
func (r *Resolver) GetDoggyDates() (*[]*DoggyDateResolver, error) {
	dates, users, dogs, err := r.Db.GetAllDoggyDates()
	if err != nil {
		log.Println(err)
		return &[]*DoggyDateResolver{{&types.Date{}, r.Db, &types.User{}, &map[graphql.ID]types.Dog{}}}, err
	}
	var ddr []*DoggyDateResolver
	for _, v := range dates {
		date := new(types.Date)
		*date = v
		u := users[v.User]
		ddr = append(ddr, &DoggyDateResolver{date, r.Db, &u, &dogs})
	}
	log.Println("Resolve: getDoggyDates graphql query")
	return &ddr, nil
}

// PlanDate graphql mutation
func (r *Resolver) PlanDate(args *struct {
	Date        graphql.Time
	Description string
	Dogs        []graphql.ID
	Location    string
	User        graphql.ID
}) (*DoggyDateResolver, error) {
	date, err := r.Db.InsertDoggyDate(args.Date, args.Description, args.Dogs, args.Location, args.User)
	if err != nil {
		log.Println(err)
		return &DoggyDateResolver{}, err
	}
	dogMap, uMap, _ := r.Db.GetDogsByArray(args.Dogs)
	u := uMap[date.User]
	log.Println("Resolve: planDate graphql mutation")
	return &DoggyDateResolver{&date, r.Db, &u, &dogMap}, err
}

// ID function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) ID() graphql.ID {
	return r.date.ID
}

// Date function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Date() *graphql.Time {
	return &r.date.Date
}

// Description function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Description() *string {
	return &r.date.Description
}

// Dogs function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Dogs() *[]*DogResolver {
	var dogs []*DogResolver
	var doga []types.Dog
	for i := 0; i < len(r.u.Dogs); i++ {
		doga = append(doga, (*r.dogMap)[r.u.Dogs[i]])
	}
	for i := 0; i < len(r.date.Dogs); i++ {
		d := (*r.dogMap)[r.date.Dogs[i]]
		dogs = append(dogs, &DogResolver{&d, &doga, r.Db, r.u})
	}

	return &dogs
}

// Location function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) Location() *string {
	return &r.date.Location
}

// User function required by graphql to return DoggyDates's ID
func (r *DoggyDateResolver) User() *UserResolver {
	var dogs []types.Dog
	for i := 0; i < len(r.u.Dogs); i++ {
		d := (*r.dogMap)[r.u.Dogs[i]]
		dogs = append(dogs, d)
	}
	return &UserResolver{r.u, &dogs, r.Db}
}

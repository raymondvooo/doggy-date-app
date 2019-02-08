package gql

import (
	// "fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
	"github.com/raymondvooo/doggy-date-app/server/types"
)

type Resolver struct {
	Db *postgres.Db
}

// func (r *Resolver) User(args struct{ ID graphql.ID }) (*userResolver, error) {
// 	if user := userData[args.ID]; user != nil {
// 		data := &userResolver{user}
// 		return data, nil
// 	}
// 	return &userResolver{}, nil
// }

func (r *Resolver) User(args struct{ ID graphql.ID }) (*userResolver, error) {
	user, err := r.Db.GetUsersById(args.ID)
	if err != nil {
		return &userResolver{nil, r.Db}, err
	}
	data := &userResolver{&user, r.Db}
	return data, nil
}

func (r *Resolver) Dog(args struct{ ID graphql.ID }) (*dogResolver, error) {
	dog, user, err := r.Db.GetDogById(args.ID)
	if err != nil {
		return &dogResolver{}, err
	}
	data := &dogResolver{&dog, r.Db, &user}
	return data, nil
}

func (r *Resolver) CreateUser(args *struct {
	ID       graphql.ID
	Name     string
	Email    string
	DogID    graphql.ID
	DogName  string
	DogAge   int32
	DogBreed string
}) (*userResolver, error) {
	var err error
	if dogID, err := r.Db.InsertDog(args.DogID, args.DogName, args.DogAge, args.DogBreed, args.ID); err == nil {
		if u, err := r.Db.InsertUser(args.ID, args.Name, args.Email, dogID); err == nil {
			return &userResolver{&u, r.Db}, nil
		}
		return &userResolver{}, err
	}
	return &userResolver{}, err
}

type userResolver struct {
	u  *types.User
	db *postgres.Db
}

type dogResolver struct {
	d  *types.Dog
	db *postgres.Db
	o  *types.User
}

func (r *userResolver) ID() graphql.ID {
	return r.u.ID
}

func (r *userResolver) Name() *string {
	return &r.u.Name
}

func (r *userResolver) Email() *string {
	return &r.u.Email
}

func (r *userResolver) Dogs() *[]*dogResolver {
	if dogs, err := r.db.GetDogsByArray(r.u.Dogs); err == nil {
		var dr []*dogResolver
		for _, d := range dogs {
			dr = append(dr, &dogResolver{d: &d, db: r.db, o: r.u})
		}
		return &dr
	}
	return &[]*dogResolver{}
}

// func resolveDogs(ids []graphql.ID) *[]*dogResolver {
// 	var dogs []*dogResolver
// 	for _, id := range ids {
// 		if d := resolveDog(id); d != nil {
// 			dogs = append(dogs, d)
// 		}
// 	}
// 	return &dogs
// }

// func resolveDog(id graphql.ID) *dogResolver {
// 	if d, ok := dogData[id]; ok {
// 		return &dogResolver{d}
// 	}
// 	return nil
// }

func (r *dogResolver) ID() graphql.ID {
	return r.d.ID
}

func (r *dogResolver) Name() *string {
	return &r.d.Name
}

func (r *dogResolver) Age() *int32 {
	return &r.d.Age
}

func (r *dogResolver) Breed() *string {
	return &r.d.Breed
}

func (r *dogResolver) Owner() *userResolver {
	if u := r.o; u != nil {
		return &userResolver{u: u, db: r.db}
	}
	return nil
}

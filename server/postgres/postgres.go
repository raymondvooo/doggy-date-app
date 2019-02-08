package postgres

import (
	"database/sql"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/types"

	// postgres driver
	"github.com/lib/pq"
)

// Db is our database struct used for interacting with the database
type Db struct {
	*sql.DB
}

// NewConnection makes a new database using the connection string and
// returns it, otherwise returns the error
func NewConnection(connect string) (*Db, error) {
	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	// Check that our connection is good
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}

// GetUsersById is called within our user query for graphql
func (d *Db) GetUsersById(id graphql.ID) (types.User, error) {
	var u types.User
	var qString = "SELECT * FROM users WHERE id=$1"
	var qDog []string
	// Make database query
	err := d.QueryRow(qString, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		pq.Array(&qDog),
	)
	StringToGraphqlID(qDog, &u.Dogs)

	if err == sql.ErrNoRows {
		fmt.Println("GetUsersById Query Err: ", err)
		return u, err
	}
	return u, nil
}

// GetDogById is called within our user query for graphql
func (d *Db) GetDogById(id graphql.ID) (types.Dog, types.User, error) {
	// Prepare query, takes a name argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM dogs WHERE id=$1")
	if err != nil {
		fmt.Println("GetDogById Preperation Err: ", err)
	}
	defer stmt.Close()
	var dog types.Dog
	var u types.User
	// Make query with our stmt, passing in name argument
	err = stmt.QueryRow(id).Scan(
		&dog.ID,
		&dog.Name,
		&dog.Age,
		&dog.Breed,
		&dog.Owner,
	)
	if err == sql.ErrNoRows {
		fmt.Println("GetDogById Query Err: ", err)
		return dog, u, err
	}
	u, err = d.GetUsersById(dog.Owner)

	if err == sql.ErrNoRows {
		fmt.Println("GetUsersById Query Err: ", err)
		return dog, u, err
	}
	return dog, u, nil
}

// GetDogsByArray is called within our dogs query for graphql
func (d *Db) GetDogsByArray(dogIds []graphql.ID) ([]types.Dog, error) {
	// Prepare query, takes a name argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM dogs WHERE id= ANY($1)")
	if err != nil {
		fmt.Println("GetDogsByArray Preperation Err: ", err)
	}
	defer stmt.Close()
	var ds []string
	ds = GraphqlIDToString(dogIds, ds)
	// dogString := "{" + strings.Join(ds, ", ") + "}"
	// Make query with our stmt, passing in name argument
	rows, err := stmt.Query(pq.Array(ds))
	if err != nil {
		fmt.Println("GetDogByName Query Err: ", err)
	}

	// Create User struct for holding each row's data
	var r types.Dog
	// Create slice of Users for our response
	dogs := []types.Dog{}
	// Copy the columns from row into the values pointed at by r (User)
	for rows.Next() {
		err = rows.Scan(
			&r.ID,
			&r.Name,
			&r.Age,
			&r.Breed,
			&r.Owner,
		)
		if err != nil {
			fmt.Println("Error scanning rows: ", err)
			return []types.Dog{}, err
		}
		dogs = append(dogs, r)
	}
	return dogs, nil
}

// InsertUser queries database to insert user row
func (d *Db) InsertUser(id graphql.ID, name string, email string, dogId graphql.ID) (types.User, error) {
	// Prepare query, takes a name argument, protects from sql injection
	var di = []string{string(dogId)}
	stmt, err := d.Prepare("INSERT INTO users VALUES ($1, $2, $3, $4)")
	if err != nil {
		fmt.Println("InsertUser Preperation Err: ", err)
	}
	defer stmt.Close()
	fmt.Println("YEH EXECUTE USER")
	if _, err := stmt.Exec(string(id), name, email, pq.Array(di)); err != nil {
		return types.User{}, err
	}
	return types.User{ID: id, Name: name, Email: email, Dogs: []graphql.ID{dogId}}, nil
}

// InsertDog queries database to insert dog row
func (d *Db) InsertDog(id graphql.ID, name string, age int32, breed string, ownerId graphql.ID) (graphql.ID, error) {
	// Prepare query, takes a name argument, protects from sql injection
	stmt, err := d.Prepare("INSERT INTO dogs VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		fmt.Println("InsertDog Preperation Err: ", err)
	}
	defer stmt.Close()

	fmt.Println("YEH EXECUTE DOG")
	if _, err := stmt.Exec(string(id), name, int64(age), breed, string(ownerId)); err != nil {
		return "", err
	}
	return id, nil
}

// StringToGraphqlID convert string array to graphqlID array
func StringToGraphqlID(s []string, gqlS *[]graphql.ID) {
	for _, id := range s {
		*gqlS = append(*gqlS, graphql.ID(id))
	}
}

// GraphqlIDToString convert graphqlID array to string array
func GraphqlIDToString(gqlS []graphql.ID, s []string) []string {
	for _, id := range gqlS {
		s = append(s, string(id))
	}
	return s
}

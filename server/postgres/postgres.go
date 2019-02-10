package postgres

import (
	"database/sql"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/types"
	uuid "github.com/satori/go.uuid"

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
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM users WHERE id=$1")
	if err != nil {
		fmt.Println("GetUsersById Preparation Err: ", err)
	}
	defer stmt.Close()
	var u types.User
	var qDog []string
	// Make database query
	uid, _ := uuid.FromString(string(id))
	err = stmt.QueryRow(uid).Scan(
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
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM dogs WHERE id=$1")
	if err != nil {
		fmt.Println("GetDogById Preparation Err: ", err)
	}
	defer stmt.Close()
	var dog types.Dog
	var u types.User
	// Make query with our stmt, passing in id argument
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
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM dogs WHERE id= ANY($1)")
	if err != nil {
		fmt.Println("GetDogsByArray Preparation Err: ", err)
	}
	defer stmt.Close()
	var dus []uuid.UUID
	dus = GraphqlIDToUUID(dogIds, dus)
	// Make query with our stmt, passing in id argument
	rows, err := stmt.Query(pq.Array(dus))
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
	// Prepare query, takes arguments, protects from sql injection
	did, _ := uuid.FromString(string(dogId))
	var di = []uuid.UUID{did}
	stmt, err := d.Prepare("INSERT INTO users VALUES ($1, $2, $3, $4)")
	if err != nil {
		fmt.Println("InsertUser Preparation Err: ", err)
	}
	defer stmt.Close()
	uid, _ := uuid.FromString(string(id))
	if _, err := stmt.Exec(uid, name, email, pq.Array(di)); err != nil {
		fmt.Println("InsertUser Execution Err: ", err)
		return types.User{}, err
	}
	return types.User{ID: id, Name: name, Email: email, Dogs: []graphql.ID{dogId}}, nil
}

// InsertDog queries database to insert dog row
func (d *Db) InsertDog(id graphql.ID, name string, age int32, breed string, ownerId graphql.ID) (graphql.ID, error) {
	// Prepare query, takes arguments, protects from sql injection
	stmt, err := d.Prepare("INSERT INTO dogs VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		fmt.Println("InsertDog Preparation Err: ", err)
	}
	defer stmt.Close()
	did, _ := uuid.FromString(string(id))
	uid, _ := uuid.FromString(string(ownerId))
	if _, err := stmt.Exec(did, name, int64(age), breed, uid); err != nil {
		fmt.Println("InsertDog Execution Err: ", err)
		return "", err
	}
	return id, nil
}

// InsertDoggyDate queries database to insert dog row
func (d *Db) InsertDoggyDate(id graphql.ID, date string, description string, dogIds []graphql.ID, location string, user graphql.ID) (types.Date, error) {
	stmt, err := d.Prepare("INSERT INTO doggy_dates VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		fmt.Println("InsertInsertDoggyDateDog Preparation Err: ", err)
	}
	defer stmt.Close()
	var dus []uuid.UUID
	dus = GraphqlIDToUUID(dogIds, dus)
	did, _ := uuid.FromString(string(id))
	uid, _ := uuid.FromString(string(user))
	if _, err := stmt.Exec(did, date, description, pq.Array(dus), location, uid); err != nil {
		fmt.Println("InsertDoggyDate Exec Err: ", err)
		return types.Date{}, err
	}
	return types.Date{ID: id, Date: date, Description: description, Dogs: dogIds, Location: location, User: user}, nil
}

// DeleteDog queries database to insert dog row
func (d *Db) DeleteDog(id graphql.ID) (bool, error) {
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare("DELETE FROM dogs WHERE id=$1")
	if err != nil {
		fmt.Println("DeleteDog Preparation Err: ", err)
	}
	defer stmt.Close()
	did, _ := uuid.FromString(string(id))
	if _, err := stmt.Exec(did); err != nil {
		return true, err
	}
	return false, nil
}

// CheckIDExists queries database if user or dog ID exists
func (d *Db) CheckIDExists(tableType string, id graphql.ID) (bool, error) {
	q := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", tableType)
	fmt.Println(q)
	stmt, err := d.Prepare(q)
	if err != nil {
		fmt.Println("CheckIDExists Preparation Err: ", err)
	}
	var exists string
	uid, _ := uuid.FromString(string(id))
	err = stmt.QueryRow(uid).Scan(&exists)
	fmt.Println(err)
	if err == sql.ErrNoRows {
		fmt.Println("CheckIDExists Query Err: ", err)
		return false, err
	}
	return true, nil
}

// CheckEmailExists queries database if email exists
func (d *Db) CheckEmailExists(email string) (bool, error) {
	stmt, err := d.Prepare("SELECT email FROM users WHERE email=$1")
	if err != nil {
		fmt.Println("CheckEmailExists Preparation Err: ", err)
	}
	var exists string
	err = stmt.QueryRow(email).Scan(&exists)
	if err == sql.ErrNoRows {
		fmt.Println("CheckEmailExists Query Err: ", err)
		return false, err
	}
	return true, nil
}

// StringToGraphqlID convert string array to graphqlID array
func StringToGraphqlID(s []string, gqlS *[]graphql.ID) {
	for _, id := range s {
		*gqlS = append(*gqlS, graphql.ID(id))
	}
}

// GraphqlIDToUUID convert graphqlID array to string array
func GraphqlIDToUUID(gqlS []graphql.ID, u []uuid.UUID) []uuid.UUID {
	for _, id := range gqlS {
		x, _ := uuid.FromString(string(id))
		u = append(u, x)
	}
	return u
}

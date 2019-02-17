package postgres

import (
	"database/sql"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/raymondvooo/doggy-date-app/server/types"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"

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

// GetUserByEmail is called within our user query for graphql
func (d *Db) GetUserByEmail(email string) (types.User, []types.Dog, error) {
	log.Println("Starting: GetUserByEmail Query")
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare(`SELECT
	u.id,
	u.name,
	u.email,
	u.profile_image,
	u.join_date,
	d.id,
	d.name,
	d.age,
	d.breed,
	d.profile_image
	FROM users u INNER JOIN dogs d ON u.id = d.owner
	WHERE u.email = $1
	ORDER BY u.name;`)
	if err != nil {
		log.Println("GetUserByEmail Preparation Error: ", err)
	}
	defer stmt.Close()

	var u types.User
	var dog types.Dog
	var dogs []types.Dog
	// Make database query
	rows, err := stmt.Query(email)
	if err != nil {
		log.Println("GetUserByEmail Query Error: ", err)
		return u, dogs, err
	}
	for rows.Next() {
		var joinDate time.Time
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.ProfileImageURL,
			&joinDate, // readable Time type
			&dog.ID,
			&dog.Name,
			&dog.Age,
			&dog.Breed,
			&dog.ProfileImageURL,
		)
		if err != nil {
			log.Println(" GetUserByEmail error scanning rows: ", err)
			return u, dogs, err
		}
		u.JoinDate = graphql.Time{Time: joinDate} // convert Time to graphql.Time
		dogs = append(dogs, dog)
	}
	log.Println("Success: GetUserByEmail Query")
	return u, dogs, nil
}

// GetUserByID is called within our user query for graphql
func (d *Db) GetUserByID(id uuid.UUID) (types.User, []types.Dog, error) {
	log.Println("Starting: GetUserByID Query")
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare(`SELECT
	u.id,
	u.name,
	u.profile_image,
	u.join_date,
	d.id,
	d.name,
	d.age,
	d.breed,
	d.profile_image
	FROM users u INNER JOIN dogs d ON u.id = d.owner
	WHERE u.id = $1
	ORDER BY u.name;`)
	if err != nil {
		log.Println("GetUserByID Preparation Error: ", err)
	}
	defer stmt.Close()
	var u types.User
	var dog types.Dog
	var dogs []types.Dog
	// Make database query
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("GetUserByID Query Error: ", err)
		return u, dogs, err
	}
	for rows.Next() {
		var joinDate time.Time
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.ProfileImageURL,
			&joinDate, // readable Time type
			&dog.ID,
			&dog.Name,
			&dog.Age,
			&dog.Breed,
			&dog.ProfileImageURL,
		)
		if err != nil {
			log.Println("GetUserByID error scanning rows: ", err)
			return u, dogs, err
		}
		u.JoinDate = graphql.Time{Time: joinDate} // convert Time to graphql.Time
		dogs = append(dogs, dog)
	}
	log.Println("Success: GetUserByID Query")
	return u, dogs, nil
}

// GetDogByID is called within our user query for graphql
func (d *Db) GetDogByID(id uuid.UUID) ([]types.Dog, types.User, error) {
	log.Println("Starting: GetDogByID Query")
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare(`SELECT
	d.id,
	d.name,
	d.age,
	d.breed,
	d.profile_image,
	u.id,
	u.name,
	u.profile_image
	FROM users u INNER JOIN dogs d ON u.id = d.owner
	WHERE u.dogs::text LIKE '%' || $1 || '%'
	ORDER BY d.id::text = $1 DESC;`)
	if err != nil {
		log.Println("GetDogByID Preparation Error: ", err)
	}
	defer stmt.Close()
	var dog types.Dog
	var dogs []types.Dog
	var u types.User
	// Make query with our stmt, passing in id argument
	rows, err := stmt.Query(id)
	if err != nil {
		log.Println("GetDogByID Query Error: ", err)
		return dogs, u, err
	}
	for rows.Next() {
		err = rows.Scan(
			&dog.ID,
			&dog.Name,
			&dog.Age,
			&dog.Breed,
			&dog.ProfileImageURL,
			&u.ID,
			&u.Name,
			&u.ProfileImageURL,
		)
		if err != nil {
			log.Println("GetDogByID error scanning rows: ", err)
			return dogs, u, err
		}
		dogs = append(dogs, dog)
	}
	log.Println("Success: GetDogByID Query")
	return dogs, u, nil
}

// GetDogsByArray is called within our dogs query for graphql
func (d *Db) GetDogsByArray(dogIds []graphql.ID) (map[graphql.ID]types.Dog, map[graphql.ID]types.User, error) {
	log.Println("Starting: GetDogsByArray Query")
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare(`SELECT
		d.id,
		d.name,
		d.age,
		d.breed,
		d.profile_image,
		u.id,
		u.name,
		u.profile_image,
		u.join_date
		FROM users u INNER JOIN dogs d ON u.id = d.owner
		WHERE d.id= ANY($1);`)
	if err != nil {
		log.Println("GetDogsByArray Preparation Error: ", err)
	}
	defer stmt.Close()
	var dog types.Dog
	var dus []uuid.UUID
	dogMap := map[graphql.ID]types.Dog{}
	var u types.User
	uMap := map[graphql.ID]types.User{}
	GraphqlIDToUUID(dogIds, &dus)
	// Make query with our stmt, passing in id argument
	rows, err := stmt.Query(pq.Array(dus))
	if err != nil {
		log.Println("GetDogsByArray Query Error: ", err)
	}
	// Copy the columns from row into the values pointed at by r (User)
	for rows.Next() {
		var joinDate time.Time
		err = rows.Scan(
			&dog.ID,
			&dog.Name,
			&dog.Age,
			&dog.Breed,
			&dog.ProfileImageURL,
			&u.ID,
			&u.Name,
			&u.ProfileImageURL,
			&joinDate, // readable Time type
		)
		if err != nil {
			log.Println("GetDogsByArray error scanning rows: ", err)
			return dogMap, uMap, err
		}
		u.JoinDate = graphql.Time{Time: joinDate}
		dogMap[dog.ID] = dog
		uMap[u.ID] = u
	}
	log.Println("Success: GetDogsByArray Query")
	return dogMap, uMap, nil
}

// GetAllDoggyDates is called within our doggydate query for graphql
func (d *Db) GetAllDoggyDates() (map[graphql.ID]types.Date, map[graphql.ID]types.User, map[graphql.ID]types.Dog, error) {
	log.Println("Starting: GetAllDoggyDates Query")
	// Prepare query, takes a id argument, protects from sql injection
	stmt, err := d.Prepare(`SELECT
	dd.*,
	u.id,
	u.name,
	u.dogs,
	u.profile_image,
	u.join_date,
	d.id,
	d.name,
	d.age,
	d.breed,
	d.profile_image
	FROM doggy_dates dd
		JOIN users u ON dd.user = u.id
		JOIN dogs d ON d.owner = u.id;`)
	if err != nil {
		log.Println("GetAllDoggyDates Preparation Error: ", err)
	}
	defer stmt.Close()

	// Make query with our stmt, passing in id argument
	rows, err := stmt.Query()
	if err != nil {
		log.Println("GetAllDoggyDates Query Error: ", err)
	}

	var date types.Date
	var u types.User
	uMap := map[graphql.ID]types.User{}
	var dog types.Dog
	dogMap := map[graphql.ID]types.Dog{}
	dates := map[graphql.ID]types.Date{}
	// Copy the columns from row into the values pointed at by r (User)
	for rows.Next() {
		var createDate time.Time
		var userJoinDate time.Time
		var dateDogs []string
		var userDogs []string
		err = rows.Scan(
			&date.ID,
			&createDate, // readable Time type
			&date.Description,
			pq.Array(&dateDogs), // readable [] string type
			&date.Location,
			&date.User,
			&u.ID,
			&u.Name,
			pq.Array(&userDogs), // readable [] string type
			&u.ProfileImageURL,
			&userJoinDate, // readable Time type
			&dog.ID,
			&dog.Name,
			&dog.Age,
			&dog.Breed,
			&dog.ProfileImageURL,
		)
		if err != nil {
			log.Println("GetAllDoggyDates error scanning rows: ", err)
			return dates, uMap, dogMap, err
		}
		date.Date = graphql.Time{Time: createDate}    // convert Time to graphql.Time
		u.JoinDate = graphql.Time{Time: userJoinDate} // convert Time to graphql.Time
		date.Dogs = nil                               // reset values before reassigning
		u.Dogs = nil                                  // reset values before reassigning
		StringToGraphqlID(dateDogs, &date.Dogs)
		StringToGraphqlID(userDogs, &u.Dogs)
		dates[date.ID] = date
		dogMap[dog.ID] = dog
		uMap[u.ID] = u
	}
	log.Println("Success: GetAllDoggyDates Query")
	return dates, uMap, dogMap, nil
}

// InsertUserDog queries database to insert user row
func (d *Db) InsertUserDog(name string, email string, uImg string, dname string,
	age int32, breed string, dImg string) (types.User, types.Dog, error) {
	log.Println("Starting: InsertUserDog Execution")
	// Prepare query, takes arguments, protects from sql injection
	stmt, err := d.Prepare(`WITH createAccount AS (
		INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6)
	  ) INSERT INTO dogs VALUES ($7, $8, $9, $10, $11, $12);`)
	if err != nil {
		log.Println("InsertUserDog Preparation Error: ", err)
	}
	defer stmt.Close()
	did, _ := uuid.NewV1()
	var di = []uuid.UUID{did}
	uid, _ := uuid.NewV1() // Generate new uuid
	joinDate := time.Now() // Generate timestamp
	if _, err := stmt.Exec(uid, name, email, pq.Array(di), uImg, joinDate, did, dname, age, breed, uid, dImg); err != nil {
		log.Println("InsertUserDog Execution Error: ", err)
		return types.User{}, types.Dog{}, err
	}
	log.Println("Success: InsertUserDog Execution")
	return types.User{
			ID:              graphql.ID(uid.String()),
			Name:            name,
			Email:           email,
			Dogs:            []graphql.ID{graphql.ID(did.String())},
			ProfileImageURL: uImg, JoinDate: graphql.Time{Time: joinDate}},
		types.Dog{
			ID:              graphql.ID(did.String()),
			Name:            dname,
			Age:             age,
			Breed:           breed,
			Owner:           graphql.ID(uid.String()),
			ProfileImageURL: dImg}, nil
}

// InsertDoggyDate queries database to insert dog row
func (d *Db) InsertDoggyDate(date graphql.Time, description string, dogIds []graphql.ID, location string, user graphql.ID) (types.Date, error) {
	log.Println("Starting: InsertDoggyDate Execution")
	// Prepare query, takes arguments, protects from sql injection
	stmt, err := d.Prepare("INSERT INTO doggy_dates VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Println("InsertDoggyDateDog Preparation Error: ", err)
	}
	defer stmt.Close()

	var dus []uuid.UUID
	GraphqlIDToUUID(dogIds, &dus)
	did, _ := uuid.NewV1()
	uid, _ := uuid.FromString(string(user))
	gDate := date.Local() // convert graphql.Time to golang Time
	if _, err := stmt.Exec(did, gDate, description, pq.Array(dus), location, uid); err != nil {
		log.Println("InsertDoggyDate Execution Error: ", err)
		return types.Date{}, err
	}
	log.Println("Success: InsertDoggyDateDog Execution")
	return types.Date{ID: graphql.ID(did.String()), Date: date, Description: description, Dogs: dogIds, Location: location, User: user}, nil
}

// UpdateProfilePic queries database if email exists
func (d *Db) UpdateProfilePic(tableType string, id graphql.ID, imgURL string) (bool, error) {
	// Prepare query, takes arguments, protects from sql injection
	q := fmt.Sprintf("UPDATE %s SET profile_image=$1 WHERE id=$2", tableType)
	log.Println("Starting: UpdateProfilePic Execution")
	stmt, err := d.Prepare(q)
	if err != nil {
		log.Println("UpdateProfilePic Preparation Error: ", err)
	}
	defer stmt.Close()
	uid, _ := uuid.FromString(string(id))
	if _, err := stmt.Exec(imgURL, uid); err != nil {
		log.Println("UpdateProfilePic Execution Error: ", err)
		return false, err
	}
	log.Println("Success: UpdateProfilePic Execution")
	return true, nil
}

// CheckIDExists queries database if user or dog ID exists
func (d *Db) CheckIDExists(tableType string, id graphql.ID) (bool, error) {
	q := fmt.Sprintf("SELECT id FROM %s WHERE id=$1", tableType)
	log.Println(q)
	stmt, err := d.Prepare(q)
	if err != nil {
		log.Println("CheckIDExists Preparation Error: ", err)
	}
	var exists string
	uid, _ := uuid.FromString(string(id))
	err = stmt.QueryRow(uid).Scan(&exists)
	log.Println(err)
	if err == sql.ErrNoRows {
		log.Println("CheckIDExists Query: ID does not exists ", err)
		return false, err
	}
	log.Println("CheckIDExists Query: ID exists!")
	return true, nil
}

// CheckEmailExists queries database if email exists
func (d *Db) CheckEmailExists(email string) (bool, error) {
	stmt, err := d.Prepare("SELECT email FROM users WHERE email=$1")
	if err != nil {
		log.Println("CheckEmailExists Preparation Error: ", err)
	}
	var exists string
	err = stmt.QueryRow(email).Scan(&exists)
	if err == sql.ErrNoRows {
		log.Println("CheckEmailExists Query email does not exists ", err)
		return false, err
	}
	log.Println("CheckEmailExists Query: email exists!")
	return true, nil
}

// StringToGraphqlID convert string array to graphqlID array
func StringToGraphqlID(s []string, gqlS *[]graphql.ID) {
	for i := 0; i < len(s); i++ {
		*gqlS = append(*gqlS, graphql.ID(s[i]))
	}
}

// GraphqlIDToUUID convert graphqlID array to UUID array
func GraphqlIDToUUID(gqlS []graphql.ID, u *[]uuid.UUID) {
	for i := 0; i < len(gqlS); i++ {
		x, _ := uuid.FromString(string(gqlS[i]))
		*u = append(*u, x)
	}
}

// UUIDToGraphqlID convert UUID array to graphqlID array
func UUIDToGraphqlID(uid []uuid.UUID, gid *[]graphql.ID) {
	for i := 0; i < len(uid); i++ {
		x := graphql.ID(uid[i].String())
		*gid = append(*gid, x)
	}
}

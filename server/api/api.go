package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/minio/minio-go"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
	"net/http"
	"strings"
	"time"
)

type Email struct {
	Email string `json:"email"`
}

type ProfileBuilder struct {
	ID       graphql.ID
	ImageURL string
}

//CheckEmailExists checks against the database to see if email exists in the system
func CheckEmailExists(w http.ResponseWriter, req *http.Request, db *postgres.Db) {
	var e Email
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if e.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: Please provide valid email!"))
	} else {
		exists, _ := db.CheckEmailExists(e.Email)
		if !exists {
			// exists, _ := json.Marshal(exists)
			// w.Write(exists)
			w.Write([]byte("false"))
		} else {
			// exists, _ := json.Marshal(exists)
			// w.Write(exists)
			w.Write([]byte("true"))
		}
	}
}

// UploadAnyS3 upload any file to S3
func (pb *ProfileBuilder) UploadAnyS3(w http.ResponseWriter, req *http.Request, minioClient *minio.Client, db *postgres.Db, tableType string, id graphql.ID) {
	// Create context for cancel deadline signal
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Gets multipart file and header
	file, header, err := req.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	contentType := header.Header["Content-Type"][0]
	contentDisposition := header.Header["Content-Disposition"][0]                                      // returns form-data
	contentDisposition = "inline;" + contentDisposition[strings.IndexByte(contentDisposition, ';')+1:] //replace form-data with inline
	fmt.Println("Content-Type ", contentType)
	fmt.Println("Content-Disposition ", contentDisposition)
	defer file.Close()
	bucketName := "doggy-date-app"
	path := "dogs/" //path to folder inside bucket
	success, err := minioClient.PutObjectWithContext(ctx, bucketName, path+header.Filename, file, header.Size, minio.PutObjectOptions{
		ContentType:        contentType,
		ContentDisposition: contentDisposition,
		// ContentType:        "application/octet-stream",
		// ContentDisposition: "inline",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	imgURL := fmt.Sprintf("https://d2m79q3ctf5ck3.cloudfront.net/%s", header.Filename)
	pb.UpdateProfilePic(db, tableType, id, imgURL)
	fmt.Println("Successfully uploaded bytes: ", success)
	w.Write([]byte(imgURL))
}

// UpdateProfilePic updates profile picture row in postgres
func (pb *ProfileBuilder) UpdateProfilePic(db *postgres.Db, tableType string, id graphql.ID, imgURL string) {
	if _, err := db.UpdateProfilePic(tableType, id, imgURL); err != nil {
		fmt.Println("UpdateProfilePic Err ", err)
	}
}

// UploadImage DEPRECATED uses old way to send image up to 10MB
// func UploadImage(w http.ResponseWriter, req *http.Request, db *postgres.Db, s3bucket *s3.S3) {
// var imgData string
// decoder := json.NewDecoder(req.Body)
// err := decoder.Decode(&imgData)
// if err != nil {
// 	w.WriteHeader(http.StatusBadRequest)
// 	w.Write([]byte(fmt.Sprintln(err)))
// }
// imgType := imgData[strings.IndexByte(imgData, '/')+1 : strings.IndexByte(imgData, ';')] //parse base64 url and get image type
// b64data := imgData[strings.IndexByte(imgData, ',')+1:]                                  // parse base64 url and get raw b64 data
// buff, err := base64.StdEncoding.DecodeString(b64data)                                   //decode b64 into array buffer
// if err != nil {
// 	fmt.Println("err", err)
// 	w.WriteHeader(http.StatusBadRequest)
// 	w.Write([]byte(fmt.Sprintln(err)))
// }
// reader := bytes.NewReader(buff) // convert array buffer into file
// i, _ := uuid.NewV4()
// imgName := i.String() + "." + strings.ToLower(imgType)
// // put file into s3 bucket
// resp, err := s3bucket.PutObject(&s3.PutObjectInput{
// 	Bucket: aws.String("doggy-date-app/dogs/"),
// 	Key:    aws.String(imgName),
// 	Body:   reader,
// })
// if err != nil {
// 	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
// 		// If the SDK can determine the request or retry delay was canceled
// 		// by a context the CanceledErrorCode error code will be returned.
// 		fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
// 	} else {
// 		fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
// 	}
// 	w.WriteHeader(http.StatusForbidden)
// }
// fmt.Printf("response %s", awsutil.StringValue(resp))
// }

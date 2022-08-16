package mongo_client

import (
	"log"
	"testing"
)

func TestFileUpload(t *testing.T) {
	collectionName := "file_coll"
	repository := SetMongoConfiguration("admin", "123456", "localhost", "27017", "model_core", collectionName)
	if id, ok := repository.FileUpload(collectionName, "C:/Users/AC/Desktop/1158.jpg_wh860.jpg"); ok {
		repository.Download(collectionName, id)
	} else {
		log.Println("error")
	}
}

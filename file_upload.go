package mongo_client

import (
	"bytes"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"path"
)

type mongoFileRepository struct {
	mongoRepository
	MongoBucket *gridfs.Bucket
}

type FileInfo struct {
	FileName    string
	Size        int
	ContentType string
	FileId      string
}

func (mc *mongoRepository) getMongoBucket(collectionName string) (mongoFileRepository, error) {
	bucketOptions := options.GridFSBucket().SetName(collectionName)
	bucket, err := gridfs.NewBucket(mc.MDatabase, bucketOptions)
	return mongoFileRepository{MongoBucket: bucket}, err
}

func (mc *mongoRepository) FileUpload(collectionName, filename string) (string, bool) {
	if repository, err := mc.getMongoBucket(collectionName); err == nil {
		var bts []byte
		bts, err = ioutil.ReadFile(filename)
		fileSimpleName := path.Base(filename)
		tp := path.Ext(fileSimpleName)
		id := primitive.NewObjectID()
		err = repository.MongoBucket.UploadFromStreamWithID(id, filename, bytes.NewBuffer(bts))
		fInfo := FileInfo{FileName: fileSimpleName, Size: len(bts), ContentType: tp, FileId: id.Hex()}
		result := mc.Save(fInfo, collectionName)
		return result, err == nil && result != ""
	}
	return "", false
}

func (mc *mongoRepository) Download(collectionName, fileId, downloadPath string) {
	fileInfo := FileInfo{}
	err := mc.FindById(fileId, collectionName).Decode(&fileInfo)
	repository, err := mc.getMongoBucket(collectionName)
	fileBuffer := bytes.NewBuffer(nil)
	var id primitive.ObjectID
	id, err = primitive.ObjectIDFromHex(fileInfo.FileId)
	_, err = repository.MongoBucket.DownloadToStream(id, fileBuffer)
	if err != nil {
		log.Printf("download failed ,err : %s ", err.Error())
	} else {
		err = ioutil.WriteFile(downloadPath+uuid.New().String()+fileInfo.ContentType, fileBuffer.Bytes(), 0666)
		if err != nil {
			log.Println("write file failed !")
		}
	}
}

func (mc *mongoRepository) previewFile() {

}

package mongo_client

import (
	"log"
	"testing"
)

func TestMongoDbOperation(t *testing.T) {
	collection := "go_collection"
	repository := SetMongoConfiguration("admin", "123456", "localhost", "27017", "model_core", "go_collection")

	entity := make(map[string]interface{})
	entity["test"] = "7666"
	entity["aaa"] = "ooo"
	id := repository.Save(entity, "go_collection")
	log.Println(id)

	res := repository.FindById(id, collection)
	br, err := res.DecodeBytes()
	log.Println(br)

	findRes, err := repository.FindAll(collection)
	if err != nil {
		log.Fatalln("find all error ")
	}

	for findRes.TryNext(repository.GetCtx()) {
		re, _ := findRes.Current.Elements()
		log.Println(re)
	}
	entity["test"] = 8888
	entity["aaa"] = "666666"
	if re, err := repository.Update(id, entity, collection); err != nil {
		log.Printf("update failed , %s", err.Error())
	} else {
		log.Println(re)
	}
	if resul, err := repository.FindById(id, collection).DecodeBytes(); err == nil {
		log.Println(resul)
	} else {
		log.Println("err")
	}

	r := repository.Delete(id, collection)
	br, _ = r.DecodeBytes()
	log.Printf("delete %s", br)

	az := repository.FindById(id, collection)
	log.Println(az.DecodeBytes())

	//repository.FindById(result.InsertedID,collection)
}

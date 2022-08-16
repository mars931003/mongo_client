package mongo_client

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoRepository interface {
	Save(entity interface{}, collection string) string
	Update(id string, entity map[string]interface{}, collection string) (*mongo.UpdateResult, error)
	Delete(id string, collection string) *mongo.SingleResult
	FindById(id, collection string) *mongo.SingleResult
	FindAll(collection string) (*mongo.Cursor, error)
	Find(condition map[string]interface{}, collection string) (*mongo.Cursor, error)
	FileUpload(collectionName, filename string) (string, bool)
	Download(collectionName, fileId, downloadPath string)
	previewFile()
	GetCtx() context.Context //FIXME delete
}

type mongoRepository struct {
	MongoRepository
	MDatabase   *mongo.Database
	Ctx         context.Context
	MCollection *mongo.Collection
}

func (mc *mongoRepository) GetCtx() context.Context {
	return mc.Ctx
}

// SetMongoConfiguration 配置并连接mongodb,
// uri mongodb://username:password@host:port/?authSource=admin
func SetMongoConfiguration(username, pwd, host, port, db, collection string) MongoRepository {
	uri := "mongodb://" + username + ":" + pwd + "@" + host + ":" + port + "/?authSource=admin"
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln("db collect failed !")
	}
	return &mongoRepository{MDatabase: client.Database(db), MCollection: client.Database(db).Collection(collection), Ctx: ctx}
}

func (mc *mongoRepository) Save(entity interface{}, collection string) string {
	if res, err := mc.MDatabase.Collection(collection).InsertOne(mc.Ctx, entity); err == nil {
		return convert2string(res)
	}

	log.Println("insert error ")
	return ""
}

func convert2string(result *mongo.InsertOneResult) string {
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func (mc *mongoRepository) Update(id string, entity map[string]interface{}, collection string) (*mongo.UpdateResult, error) {
	if objectId, err := primitive.ObjectIDFromHex(id); err != nil {
		log.Printf("convert id to objectId failed ! err : %s", err)
		return nil, err
	} else {
		return mc.MDatabase.Collection(collection).UpdateByID(mc.Ctx, objectId, setUpdate(entity))
	}
}

func (mc *mongoRepository) Delete(id string, collection string) *mongo.SingleResult {
	objectId := id2ObjectId(id)
	return mc.MDatabase.Collection(collection).FindOneAndDelete(mc.Ctx, objectId)
}

func (mc *mongoRepository) FindById(id, collection string) *mongo.SingleResult {
	objectId := id2ObjectId(id)
	return mc.MDatabase.Collection(collection).FindOne(mc.Ctx, objectId)
}

func (mc *mongoRepository) FindAll(collection string) (*mongo.Cursor, error) {
	return mc.MDatabase.Collection(collection).Find(mc.Ctx, bson.M{})
}

func (mc *mongoRepository) Find(condition map[string]interface{}, collection string) (*mongo.Cursor, error) {
	return mc.MDatabase.Collection(collection).Find(mc.Ctx, condition)
}

func id2ObjectId(id string) primitive.M {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("id error , can not transfer id string to objectId !")
		return nil
	}
	return bson.M{"_id": objectId}
}

func setUpdate(entity map[string]interface{}) *primitive.M {
	value := bson.M{}
	for k, v := range entity {
		value[k] = v
	}
	return &bson.M{"$set": value}
}

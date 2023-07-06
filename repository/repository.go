package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Save(ctx context.Context, entity BimEntity) error
	Update(ctx context.Context, entity BimEntity) error
	Delete(ctx context.Context, title string) error
	GetById(ctx context.Context, title string) (BimEntity, error)
	FullTextSearch(ctx context.Context, inputMsg string) []BimEntity
}

//
//type pgRepository struct {
//	db *sqlx.DB
//}
//
//func NewPgRepository(db *sqlx.DB) Repository {
//	return &pgRepository{db: db}
//}
//
//func (p *pgRepository) GetById(ctx context.Context, input string) (*models.ResponseModel, error) {
//	var a = &models.ResponseModel{}
//	if err := p.db.GetContext(
//		ctx, a,
//		findById,
//		input,
//	); err != nil {
//		return nil, errors.Wrap(err, "pgRepository.GetById.GetContext")
//	}
//	return a, nil
//}

type mongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) Repository {
	return &mongoRepository{collection: collection}
}

type BimEntity struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Score       string             `bson:"score"`
}

func (mRepository mongoRepository) Save(ctx context.Context, entity BimEntity) error {
	////Get MongoDB connection using connectionhelper.
	//client, err := connectionhelper.GetMongoClient()
	//if err != nil {
	//	return err
	//}
	////Create a handle to the respective collection in the database.
	//collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	////Perform InsertOne operation & validate against the error.
	_, err := mRepository.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func (mRepository mongoRepository) Update(ctx context.Context, entity BimEntity) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "_id", Value: entity.ID}}

	//Define updater for to specifiy change to be updated.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "description", Value: entity.Description},
	}}}

	////Get MongoDB connection using connectionhelper.
	//client, err := connectionhelper.GetMongoClient()
	//if err != nil {
	//	return err
	//}
	//collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)

	//Perform UpdateOne operation & validate against the error.
	_, err := mRepository.collection.UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func (mRepository mongoRepository) Delete(ctx context.Context, title string) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "title", Value: title}}
	//Get MongoDB connection using connectionhelper.
	//client, err := connectionhelper.GetMongoClient()
	//if err != nil {
	//	return err
	//}
	////Create a handle to the respective collection in the database.
	//collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	//Perform DeleteOne operation & validate against the error.
	_, err := mRepository.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

func (mRepository mongoRepository) GetById(ctx context.Context, id string) (BimEntity, error) {
	result := BimEntity{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	//Get MongoDB connection using connectionhelper.
	//client, err := connectionhelper.GetMongoClient()
	//if err != nil {
	//	return result, err
	//}
	////Create a handle to the respective collection in the database.
	//collection := client.Database(connectionhelper.DB).Collection(connectionhelper.ISSUES)
	//Perform FindOne operation & validate against the error.
	err := mRepository.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	//Return result without any error.
	return result, nil
}

func (mRepository mongoRepository) FullTextSearch(ctx context.Context, inputMsg string) []BimEntity {
	results := []BimEntity{}

	filter := bson.D{{"$text", bson.D{{"$search", inputMsg}}}}
	sort := bson.D{{"score", bson.D{{"$meta", "textScore"}}}}
	projection := bson.D{{"description", 1}, {"score", bson.D{{"$meta", "textScore"}}}, {"_id", 0}}
	opts := options.Find().SetSort(sort).SetProjection(projection)

	cursor, err := mRepository.collection.Find(ctx, filter, opts)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	return results
}

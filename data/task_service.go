package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	// "time"

	"mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-playground/validator"
)

type taskService interface {
	GetTasks() ([]models.Task, error)
	CreateTasks(task *models.Task) (models.Task, error)
	GetTasksById(id string) (models.Task, error)
	UpdateTasksById(id string, task models.Task) (models.Task, error)
	DeleteTasksById(id string) error
}

type TaskService struct {
	validator  *validator.Validate
	client     *mongo.Client
	DataBase   *mongo.Database
	collection *mongo.Collection
}

func NewTaskService() (*TaskService, error) {
	MONGO_CONN := os.Getenv("MONGO_CONN")
	clientOptions := options.Client().ApplyURI(MONGO_CONN)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected ")
	DataBase := client.Database("test")
	collection := DataBase.Collection("tasks")
	ts := &TaskService{
		client:     client,
		validator:  validator.New(),
		DataBase:   DataBase,
		collection: collection,
	}

	return ts, nil
}



func (ts *TaskService) CreateTasks(task *models.Task) (models.Task, error, int) {

	_, err := ts.collection.InsertOne(context.TODO(), task)
	if err != nil {
		return *task, err, 500
	}
	return *task , nil ,200
}
// func (ts *TaskService) CreateTasks(task *models.Task) (models.Task, error, int) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//     fmt.Println(task)
// 	session, err := ts.client.StartSession()
// 	if err != nil {
// 		fmt.Println(err)
// 		return models.Task{}, err, 500
// 	}
// 	defer session.EndSession(ctx)

// 	resultTask := models.Task{}
// 	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
// 		err := session.StartTransaction()
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}

// 		insertResult, err := ts.collection.InsertOne(sc, task)
// 		if err != nil {
// 			fmt.Println(err)
// 			session.AbortTransaction(sc)
// 			return err
// 		}

// 		var fetched models.Task
// 		err = ts.collection.FindOne(sc, bson.D{{"_id", insertResult.InsertedID.(primitive.ObjectID)}}).Decode(&fetched)

// 		if err != nil {
// 			fmt.Println(err)
// 			session.AbortTransaction(sc)
// 			return errors.New("Task Not Created")
// 		}
// 		var date_not_match = !fetched.DueDate.Equal(task.DueDate.In(fetched.DueDate.Location()))
// 		if fetched.Title != task.Title || fetched.Description != task.Description || fetched.Status != task.Status || date_not_match {
// 			session.AbortTransaction(sc)
// 			return errors.New("Task Not Created")
// 		}

// 		err = session.CommitTransaction(sc)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}

// 		resultTask = fetched
// 		return nil
// 	})

// 	if err != nil {
// 		fmt.Println(err)
// 		return models.Task{}, err, 500
// 	}

// 	fmt.Println("Inserted a single document: ", resultTask.ID)
// 	return resultTask, nil, 201
// }

func (ts *TaskService) GetTasks() ([]*models.Task, error, int) {
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return nil, err1, 500
	}

	findOptions := options.Find()
	filter := bson.D{{}}

	var results []*models.Task

	cur, err := ts.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
		return []*models.Task{}, err, 0
	}

	for cur.Next(context.TODO()) {
		var elem models.Task
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err.Error())
			return []*models.Task{}, err, 500
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return []*models.Task{}, err, 500
	}

	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	return results, nil, 200
}

func (ts *TaskService) GetTasksById(id primitive.ObjectID) (models.Task, error, int) {
	_, err1 := ts.collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{"_id", 1}},
	})
	if err1 != nil {
		return models.Task{}, err1, 500
	}

	filter := bson.D{{"_id", id}}
	var result models.Task
	err := ts.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.Task{}, errors.New("Task not found"), http.StatusNotFound
	}
	return result, nil, 200
}

func (ts *TaskService) UpdateTasksById(id primitive.ObjectID, task models.Task) (models.Task, error, int) {
	session, err := ts.client.StartSession()
	if err != nil {
		fmt.Println(err)
		return models.Task{}, err, 500
	}
	defer session.EndSession(context.Background())

	var NewTask models.Task
	statusCode := 200

	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		err = session.StartTransaction()
		if err != nil {
			fmt.Println(err)
			return err
		}

		NewTask, err, statusCode = ts.GetTasksById(id)
		if err != nil {
			_ = session.AbortTransaction(sc)
			return errors.New("task does not exist")
		}

		if task.Title != "" {
			NewTask.Title = task.Title
		}
		if task.Description != "" {
			NewTask.Description = task.Description
		}
		if task.Status != "" {
			NewTask.Status = task.Status
		}
		if !task.DueDate.IsZero() {
			NewTask.DueDate = task.DueDate
		}

		filter := bson.D{{"_id", id}}
		update := bson.D{
			{"$set", bson.D{
				{"title", NewTask.Title},
				{"description", NewTask.Description},
				{"status", NewTask.Status},
				{"due_date", NewTask.DueDate},
			}},
		}

		updateResult, err := ts.collection.UpdateOne(sc, filter, update)
		if err != nil {
			_ = session.AbortTransaction(sc)
			fmt.Println(err)
			return err
		}
		if updateResult.ModifiedCount == 0 {
			return errors.New("task does not exist")
		}

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		err = session.CommitTransaction(sc)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		return models.Task{}, err, 500
	}

	return NewTask, nil, statusCode
}

func (ts *TaskService) DeleteTasksById(id primitive.ObjectID) (error, int) {
	filter := bson.D{{"_id", id}}

	deleteResult, err := ts.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return err, 500
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("Task does not exist"), http.StatusNotFound
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil, 200
}

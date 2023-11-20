package repository

import (
	"context"

	"delta-core/domain"
	"delta-core/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *taskRepository) Create(c context.Context, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.InsertOne(c, task)

	return err
}

func (tr *taskRepository) FetchById(c context.Context, taskId string) (domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	var task domain.Task

	idHex, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return domain.Task{}, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&task)

	return task, err

}

func (tr *taskRepository) FetchByUserID(c context.Context, userID string) ([]domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	var tasks []domain.Task

	idHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return tasks, err
	}

	cursor, err := collection.Find(c, bson.M{"userID": idHex})
	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &tasks)
	if tasks == nil {
		return []domain.Task{}, err
	}

	return tasks, err
}

func (tr *taskRepository) FetchAll(c context.Context) ([]domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	var tasks []domain.Task

	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(c, &tasks)
	if tasks == nil {
		return []domain.Task{}, err
	}

	return tasks, err
}
func (tr *taskRepository) Delete(c context.Context, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)

	_, err := collection.DeleteOne(c, task)
	return err
}

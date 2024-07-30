package route

import (
	"context"
	"fmt"
	"time"

	"delta-core/api/controller"
	"delta-core/bootstrap"
	"delta-core/domain"
	"delta-core/mongo"
	"delta-core/repository"
	"delta-core/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(db, domain.CollectionTask)
	pr := repository.NewUserRepository(db, domain.CollectionUser)
	puc := usecase.NewProfileUsecase(pr, timeout)

	tc := &controller.TaskController{
		TaskUsecase:      usecase.NewTaskUsecase(tr, timeout),
		SignalSubUsecase: usecase.NewSignalSubUsecase(tr, timeout),
		ProfileUsecase:   usecase.NewProfileUsecase(pr, timeout),
	}
	var taskCollection = db.Collection(domain.CollectionTask)
	var tasks []domain.Task

	cursor, err := taskCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Printf("Cannot fetch all tasks, error %v", err)
	}

	cursor.All(context.TODO(), &tasks)
	tc.SignalSubUsecase.InitialiseSingalSubs(context.TODO(), env, puc, tasks)

	group.GET("/task", tc.Fetch)
	group.POST("/task", tc.Create)
	group.DELETE("/task", tc.Cancel)
}

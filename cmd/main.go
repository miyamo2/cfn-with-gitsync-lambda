package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/miyamo2/dynmgrm"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	usersTableName string
)

type Users struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pathParam := request.PathParameters
	id, ok := pathParam["id"]
	if !ok {
		slog.Error("invalid path parameter")
		return events.APIGatewayProxyResponse{
			Body:       string("Invalid Request"),
			StatusCode: 400,
		}, nil
	}

	u := &Users{}
	if err := db.Table("users").Where("id=?", id).Scan(u).Error; err != nil {
		slog.Error("failed to get user", slog.String("error", err.Error()))
		return events.APIGatewayProxyResponse{
			Body:       string("Internal Server Error"),
			StatusCode: 500,
		}, nil
	}

	if u.ID == "" {
		slog.Warn("user not found", slog.String("id", id))
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("%s are not found", id),
			StatusCode: 404,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%s: %s", u.ID, u.Name),
		StatusCode: 200,
	}, nil
}

func init() {
	usersTableName = os.Getenv("USERS_TABLE_NAME")
	if usersTableName == "" {
		log.Fatal("USERS_TABLE_NAME are not set")
	}
	g, err := gorm.Open(dynmgrm.New())
	if err != nil {
		log.Fatal(err.Error())
	}
	db = g
}

func main() {
	lambda.Start(handler)
}

package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	_diaryLambdaDeliver "github.com/hareku/pastory-api/diary/delivery/lambda"
	_diaryRepo "github.com/hareku/pastory-api/diary/repository"
	_diaryUcase "github.com/hareku/pastory-api/diary/usecase"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	awsConf := &aws.Config{
		Region: aws.String("ap-northeast-1"),
	}

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint != "" {
		awsConf.Endpoint = aws.String(awsEndpoint)
	}

	db := dynamo.New(session.New(), awsConf)

	dr := _diaryRepo.NewDynamoDiaryRepository(db)
	du := _diaryUcase.NewDiaryUsecase(dr)
	dd := _diaryLambdaDeliver.NewDiaryLambdaHandler(du)

	res, err := dd.HandleRequest(&request)
	return *res, err
}

func main() {
	lambda.Start(handler)
}

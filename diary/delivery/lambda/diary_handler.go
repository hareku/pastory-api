package lambda

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/hareku/pastory-api/auth"
	"github.com/hareku/pastory-api/diary"
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

// lambdaDiaryHandler  represent the lambdaHandler for diary
type lambdaDiaryHandler struct {
	DUsecase diary.Usecase
}

type LambdaContext struct {
	Req *events.APIGatewayProxyRequest
}

func NewDiaryLambdaHandler(us diary.Usecase) *lambdaDiaryHandler {
	return &lambdaDiaryHandler{
		DUsecase: us,
	}
}

func (d *lambdaDiaryHandler) HandleRequest(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	res, err := d.makeResponse(request)

	if res.Headers == nil {
		res.Headers = map[string]string{}
	}

	res.Headers["Access-Control-Allow-Origin"] = "*"
	res.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS"
	res.Headers["Access-Control-Allow-Headers"] = "Authorization, Content-Type"

	return res, err
}

func (d *lambdaDiaryHandler) makeResponse(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	lambdaCtx := &LambdaContext{
		Req: request,
	}

	switch {
	case lambdaCtx.Req.HTTPMethod == "OPTIONS":
		return d.optionsResponse()
	case lambdaCtx.Req.HTTPMethod == "GET" && lambdaCtx.Req.Path == "/diaries":
		return d.FetchDiaries(lambdaCtx)
	case lambdaCtx.Req.HTTPMethod == "POST" && lambdaCtx.Req.Path == "/diaries":
		return d.CreateDiary(lambdaCtx)
	case lambdaCtx.Req.HTTPMethod == "PUT" && strings.HasPrefix(lambdaCtx.Req.Path, "/diaries/"):
		return d.UpdateDiary(lambdaCtx)
	case lambdaCtx.Req.HTTPMethod == "DELETE" && strings.HasPrefix(lambdaCtx.Req.Path, "/diaries/"):
		return d.DeleteDiary(lambdaCtx)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       "Not found",
	}, nil
}

func (d *lambdaDiaryHandler) forbiddenResponse(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 403,
		Body:       "Forbidden",
	}, err
}

func (d *lambdaDiaryHandler) internalErrorResponse(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Internal error",
	}, err
}

func (d *lambdaDiaryHandler) optionsResponse() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "",
	}, nil
}

func (d *lambdaDiaryHandler) FetchDiaries(c *LambdaContext) (*events.APIGatewayProxyResponse, error) {
	token, err := auth.AuthFirebaseForLambda(c.Req)
	if err != nil {
		return d.forbiddenResponse(err)
	}

	diaries, err := d.DUsecase.FetchMany(token.UID)
	if err != nil {
		return d.internalErrorResponse(err)
	}

	jsonBytes, err := json.Marshal(diaries)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to parse diaries to json"))
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func (d *lambdaDiaryHandler) CreateDiary(c *LambdaContext) (*events.APIGatewayProxyResponse, error) {
	token, err := auth.AuthFirebaseForLambda(c.Req)
	if err != nil {
		return d.forbiddenResponse(err)
	}

	createInput := new(diary.CreateInput)
	err = json.Unmarshal([]byte(c.Req.Body), createInput)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to parse json of request body"))
	}
	createInput.UserID = token.UID

	validate := validator.New()
	validate.RegisterValidation("date", diary.DateValidation)
	err = validate.Struct(createInput)

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 422,
			Body:       "InvalidParameter",
		}, nil
	}

	diary, err := d.DUsecase.Create(createInput)

	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to create diary"))
	}

	jsonBytes, err := json.Marshal(diary)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to parse diary to json"))
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func (d *lambdaDiaryHandler) UpdateDiary(c *LambdaContext) (*events.APIGatewayProxyResponse, error) {
	token, err := auth.AuthFirebaseForLambda(c.Req)
	if err != nil {
		return d.forbiddenResponse(err)
	}

	updateInput := new(diary.UpdateInput)
	err = json.Unmarshal([]byte(c.Req.Body), updateInput)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to parse json of request body"))
	}

	updateInput.DiaryID = strings.TrimPrefix(c.Req.Path, "/diaries/")
	updateInput.UserID = token.UID

	validate := validator.New()
	err = validate.Struct(updateInput)

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 422,
			Body:       "InvalidParameter",
		}, nil
	}

	diary, err := d.DUsecase.Update(updateInput)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to update diary"))
	}

	jsonBytes, err := json.Marshal(diary)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to parse diary to json"))
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func (d *lambdaDiaryHandler) DeleteDiary(c *LambdaContext) (*events.APIGatewayProxyResponse, error) {
	token, err := auth.AuthFirebaseForLambda(c.Req)
	if err != nil {
		return d.forbiddenResponse(err)
	}

	deleteInput := &diary.DeleteInput{
		ID:     strings.TrimPrefix(c.Req.Path, "/diaries/"),
		UserID: token.UID,
	}

	err = d.DUsecase.Delete(deleteInput)
	if err != nil {
		return d.internalErrorResponse(errors.Wrap(err, "Failed to delete diary"))
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 204,
		Body:       "",
	}, nil
}

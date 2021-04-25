package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	req := new(openapi.DeleteTagsRequest)

	headers := map[string]string{"Access-Control-Allow-Origin": "*"}
	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	}

	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	usedTags, err := taggedImageUsecase.ValidateDeleteTags(*userName, *req.Tags)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers: headers,
		}, err
	} else if len(usedTags) != 0 {
		message := "Some Tags are now used. Tags:"
		for _, e := range usedTags {
			message += e + " "
		}
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       message,
			StatusCode: 400,
		}, err
	}

	tagUsecase := usecase.NewTagUsecase()
	tagUsecase.Delete(*userName, *req.Tags)

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

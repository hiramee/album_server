package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	TagUsecase := usecase.NewTagUsecase(ctx)
	userName, _ := util.GetUsernameFromHeader(request)

	results, err := TagUsecase.ListAll(*userName)

	if err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}
	response := new(openapi.GetTagsResponse)
	var tags []string
	for _, e := range results {
		tags = append(tags, e.TagName)
	}
	response.Tags = &tags

	return util.CreateOKResponse(response)
}

func main() {
	lambda.Start(handler)
}

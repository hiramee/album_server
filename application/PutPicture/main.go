package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	id := request.PathParameters["id"]
	req := new(openapi.PutPictureRequest)

	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, nil)
	}

	tagUsecase := usecase.NewTagUsecase(ctx)
	if err := tagUsecase.SaveTagIfAbsent(*userName, req.Tags); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, nil)
	}
	taggedImageUsecase := usecase.NewTaggedImageUsecase(ctx)
	taggedImageUsecase.UpdateTaggedImage(*userName, id, req.Tags)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}

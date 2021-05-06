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
	id := request.PathParameters["id"]
	req := new(openapi.PutPictureRequest)

	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, nil)
	}

	tagUsecase := usecase.NewTagUsecase()
	if err := tagUsecase.SaveTagIfAbsent(*userName, req.Tags); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, nil)
	}
	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	taggedImageUsecase.UpdateTaggedImage(*userName, id, req.Tags)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}

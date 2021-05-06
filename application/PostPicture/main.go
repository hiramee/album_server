package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	req := new(openapi.PostPicturesRequest)

	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}

	tagUsecase := usecase.NewTagUsecase()
	if err := tagUsecase.SaveTagIfAbsent(*userName, req.Tags); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}
	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	decoded, err := base64.StdEncoding.DecodeString(req.Picture)
	if err != nil {
		return util.CreateErrorResponse(nil, util.APPLICATION_ERROR, err)
	}
	taggedImageUsecase.SaveTaggedImage(*userName, req.Tags, req.Ext, decoded)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}

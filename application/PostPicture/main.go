package main

import (
	"album-server/application/usecase"
	"album-server/openapi"
	"album-server/util"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const image_upper_size = 2 * 1000 * 1000

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userName, _ := util.GetUsernameFromHeader(request)
	req := new(openapi.PostPicturesRequest)

	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}

	decoded, err := base64.StdEncoding.DecodeString(req.Picture)
	if err != nil {
		return util.CreateErrorResponse(nil, util.APPLICATION_ERROR, err)
	}
	if len(decoded) > image_upper_size {
		return util.CreateErrorResponse("file size is over 2MB", util.VALIDATION_ERROR, fmt.Errorf("file size is over 2MB. size = %d B", len(decoded)))
	}
	tagUsecase := usecase.NewTagUsecase()

	if err := tagUsecase.SaveTagIfAbsent(*userName, req.Tags); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}
	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	taggedImageUsecase.SaveTaggedImage(*userName, req.Tags, req.Ext, decoded)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}

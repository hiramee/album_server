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

	if err := json.Unmarshal([]byte(request.Body), req); err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	}

	taggedImageUsecase := usecase.NewTaggedImageUsecase()
	usedTags, err := taggedImageUsecase.ValidateDeleteTags(*userName, *req.Tags)
	if err != nil {
		return util.CreateErrorResponse(nil, util.VALIDATION_ERROR, err)
	} else if len(usedTags) != 0 {
		message := "some tags are now used. tags: "
		for _, e := range usedTags {
			message += e + " "
		}
		return util.CreateErrorResponse(message, util.VALIDATION_ERROR, err)
	}

	tagUsecase := usecase.NewTagUsecase()
	tagUsecase.DeleteTag(*userName, *req.Tags)

	return util.CreateOKResponse(nil)
}

func main() {
	lambda.Start(handler)
}

// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package openapi

const (
	IDTokenScopes = "IDToken.Scopes"
)

// DeleteTagsRequest defines model for DeleteTagsRequest.
type DeleteTagsRequest struct {

	// List of tags
	Tags *[]string `json:"tags,omitempty"`
}

// GetPictureResponse defines model for GetPictureResponse.
type GetPictureResponse struct {

	// Base64encoded picture
	Picture *string `json:"picture,omitempty"`
}

// GetPicturesResponse defines model for GetPicturesResponse.
type GetPicturesResponse struct {

	// List of pictures
	Pictures *[]PicturesResponseItem `json:"pictures,omitempty"`
}

// GetTagsResponse defines model for GetTagsResponse.
type GetTagsResponse struct {

	// List of tags
	Tags *[]string `json:"tags,omitempty"`
}

// PicturesResponseItem defines model for PicturesResponseItem.
type PicturesResponseItem struct {

	// fileName
	FileName *string `json:"fileName,omitempty"`

	// ID
	Id *string `json:"id,omitempty"`

	// List of tags
	Tags *[]string `json:"tags,omitempty"`
}

// PostPicturesRequest defines model for PostPicturesRequest.
type PostPicturesRequest struct {

	// extension of file
	Ext string `json:"ext"`

	// Base64Encoded picture
	Picture string `json:"picture"`

	// List of tags
	Tags []string `json:"tags"`
}

// PutPictureRequest defines model for PutPictureRequest.
type PutPictureRequest struct {

	// List of tags
	Tags []string `json:"tags"`
}

// GetPicturesParams defines parameters for GetPictures.
type GetPicturesParams struct {

	// tag
	Tag []string `json:"tag"`
}

// PostPicturesJSONBody defines parameters for PostPictures.
type PostPicturesJSONBody PostPicturesRequest

// GetPicturesIdParams defines parameters for GetPicturesId.
type GetPicturesIdParams struct {

	// true if getting thumbNail
	Thumbnail *bool `json:"thumbnail,omitempty"`
}

// PutPicturesIdJSONBody defines parameters for PutPicturesId.
type PutPicturesIdJSONBody PutPictureRequest

// PostTagsDeleteJSONBody defines parameters for PostTagsDelete.
type PostTagsDeleteJSONBody DeleteTagsRequest

// PostPicturesJSONRequestBody defines body for PostPictures for application/json ContentType.
type PostPicturesJSONRequestBody PostPicturesJSONBody

// PutPicturesIdJSONRequestBody defines body for PutPicturesId for application/json ContentType.
type PutPicturesIdJSONRequestBody PutPicturesIdJSONBody

// PostTagsDeleteJSONRequestBody defines body for PostTagsDelete for application/json ContentType.
type PostTagsDeleteJSONRequestBody PostTagsDeleteJSONBody


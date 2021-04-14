// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package openapi

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

	// Base64encoded picture
	Picture *string `json:"picture,omitempty"`

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

// GetPicturesParams defines parameters for GetPictures.
type GetPicturesParams struct {

	// tag
	Tag []string `json:"tag"`
}

// PostPicturesJSONBody defines parameters for PostPictures.
type PostPicturesJSONBody PostPicturesRequest

// PostPicturesJSONRequestBody defines body for PostPictures for application/json ContentType.
type PostPicturesJSONRequestBody PostPicturesJSONBody

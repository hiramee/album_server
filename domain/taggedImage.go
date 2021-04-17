package domain

type TaggedImage struct {
	ID          string   `dynamo:"ID,hash"`
	UserTagName string   `dynamo:"UserTagName,range" index:"GSI-UserTagName"`
	ObjectKey   string   `dynamo:",omitempty"`
	Tags        []string `dynamo:",omitempty"`
}

func NewTaggedImage(userName string, tagName string, id string, objectKey string) *TaggedImage {
	taggedImage := new(TaggedImage)
	taggedImage.UserTagName = userName + tagName
	taggedImage.ID = id
	taggedImage.ObjectKey = objectKey
	var tagSlice []string
	tagSlice = append(tagSlice, tagName)
	taggedImage.Tags = tagSlice
	return taggedImage
}

func NewTaggedImageSlice(userName string, tagNames []string, id string, objectKey string) []TaggedImage {
	taggedImageSlice := make([]TaggedImage, len(tagNames))
	for i, e := range tagNames {
		taggedImage := new(TaggedImage)
		taggedImage.UserTagName = userName + e
		taggedImage.ID = id
		taggedImage.ObjectKey = objectKey
		taggedImage.Tags = tagNames
		taggedImageSlice[i] = *taggedImage
	}
	return taggedImageSlice
}

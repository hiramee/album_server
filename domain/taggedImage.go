package domain

type TaggedImage struct {
	UserTagName string // Hash key, a.k.a. partition key
	ID          string // Range key, a.k.a. sort key
	ObjectKey   string `dynamo:",omitempty"`
}

func NewTaggedImage(userName string, tagName string, id string, objectKey string) *TaggedImage {
	taggedImage := new(TaggedImage)
	taggedImage.UserTagName = userName + tagName
	taggedImage.ID = id
	taggedImage.ObjectKey = objectKey
	return taggedImage
}

func NewTaggedImageSlice(userName string, tagNames []string, id string, objectKey string) []TaggedImage {
	taggedImageSlice := make([]TaggedImage, len(tagNames))
	for i, e := range tagNames {
		taggedImage := new(TaggedImage)
		taggedImage.UserTagName = userName + e
		taggedImage.ID = id
		taggedImage.ObjectKey = objectKey
		taggedImageSlice[i] = *taggedImage
	}
	return taggedImageSlice
}

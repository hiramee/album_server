package domain

type TaggedImage struct {
	UserTagName string // Hash key, a.k.a. partition key
	ID          string // Range key, a.k.a. sort key
	ObjectKey   string `dynamo:",omitempty"`
}

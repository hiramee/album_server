package domain

type Tag struct {
	UserName string // Hash key, a.k.a. partition key
	TagName  string // Range key, a.k.a. sort key
}

package db

type Bucket string

const (
	AccessToken  Bucket = "access_token"
	RequestToken Bucket = "request_token"
)

type TokenDb interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, token string, bucket Bucket) error
}

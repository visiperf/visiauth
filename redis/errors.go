package redis

const (
	ErrRedisNilMessage = "redis: nil"
)

func IsErrRedisNilMessage(err error) bool {
	return err.Error() == ErrRedisNilMessage
}

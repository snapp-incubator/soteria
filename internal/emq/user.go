package emq

// User contains the user information for storing in emq based on its redis-auth plugin.
// https://github.com/emqx/emqx-auth-redis
type User struct {
	Username    string
	Password    string
	IsSuperuser bool
}

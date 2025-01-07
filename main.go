package main

import "github.com/go-redis/redis/v8"

func main() {
	opt, err := redis.ParseURL("rediss://red-ctuhoid2ng1s739ecvcg:Bi7qgaBY0rnPxz8IWE0IfytMVey5b9EF@oregon-redis.render.com:6379")
	if err != nil {
		return
	}

	redis.NewClient(opt)

}

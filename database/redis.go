// database/redis.go

package database

import(
	"github.com/go-redis/redis/v8"
	// "time"
	// "net"
	// "bytes"
	// "encoding/json"
	// "strings"
	// "crypto/tls"
	// "fmt"
	// "context"

)

var rdb *redis.Client

func RDb() (rdb *redis.Client){
	

	rdb = redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
		Password: "", // no password set
		DB:		  0,  // use default DB
	})

	return rdb
}
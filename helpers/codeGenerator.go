package helpers

import (
	"crypto/rand"
	//"go-backend-ailanglearn/configs"
	"io"

	//"go.mongodb.org/mongo-driver/mongo"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
//var tokenCollection *mongo.Collection = configs.GetColletion(configs.DB, "tokens")


func HandleCodeGenerator(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}


package utilities

import (
	"context"
	"fmt"
	"net/http"
	"math/rand"
	"log"
)

type key int
const requestIdKey = key(42)

func Println(ctx context.Context, msg string) {
	id, ok := ctx.Value(requestIdKey).(int64)
	if !ok {
		log.Println("Couldn't retrieve request Id")
		return
	}
	fmt.Println("[%d] %s\n", id, msg)
}

func Decorate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := rand.Int63()
		ctx = context.WithValue(ctx, requestIdKey, id)
		f(w, r.WithContext(ctx))
	}
}
package log

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"
	"os"
)

type key int

const requestIdKey = key(42)

/*Incomplete function
 */
func Log(ctx context.Context, msg ...string) {
	f, err := os.OpenFile("local.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	id, ok := ctx.Value(requestIdKey).(int64)
	if !ok {
		ctx = context.WithValue(ctx, requestIdKey, string(time.Nanosecond))
		id, _ = ctx.Value(requestIdKey).(int64)
	}
	log.Printf("[%d] %s\n", id, msg)
}

/*
Specification: DecorateWithLog function is a wrapper over http.HandlerFunc
*/
func DecorateWithLog(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := rand.Int63()
		ctx = context.WithValue(ctx, requestIdKey, id)
		Log(ctx, "Peer's IP:%v", r.RemoteAddr)
		f(w, r.WithContext(ctx))
	}
}

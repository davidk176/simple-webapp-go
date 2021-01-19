package main

import (
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
	"log"
	"net/http"
)

func writeDataInCache(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	m := map[string]string{
		"email": "test@example.com",
	}
	log.Print("set to memchache", m)
	c, _ := json.Marshal(m)
	i := &memcache.Item{
		Key:   "a",
		Value: c,
	}

	memcache.Set(ctx, i)

	ig, err := memcache.Get(ctx, "a")
	log.Print(ig)
	log.Print(err)

}

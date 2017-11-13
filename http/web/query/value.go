package query

import (
	"http/web"
	"net/http"
	"net/url"
)

type Query struct {
	url.Values
}

func GetQuery(r *http.Request) Query {
	return Query{Values: r.URL.Query()}
}

func (q Query) MustGetString(key string) string {
	value := q.Get(key)
	if value == "" {
		panic(web.BadRequest("missing " + key))
	}
	return value
}

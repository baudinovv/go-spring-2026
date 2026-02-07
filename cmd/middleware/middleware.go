package middleware

import "net/http"

// func Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Header.get("X-API-KEY" != "secret123") {
// 			log.Printf("%s, %s, %s", r.Method, r.URL.Path, "Unauthorized")
// 			http.Error(wm)
// 		}
// 	})
// }
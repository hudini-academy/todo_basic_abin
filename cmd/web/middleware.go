package main

import (
	//"log"
	"net/http"
	"bytes"
	"fmt"
)
type responseWriter struct{
	http.ResponseWriter
	body *bytes.Buffer
}
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")


		rw := &responseWriter{
			ResponseWriter: w,
			body : 	bytes.NewBuffer(nil),
		}
		next.ServeHTTP(rw, r)
		//log.Printf(rw.body.String())
	})
}
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s -%s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		//log.Printf("%s -%s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}
func (rw *responseWriter) Write(b []byte) (int, error) {
    rw.body.Write(b)
    return rw.ResponseWriter.Write(b)
}
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer func(){
			if err := recover();
			err!= nil{
				w.Header().Set("Connection","close")
				app.serverError(w,fmt.Errorf("%s",err))
			}
		}()
			next.ServeHTTP(w,r)
	})
}
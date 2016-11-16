package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

func main() {

	http.HandleFunc("/a", generalMiddleware(handlerA))
	http.HandleFunc("/b", generalMiddleware(handlerB))
	http.HandleFunc("/setCount", generalMiddleware(handlerSetCount))
	http.HandleFunc("/", generalMiddleware(handlerRoot))
	http.ListenAndServe(":8080", nil)
}

func generalMiddleware(controllerFunc http.HandlerFunc) http.HandlerFunc {
	return handleErrors(requestNumber(timer(controllerFunc)))
}

var count int

func handleErrors(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			e := recover()
			if e != nil {
				res.WriteHeader(500)
				res.Write([]byte("there sure was an error"))
			}
		}()
		next(res, req)
	}
}

func requestNumber(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		count++
		req.Header.Set("request-number", strconv.Itoa(count))
		next(res, req)
	}
}

func timer(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		requestNumber := req.Header.Get("request-number")
		start := time.Now()
		defer func() {
			end := time.Now()
			fmt.Printf("requst #%v took %v\n", requestNumber, end.Sub(start))
		}()
		next(res, req)
	}
}

func handlerRoot(res http.ResponseWriter, req *http.Request) {
	responseBytes := []byte(fmt.Sprintf(`You're at the root.  My request number was '%v'`, req.Header.Get("request-number")))
	res.Write(responseBytes)
}

func handlerA(res http.ResponseWriter, req *http.Request) {

	responseBytes := []byte(fmt.Sprintf(`Welcome to A.  My request number was '%v'`, req.Header.Get("request-number")))
	res.Write(responseBytes)
}

func handlerB(res http.ResponseWriter, req *http.Request) {
	if !isEvenRequest(req) {
		panic("I can't even.")
	}
	responseBytes := []byte(fmt.Sprintf(`Welcome to B.  My request number was '%v'`, req.Header.Get("request-number")))
	res.Write(responseBytes)
}

func handlerSetCount(res http.ResponseWriter, req *http.Request) {
	countParam, err := strconv.Atoi(req.URL.Query().Get("count"))
	if err != nil {
		panic(err)
	}
	count = countParam
	responseBytes := []byte("count has been reset")
	res.Write(responseBytes)
}

func isEvenRequest(req *http.Request) bool {
	requestNumber, err := strconv.ParseFloat(req.Header.Get("request-number"), 64)
	if err != nil {
		panic(err)
	}

	if math.Mod(requestNumber, 2) == 0 {
		return true
	}
	return false

}

//middleware explained
//show no middleware with request number and logger repeated everywhere
//now do it with middleware
//let's handle some errors in controllers
//now let's handle these errors with a new middleware piece
//to get logging middleware to work again, we need to use defer there too, but that reads better anyway
//typed errors

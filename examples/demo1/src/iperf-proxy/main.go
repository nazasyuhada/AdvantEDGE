/*
 * MEEP Demo App API
 *
 * This is the MEEP Demo App API
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
        sw "github.com/InterDigitalInc/AdvantEDGE/iperfproxy/go"
        "net/http"
	"log"
        "github.com/gorilla/handlers"
)

func init() {
        // Initialize App
        sw.Init()
}

func main() {
        log.Printf("Demo Iperf transit App API Server started")

        router := sw.NewRouter()

        methods := handlers.AllowedMethods([]string{"OPTIONS", "DELETE", "GET", "HEAD", "POST", "PUT"})
        header := handlers.AllowedHeaders([]string{"content-type"})

        http.ListenAndServe(":30220", handlers.CORS(methods, header)(router))
}

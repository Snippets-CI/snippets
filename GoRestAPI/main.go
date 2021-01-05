// Snippets
//
// the purpose of this application is to provide an application
// that is using plain go to serve a rest api
//
//
// Terms Of Service:
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Roither | Seiberl <andi.roither@protonmail.com>
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
// swagger:meta
//go:generate swagger generate spec
package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	//createMockData()

	//fmt.Println("REST API started!")
	//r := registerRoutes()
	//log.Fatal(http.ListenAndServe(":8000", r))

	fmt.Println("> Snippets Rest API <")
	a := App{}

	// should be done by os.Getenv("")
	a.Initialize(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_HOST_NAME"), true)

	a.Run(":" + os.Getenv("REST_API_PORT"))
}

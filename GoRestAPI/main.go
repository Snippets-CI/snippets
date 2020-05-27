// Snippets
//
// The purpose of this application is to serve Snippets
//
//     Schemes: http
//     Host: localhost:8010
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Andreas | Markus <andi.roither@protonmail.com>
//
//     Consumes:
//     - text/plain
//
//     Produces:
//     - text/plain
//
// swagger:meta
//go:generate swagger generate spec
package main

import (
	"fmt"

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
	a.Initialize("admin", "123", "Snippets")

	a.Run(":8000")
}

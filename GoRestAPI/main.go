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

package main

import (
	"fmt"
	"net/http"

	"github.com/sowmik-sec/mongoapi/router"
)

// sowmiksec_db_user
// IGpzBIBEx3B9dYBN
// mongodb+srv://sowmiksec_db_user:IGpzBIBEx3B9dYBN@cluster0.ps3jgyx.mongodb.net/?appName=Cluster0

func main() {
	fmt.Println("MongoDB API")
	r := router.Router()
	fmt.Println("Server is getting started...")
	http.ListenAndServe(":4000", r)
	fmt.Println("Listening at port 4000...")
}

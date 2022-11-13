package main

import (
	"persion_test/router"
)

func main() {
	r := router.Router()
	r.Run(":8000")
}

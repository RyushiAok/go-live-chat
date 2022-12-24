package main

import "RealTime_Group_Chat/presen/router"

func main() {
	router := router.Init()
	router.Run("localhost:8080")
}

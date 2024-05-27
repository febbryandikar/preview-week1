package main

import (
	"fmt"
	"log"
	"preview-week1/config"
	"preview-week1/handler"
)

func main() {
	router, server := config.SetupServer()
	db := &handler.NewBranchHandler{DB: config.GetDatabase()}

	router.GET("/branches", db.GetAllBranches)
	router.GET("/branches/:id", db.GetBranchById)
	router.POST("/branches", db.CreateNewBranch)
	router.PUT("/branches/:id", db.UpdateBranch)
	router.DELETE("/branches/:id", db.DeleteBranch)

	fmt.Println("Server running on port :8080")
	log.Fatal(server.ListenAndServe())
}

package main

import (
 "fmt"
 "github.com/ilexsor/internal/utils"
)

func main(){
	fmt.Println("Go Final Project")

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	utils.RegisterHandlers(router)
}
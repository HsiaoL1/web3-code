package main

import "github.com/joho/godotenv"

func init(){
	
}

func main() {
	// Code
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

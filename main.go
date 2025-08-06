package main

import (
	"github.com/paul39-33/aggregator/internal/config"
	"fmt"
)

func main(){
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = cfg.SetUser("Paul")
	if err != nil {
		fmt.Println(err)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", cfg)
}
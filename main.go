package main

import (
	"github.com/paul39-33/aggregator/internal/config"
	"log"
	"os"
)

func main(){
	s := state{}
	c := commands{
		commandList: map[string]func(*state, command) error {
			"login":	handlerLogin,
		},
	}
	args := os.Args
	
	//check if there are enough arguments
	if len(args) < 2{
		log.Fatalf("Error not enough arguments!")
	}

	input := command {
		name:	args[1],
		arg:	args[2:],
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading: %v", err)
	}

	s.cfg = &cfg

	err = c.run(&s, input)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}
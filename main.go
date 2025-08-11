package main

import (
	"github.com/paul39-33/aggregator/internal/config"
	"log"
	"os"
	_ "github.com/lib/pq"
	"github.com/paul39-33/aggregator/internal/database"
	"database/sql"
)

func main(){
	s := state{}
	c := commands{
		commandList: map[string]func(*state, command) error {
			"login":	handlerLogin,
			"register":	handlerRegister,
			"reset":	handlerReset,
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

	//get db URL from environment variable
	dbURL := os.Getenv("GATOR_DATABASE_URL")
	cfg.Db_url = dbURL

	//load database URL to config struct and sql.Open() a connection to the database
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatalf("Error loading database to config: %v", err)
	}

	dbQueries := database.New(db)
	s.db = dbQueries

	err = c.run(&s, input)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}
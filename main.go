package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	curState := state{}
	myConfig := config.Read()
	curState.config = &myConfig
	comms := commands{commandList: make(map[string]func(*state, command) error)}
	comms.register("login", handlerLogin)
	comms.register("register", handlerRegister)
	comms.register("reset", handlerResetUsers)
	comms.register("users", handlerGetUsers)
	comms.register("agg", handlerAgg)
	comms.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	comms.register("feeds", handlerListFeeds)
	comms.register("follow", middlewareLoggedIn(handlerFollow))
	comms.register("following", handlerFollowing)
	comms.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	comms.register("browse", middlewareLoggedIn(handlerBrowse))
	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", myConfig.DbURL)
	if err != nil {
		fmt.Println("Error connecting to Database")
		os.Exit(1)
	}

	curState.db = database.New(db)

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	command := command{
		name:      cmdName,
		arguments: cmdArgs,
	}
	err = comms.run(&curState, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

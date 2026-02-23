package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joliverstrom-cmd/gator_boot/internal/config"
	"github.com/joliverstrom-cmd/gator_boot/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	myConfig, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Something went wrong: %+v", err)
	}

	db, err := sql.Open("postgres", myConfig.DbURL)
	if err != nil {
		log.Fatalf("Problem opening the database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{
		db:  dbQueries,
		cfg: &myConfig,
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggs)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", middlewareLoggedIn(handlerFeeds))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	inputArgs := os.Args
	if len(inputArgs) < 2 {
		log.Fatalf("Not enough arguments provided")
		os.Exit(1)
	}

	err = cmds.run(programState, command{inputArgs[1], inputArgs[2:]})
	if err != nil {
		log.Fatalf("Error: %v", err)
		os.Exit(1)
	}

	myConfig, err = config.ReadConfig()
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}

}

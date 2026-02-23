package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/gator_boot/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username supplied")
	}

	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("That username does not exist")
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("Couldn't set username: %v", err)
	}

	fmt.Printf("Username has been set to: %v\n", userName)
	return nil

}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No name supplied, usage: register <name>")
	}

	nameInput := cmd.args[0]

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      nameInput,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Couldn't add to database: %v", err)
	}

	err = s.cfg.SetUser(nameInput)
	if err != nil {
		return fmt.Errorf("Couldn't update the config json: %v", err)
	}
	fmt.Printf("Entered row with this name: %v\n", nameInput)
	printUser(user)

	return nil

}

func handlerReset(s *state, cmd command) error {

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't delete database entries: %v", err)
	}

	fmt.Println("Users database has been reset")
	return nil

}

func handlerUsers(s *state, cmd command) error {

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't delete database entries: %v", err)
	}

	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}

	}

	return nil

}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
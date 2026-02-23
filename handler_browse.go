package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/joliverstrom-cmd/gator_boot/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {

	limit := 2

	if len(cmd.args) > 0 {
		wantedLimit, err := strconv.Atoi(cmd.args[0])
		limit = wantedLimit
		if err != nil {
			return fmt.Errorf("Input needs to be a number (1, 2 3, ...)")
		}
	}

	userposts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Couldn't get posts from DB: %w", err)
	}

	for _, post := range userposts {
		fmt.Printf("* %v, URL: %v\n", post.Title, post.Url)
	}

	return nil

}

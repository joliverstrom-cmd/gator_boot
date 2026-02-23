package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	command, exists := c.commandMap[cmd.name]

	if !exists {
		return fmt.Errorf("Your supplied command <%v> doesn't exist", cmd.name)
	}

	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

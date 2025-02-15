package main

import (
	"gator/internal/config"
	"gator/internal/database"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

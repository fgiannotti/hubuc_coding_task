package main

import (
	"errors"
	"fmt"
)

var UserNotFoundError = func(username string) error { return errors.New(fmt.Sprintf("User %s not found", username)) }

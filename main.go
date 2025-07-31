// main.go

package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/signal"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to the Monkey programming language REPL.\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	sigChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		fmt.Println("\nExiting gracefully. Goodbye!")
		done <- true
	}()

	go repl.Start(os.Stdin, os.Stdout)

	<-done
}

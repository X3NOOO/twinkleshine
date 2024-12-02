package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

var envFlag string
var verboseFlag bool

func main() {
	flag.StringVar(&envFlag, "env", ".env", "Environment file")
	flag.BoolVar(&verboseFlag, "verbose", false, "Verbose mode")
	flag.Parse()

	var logfile io.Writer
	if verboseFlag {
		logfile = os.Stderr
	} else {
		logfile = io.Discard
	}

	log.SetOutput(logfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_ = godotenv.Load(envFlag)

	bot, err := NewBot(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Cannot create the bot: %v", err)
	}

	err = bot.Run()
	if err != nil {
		log.Fatalf("Error running bot: %v", err)
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	stop_signal := make(chan os.Signal, 1)
	signal.Notify(stop_signal, os.Interrupt)
	<-stop_signal

	fmt.Println("Shutting down...")
	err = bot.Stop()
	if err != nil {
		log.Fatalf("Error stopping bot: %v", err)
	}
	fmt.Println("\nBye!")
}

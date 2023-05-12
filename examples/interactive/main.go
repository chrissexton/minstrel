package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/chrissexton/minstrel"
	"log"
	"os"
)

var token = flag.String("token", "", "GCP token for your project")
var project = flag.String("project", "", "CGP project")
var prompt = flag.String("prompt", "You are a cool dude that likes things. Answer the following:", "Initial prompt to use")

func main() {
	flag.Parse()
	if *token == "" {
		log.Printf("You must provide a token with -token")
		flag.Usage()
	}
	if *project == "" {
		log.Printf("You must provide a project with -project")
		flag.Usage()
	}
	m := minstrel.New(*token, *project)
	m.SetPrompt(*prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("input text: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		out, err := m.Complete(line)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Bot: " + out)
	}
}

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/chrissexton/minstrel"
	"log"
	"os"
	"os/exec"
	"strings"
)

var token = flag.String("token", "", "GCP token for your project")
var project = flag.String("project", "", "CGP project")
var prompt = flag.String("prompt", "You are a cool dude that likes things. Answer the following:", "Initial prompt to use")

func main() {
	var err error
	flag.Parse()
	if *project == "" {
		log.Printf("You must provide a project with -project")
		flag.Usage()
	}

	if *token == "" {
		*token, err = getToken()
		if err != nil {
			log.Fatal(err)
		}
	}

	m := minstrel.New(*token, *project)
	m.SetPrompt(*prompt)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("input text: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: " + err.Error())
		}
		out, err := m.Complete(line)
		if err != nil {
			log.Fatalf("Error: " + err.Error())
		}
		fmt.Println("Bot: " + out)
	}
}

// getToken - because eff it, getting tokens is outside the scope of this program
func getToken() (string, error) {
	cmd := exec.Command("gcloud", "auth", "print-access-token")
	sout := bytes.NewBufferString("")
	cmd.Stdout = sout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(sout.String()), nil
}

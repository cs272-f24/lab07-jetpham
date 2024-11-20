package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	setup, err := newSetup()
	if err != nil {
		fmt.Println("Error during setup:", err)
		return
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter prompt: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)
		if prompt == "exit" {
			break
		}
		output := toolCallingAgent(setup, prompt)
		fmt.Println(output)
	}
}

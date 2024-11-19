package main

import (
	"bufio"
	"encoding/json"
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
		courses, filter, err := getCoursesFromPrompt(&setup, prompt)
		if err != nil {
			fmt.Println("Error getting courses:", err)
			continue
		}
		fmt.Println(filter)
		fmt.Printf("Found %d courses: \n", len(courses))
		for _, course := range courses {
			fmt.Println(course)
		}
		coursesJSON, err := json.MarshalIndent(courses, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling courses to JSON:", err)
			continue
		}
		output, err := setup.openAIClient.CreateCompletion(prompt + "\n" + string(coursesJSON))
		if err != nil {
			fmt.Println("Error creating completion:", err)
		}
		fmt.Println(output)
	}
}

package main

import (
	"fmt"
	"log"
	"time"
)

func getCourses(setup *Setup, prompt string) ([]Course, error) {
	start := time.Now()
	filterPrompt := `
		extract information from the prompt int o course filter. Only inlude information that is explicitly in the prompt and not inferred.
		`
	courseFilter := setup.openAIClient.GetCourseFilter(prompt, filterPrompt)
	log.Printf("Original %s", courseFilter)

	correctedFilter := setup.chromaDB.correctCourseFilter(setup.collections, courseFilter)
	log.Printf("Corrected %s", correctedFilter)

	filteredCourses, err := filterCourses(setup.sqlDB, correctedFilter)
	if err != nil {
		fmt.Println("Error filtering courses:", err)
		return nil, err
	}

	log.Printf("Found %v courses in %f.2", len(filteredCourses), time.Since(start).Seconds())

	return filteredCourses, nil
}

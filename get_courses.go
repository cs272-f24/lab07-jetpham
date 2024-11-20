package main

import (
	"fmt"
	"log"
)

func getCourses(setup *Setup, prompt string) ([]Course, error) {
	enhancingPrompt := `
		phrase this into a prompt that is asking about courses for example:
		"Can I learn guitar this semester?" would be phrased as "What courses are there about guitar?"
		"Where does Bioinformatics meet?" would be phrased as "What courses are there for Bioinformatics?"
		"I would like to take a Rhetoric course from Phil Choong. What can I take?" would be phrased as "What courses are there for Rhetoric with Phil Choong?"
		"Can I take a course on the weekends?" would be phrased as "What courses are there that meet on the weekends?"
		`
	enhancedPrompt, err := setup.openAIClient.CreateCompletion(prompt, enhancingPrompt)

	if err != nil {
		return nil, err
	}
	log.Println(enhancedPrompt)
	filterPrompt := `
		extract information from the prompt int o course filter. Only inlude information that is explicitly in the prompt and not inferred.
		`
	courseFilter := setup.openAIClient.GetCourseFilter(enhancedPrompt, filterPrompt)
	log.Printf("Original %s", courseFilter)

	correctedFilter := setup.chromaDB.correctCourseFilter(setup.collections, courseFilter)
	log.Printf("Corrected %s", correctedFilter)

	filteredCourses, err := filterCourses(setup.sqlDB, correctedFilter)
	if err != nil {
		fmt.Println("Error filtering courses:", err)
		return nil, err
	}
	return filteredCourses, nil
}

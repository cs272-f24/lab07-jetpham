package main

import (
	"fmt"
	"log"
)

func getCoursesFromPrompt(setup *Setup, prompt string) ([]Course, CourseFilter, error) {
	filteringPrompt := `
	Convert the following message into a prompt that goes into a program that filters over courses in a registry. It's important to specify what the parts of the prompt there are. Knowing the subject of the class if specified, knowing the title of the class if specified. If the prompt just talks about a subject, put it into the Title Short Desc\n`
	enhancedPrompt, err := setup.openAIClient.CreateCompletion(filteringPrompt + prompt)

	if err != nil {
		return nil, CourseFilter{}, err
	}
	log.Println(enhancedPrompt)
	courseFilter := setup.openAIClient.GetCourseFilter(enhancedPrompt)

	correctedFilter := setup.chromaDB.correctCourseFilter(setup.collections, courseFilter)

	filteredCourses, err := filterCourses(setup.sqlDB, correctedFilter)
	if err != nil {
		fmt.Println("Error filtering courses:", err)
		return nil, CourseFilter{}, err
	}
	return filteredCourses, correctedFilter, nil
}

package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Setup struct {
	courses      []Course
	openAIClient *OpenAIClient
	chromaDB     *chromaDB
	sqlDB        *gorm.DB
	collections  *collections
}

func newSetup() (Setup, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return Setup{}, err
	}
	courses, err := loadCSV("/home/jet/Documents/cs272/project05-jetpham/Fall 2024 Class Schedule.csv")
	if err != nil {
		fmt.Println("Error loading CSV file:", err)
		return Setup{}, err
	}
	openAIClient := NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))
	if openAIClient == nil {
		fmt.Println("Error creating OpenAI client")
		return Setup{}, err
	}
	chromaDB, err := newChroma()
	if err != nil {
		fmt.Println("Error setting up ChromaDB:", err)
		return Setup{}, err
	}
	sqlDB, err := newSqlite(courses)
	if err != nil {
		fmt.Println("Error setting up SQLite database:", err)
		return Setup{}, err
	}
	collections, err := makeCollections(chromaDB, courses) // Assuming you have a function to load collections
	if err != nil {
		fmt.Println("Error loading collections:", err)
		return Setup{}, err
	}
	return Setup{
		courses:      courses,
		openAIClient: openAIClient,
		chromaDB:     chromaDB,
		sqlDB:        sqlDB,
		collections:  collections,
	}, nil
}

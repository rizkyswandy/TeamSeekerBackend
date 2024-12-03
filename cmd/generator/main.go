package main

import (
	"github.com/rizkyswandy/TeamSeekerBackend/api"
	"github.com/rizkyswandy/TeamSeekerBackend/internal/database/postgres"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var (
	faculties = []string{
		"Computer Science",
		"Engineering",
		"Business",
		"Arts",
		"Science",
		"Medicine",
	}

	fieldOfStudy = map[string][]string{
		"Computer Science": {"Software Engineering", "Data Science", "Cybersecurity", "AI & Machine Learning", "Network Systems"},
		"Engineering":      {"Mechanical", "Electrical", "Civil", "Chemical", "Industrial"},
		"Business":        {"Marketing", "Finance", "Management", "Accounting", "International Business"},
		"Arts":            {"Fine Arts", "Design", "Music", "Theater", "Film Studies"},
		"Science":         {"Physics", "Chemistry", "Biology", "Mathematics", "Environmental Science"},
		"Medicine":        {"General Medicine", "Dentistry", "Pharmacy", "Nursing", "Public Health"},
	}

	skills = []string{
		"Python", "Java", "Go", "JavaScript", "C++",
		"React", "Node.js", "Docker", "Kubernetes", "AWS",
		"Data Analysis", "Machine Learning", "UI/UX Design",
		"Project Management", "Team Leadership", "Communication",
		"Problem Solving", "Critical Thinking", "Research",
		"Public Speaking", "Technical Writing",
	}

	focus = []string{
		"Web Development", "Mobile Development", "Cloud Computing",
		"Data Science", "Artificial Intelligence", "Cybersecurity",
		"DevOps", "Systems Design", "UI/UX", "Product Management",
		"Research", "Innovation", "Entrepreneurship",
	}
)

func generateRandomProfile(id int) *api.StudentProfile {
	rand.Seed(time.Now().UnixNano() + int64(id))

	faculty := faculties[rand.Intn(len(faculties))]
	
	possibleFields := fieldOfStudy[faculty]
	field := possibleFields[rand.Intn(len(possibleFields))]

	numSkills := rand.Intn(5) + 3
	selectedSkills := make([]string, 0)
	for i := 0; i < numSkills; i++ {
		skill := skills[rand.Intn(len(skills))]
		if !contains(selectedSkills, skill) {
			selectedSkills = append(selectedSkills, skill)
		}
	}

	numFocus := rand.Intn(3) + 2
	selectedFocus := make([]string, 0)
	for i := 0; i < numFocus; i++ {
		f := focus[rand.Intn(len(focus))]
		if !contains(selectedFocus, f) {
			selectedFocus = append(selectedFocus, f)
		}
	}

	return &api.StudentProfile{
		Name:         fmt.Sprintf("Student %d", id),
		Email:        fmt.Sprintf("student%d@university.edu", id),
		Faculty:      faculty,
		FieldOfStudy: field,
		Semester:     rand.Intn(8) + 1, 
		Skills:       selectedSkills,
		Focus:        selectedFocus,
		IsAvailable:  rand.Float32() > 0.3, 
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	connString := "postgres://ilb:@localhost:5432/team_seeker?sslmode=disable"
	
	db, err := postgres.NewPostgresDB(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	for i := 1; i <= 10000; i++ {
		profile := generateRandomProfile(i)
		
		err := db.CreateProfile(profile)
		if err != nil {
			log.Printf("Failed to create profile %d: %v", i, err)
			continue
		}
		
		if i%50 == 0 {
			log.Printf("Created %d profiles", i)
		}
	}

	log.Println("Data generation completed!")
}

package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/invopop/jsonschema"
)

type CourseFilter struct {
	SubjectCodes            []string `json:"subject_codes" jsonschema_description:"The subject codes of the courses, if provided as codes. Example: ['CS', 'MATH']"`
	SubjectNames            []string `json:"subject_names" jsonschema_description:"The subject names of the courses, if provided as names. Example: ['Computer Science', 'Mathematics']"`
	CourseNumbers           []string `json:"course_numbers" jsonschema_description:"The course numbers. Example: ['101', '202']"`
	Sections                []string `json:"sections" jsonschema_description:"The sections of the courses. Example: ['001', '002']"`
	CRNs                    []int    `json:"crns" jsonschema_description:"The course registration numbers. Example: [12345, 67890]"`
	ScheduleTypeCodes       []string `json:"schedule_type_codes" jsonschema_description:"The schedule type codes, if provided as codes. Example: ['LEC', 'LAB']"`
	CampusCodes             []string `json:"campus_codes" jsonschema_description:"The campus codes, if provided as codes. Example: ['MAIN', 'SAT']"`
	TitleShortDescs         []string `json:"title_short_descs" jsonschema_description:"The short descriptions of the course titles. Example: ['Intro to CS', 'Calculus I']"`
	InstructionModeDescs    []string `json:"instruction_mode_descs" jsonschema_description:"The descriptions of the instruction modes. Example: ['In Person', 'Online']"`
	MeetingTypeCodes        []string `json:"meeting_type_codes" jsonschema_description:"The meeting type codes, if provided as codes. Example: ['CLAS', 'LAB']"`
	MeetingTypeNames        []string `json:"meeting_type_names" jsonschema_description:"The meeting type names, if provided as names. Example: ['Class', 'Laboratory']"`
	MeetDays                []string `json:"meet_days" jsonschema_description:"The meeting days. Example: ['M', 'W', 'F']"`
	MeetDaysFull            []string `json:"meet_days_full" jsonschema_description:"The full meeting days. Example: ['Monday', 'Wednesday', 'Friday']"`
	BeginTimes              []string `json:"begin_times" jsonschema_description:"The begin times of the meetings. Example: ['08:00', '10:00']"`
	EndTimes                []string `json:"end_times" jsonschema_description:"The end times of the meetings. Example: ['09:00', '11:00']"`
	MeetStarts              []string `json:"meet_starts" jsonschema_description:"The start dates of the meetings. Example: ['2023-01-10', '2023-01-12']"`
	MeetEnds                []string `json:"meet_ends" jsonschema_description:"The end dates of the meetings. Example: ['2023-05-10', '2023-05-12']"`
	Buildings               []string `json:"buildings" jsonschema_description:"The buildings where the courses are held. Example: ['Engineering Hall', 'Science Building']"`
	Rooms                   []string `json:"rooms" jsonschema_description:"The rooms where the courses are held. Example: ['101', '202']"`
	ActualEnrollments       []int    `json:"actual_enrollments" jsonschema_description:"The actual enrollments in the courses. Example: [30, 25]"`
	PrimaryInstructorFirsts []string `json:"primary_instructor_first_names" jsonschema_description:"The first names of the primary instructors. Only fill out if full name not given. Example: ['John', 'Jane']"`
	PrimaryInstructorLasts  []string `json:"primary_instructor_last_names" jsonschema_description:"The last names of the primary instructors. Only if exlicitally told it's a last name. If given both first and last then fill out full name. Example: ['Doe', 'Smith']"`
	PrimaryInstructorFulls  []string `json:"primary_instructor_full_names" jsonschema_description:"The full names of the primary instructors. Only for the full names of the instructors. Example: ['John Doe', 'Jane Smith']"`
	PrimaryInstructorEmails []string `json:"primary_instructor_emails" jsonschema_description:"The emails of the primary instructors. Example: ['jdoe@example.com', 'jsmith@example.com']"`
	Colleges                []string `json:"colleges" jsonschema_description:"The colleges offering the courses. Example: ['College of Engineering', 'College of Science']"`
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func (d *chromaDB) correctCourseFilter(collections *collections, filter CourseFilter) CourseFilter {
	filter.SubjectNames = d.correctSubjectNames(collections, filter)
	filter.TitleShortDescs = d.correctTitleShortDescs(collections, filter)
	filter.PrimaryInstructorFirsts = d.correctInstructorFirstNames(collections, filter)
	filter.PrimaryInstructorLasts = d.correctInstructorLastNames(collections, filter)
	filter.PrimaryInstructorFulls = d.correctInstructorFullNames(collections, filter)
	return filter
}

func (d *chromaDB) correctSubjectNames(collections *collections, filter CourseFilter) []string {
	tempSubjectNames := make([]string, 0)
	for _, code := range filter.SubjectCodes {
		if name, exists := courseSubjects[code]; exists {
			tempSubjectNames = append(tempSubjectNames, name)
		}
	}

	// Query the collections for subject names
	correctedSubjectNames := make([]string, 0, len(filter.SubjectNames))
	for _, subjectName := range filter.SubjectNames {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.SubjectNameList, subjectName) {
			correctedSubjectNames = append(correctedSubjectNames, subjectName)
			continue
		}
		result, err := d.query(collections.SubjectNameCollection.Name, subjectName, 1)
		if err != nil {
			log.Fatalf("Error querying subject name '%s': %s \n", subjectName, err)
		}
		correctedSubjectNames = append(correctedSubjectNames, result...)
	}

	// Add the converted subject names to the corrected ones
	correctedSubjectNames = append(correctedSubjectNames, tempSubjectNames...)

	return correctedSubjectNames
}

func (d *chromaDB) correctTitleShortDescs(collections *collections, filter CourseFilter) []string {
	correctedTitleShortDescs := make([]string, 0, len(filter.TitleShortDescs))
	for _, titleShortDesc := range filter.TitleShortDescs {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.TitleShortDescList, titleShortDesc) {
			correctedTitleShortDescs = append(correctedTitleShortDescs, titleShortDesc)
			continue
		}
		result, err := d.query(collections.TitleShortDescCollection.Name, titleShortDesc, 1)
		if err != nil {
			log.Fatalf("Error querying title short description '%s': %s \n", titleShortDesc, err)
		}
		correctedTitleShortDescs = append(correctedTitleShortDescs, result...)
	}
	return correctedTitleShortDescs
}

func (d *chromaDB) correctInstructorFirstNames(collections *collections, filter CourseFilter) []string {
	correctedFirstNames := make([]string, 0, len(filter.PrimaryInstructorFirsts))
	for _, firstName := range filter.PrimaryInstructorFirsts {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.InstructorFirstNameList, firstName) {
			correctedFirstNames = append(correctedFirstNames, firstName)
			log.Println("Found in collection: ", firstName)
			continue
		}
		log.Println("Not found in collection: ", firstName)
		result, err := d.query(collections.instructorFirstNameCollection.Name, firstName, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor first name '%s': %s \n", firstName, err)
		}
		log.Println("Result: ", result)
		correctedFirstNames = append(correctedFirstNames, result...)
	}
	return correctedFirstNames
}

func (d *chromaDB) correctInstructorLastNames(collections *collections, filter CourseFilter) []string {
	correctedLastNames := make([]string, 0, len(filter.PrimaryInstructorLasts))
	for _, lastName := range filter.PrimaryInstructorLasts {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.InstructorLastNameList, lastName) {
			correctedLastNames = append(correctedLastNames, lastName)
			continue
		}
		result, err := d.query(collections.instructorLastNameCollection.Name, lastName, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor last name '%s': %s \n", lastName, err)
		}
		correctedLastNames = append(correctedLastNames, result...)
	}
	return correctedLastNames
}

func (d *chromaDB) correctInstructorFullNames(collections *collections, filter CourseFilter) []string {
	correctedFullNames := make([]string, 0, len(filter.PrimaryInstructorFulls))
	for _, fullName := range filter.PrimaryInstructorFulls {
		// Shortcut if it's in the collection already
		if slices.Contains(collections.InstructorFullNameList, fullName) {
			correctedFullNames = append(correctedFullNames, fullName)
			continue
		}
		result, err := d.query(collections.InstructorFullNameCollection.Name, fullName, 1)
		if err != nil {
			log.Fatalf("Error querying primary instructor full name '%s': %s \n", fullName, err)
		}
		correctedFullNames = append(correctedFullNames, result...)
	}
	return correctedFullNames
}

func (f CourseFilter) String() string {
	var sb strings.Builder
	sb.WriteString("CourseFilter:\n")
	if len(f.SubjectCodes) > 0 {
		sb.WriteString(fmt.Sprintf("  SubjectCodes: %v\n", f.SubjectCodes))
	}
	if len(f.SubjectNames) > 0 {
		sb.WriteString(fmt.Sprintf("  SubjectNames: %v\n", f.SubjectNames))
	}
	if len(f.CourseNumbers) > 0 {
		sb.WriteString(fmt.Sprintf("  CourseNumbers: %v\n", f.CourseNumbers))
	}
	if len(f.Sections) > 0 {
		sb.WriteString(fmt.Sprintf("  Sections: %v\n", f.Sections))
	}
	if len(f.CRNs) > 0 {
		sb.WriteString(fmt.Sprintf("  CRNs: %v\n", f.CRNs))
	}
	if len(f.ScheduleTypeCodes) > 0 {
		sb.WriteString(fmt.Sprintf("  ScheduleTypeCodes: %v\n", f.ScheduleTypeCodes))
	}
	if len(f.CampusCodes) > 0 {
		sb.WriteString(fmt.Sprintf("  CampusCodes: %v\n", f.CampusCodes))
	}
	if len(f.TitleShortDescs) > 0 {
		sb.WriteString(fmt.Sprintf("  TitleShortDescs: %v\n", f.TitleShortDescs))
	}
	if len(f.InstructionModeDescs) > 0 {
		sb.WriteString(fmt.Sprintf("  InstructionModeDescs: %v\n", f.InstructionModeDescs))
	}
	if len(f.MeetingTypeCodes) > 0 {
		sb.WriteString(fmt.Sprintf("  MeetingTypeCodes: %v\n", f.MeetingTypeCodes))
	}
	if len(f.MeetingTypeNames) > 0 {
		sb.WriteString(fmt.Sprintf("  MeetingTypeNames: %v\n", f.MeetingTypeNames))
	}
	if len(f.MeetDays) > 0 {
		sb.WriteString(fmt.Sprintf("  MeetDays: %v\n", f.MeetDays))
	}
	if len(f.MeetDaysFull) > 0 {
		sb.WriteString(fmt.Sprintf("  MeetDaysFull: %v\n", f.MeetDaysFull))
	}
	if len(f.BeginTimes) > 0 {
		sb.WriteString(fmt.Sprintf("  BeginTimes: %v\n", f.BeginTimes))
	}
	if len(f.EndTimes) > 0 {
		sb.WriteString(fmt.Sprintf("  EndTimes: %v\n", f.EndTimes))
	}
	if len(f.MeetStarts) > 0 {
		sb.WriteString(fmt.Sprintf("  MeetStarts: %v\n", f.MeetStarts))
	}
	if len(f.MeetEnds) > 0 {
		sb.WriteString(fmt.Sprintf("  MeetEnds: %v\n", f.MeetEnds))
	}
	if len(f.Buildings) > 0 {
		sb.WriteString(fmt.Sprintf("  Buildings: %v\n", f.Buildings))
	}
	if len(f.Rooms) > 0 {
		sb.WriteString(fmt.Sprintf("  Rooms: %v\n", f.Rooms))
	}
	if len(f.ActualEnrollments) > 0 {
		sb.WriteString(fmt.Sprintf("  ActualEnrollments: %v\n", f.ActualEnrollments))
	}
	if len(f.PrimaryInstructorFirsts) > 0 {
		sb.WriteString(fmt.Sprintf("  PrimaryInstructorFirsts: %v\n", f.PrimaryInstructorFirsts))
	}
	if len(f.PrimaryInstructorLasts) > 0 {
		sb.WriteString(fmt.Sprintf("  PrimaryInstructorLasts: %v\n", f.PrimaryInstructorLasts))
	}
	if len(f.PrimaryInstructorFulls) > 0 {
		sb.WriteString(fmt.Sprintf("  PrimaryInstructorFulls: %v\n", f.PrimaryInstructorFulls))
	}
	if len(f.PrimaryInstructorEmails) > 0 {
		sb.WriteString(fmt.Sprintf("  PrimaryInstructorEmails: %v\n", f.PrimaryInstructorEmails))
	}
	if len(f.Colleges) > 0 {
		sb.WriteString(fmt.Sprintf("  Colleges: %v\n", f.Colleges))
	}
	return sb.String()
}

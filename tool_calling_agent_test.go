package main

import "testing"

func Test_toolCallingAgent(t *testing.T) {
	setup, err := newSetup()
	if err != nil {
		t.Errorf("Error during setup: %v", err)
		return
	}
	tests := []struct {
		name    string
		prompt  string
		want    string
		wantErr bool
	}{
		{
			name:   "TestPhil",
			prompt: "What courses is Phil Peterson teaching in Fall 2024?",
			want: `
				In Fall 2024, Phil Peterson is teaching the following courses:	
				1. **CS 272: Software Development**
				- **Section 03**
					- Schedule: Tuesday and Thursday, 2:40 PM - 4:25 PM
					- Location: Building LS, Room G12
					- Enrollment: 26 students

				- **Section 04**
					- Schedule: Tuesday and Thursday, 8:00 AM - 9:45 AM
					- Location: Building LS, Room G12
					- Enrollment: 19 students

				2. **CS 272L: Software Development Lab**
				- **Section 01**
					- Schedule: Wednesday, 1:00 PM - 2:30 PM
					- Location: Building MH, Room 122
					- Enrollment: 21 students

				- **Section 02**
					- Schedule: Wednesday, 2:55 PM - 4:25 PM
					- Location: Building MH, Room 122
					- Enrollment: 24 students

				All classes are conducted in-person.
				`,
			wantErr: false,
		},
		{
			name:   "TestPHIL",
			prompt: "Which philosophy courses are offered this semester?",
			want: `
				Here are the philosophy courses offered this semester:

				1. **Great Philosophical Questions**
				- Instructor: Deena Lin
				- Schedule: Monday, Wednesday, Friday (9:15 AM - 10:20 AM)
				- Format: In-Person

				2. **Great Philosophical Questions**
				- Instructor: Jea Oh
				- Schedule: Monday, Wednesday, Friday (1:00 PM - 2:05 PM)
				- Format: In-Person

				3. **FYS: Lovers of Wisdom**
				- Instructor: Thomas Cavanaugh
				- Schedule: Tuesday, Thursday (8:00 AM - 9:45 AM)
				- Format: In-Person

				4. **Philosophy of Religion**
				- Instructor: Deena Lin
				- Schedule: Monday, Wednesday, Friday (1:00 PM - 2:05 PM)
				- Format: In-Person

				5. **Philosophy of Science**
				- Instructor: Krupa Patel
				- Schedule: Monday, Wednesday (4:45 PM - 6:25 PM)
				- Format: In-Person

				6. **Philosophy of Biology**
				- Instructor: Stephen Friesen
				- Schedule: Monday, Wednesday, Friday (10:30 AM - 11:35 AM)
				- Format: In-Person

				7. **The Human Animal**
				- Instructor: Jennifer Fisher
				- Schedule: Tuesday, Thursday (2:40 PM - 4:25 PM)
				- Format: In-Person

				8. **Aesthetics**
				- Instructor: Laurel Scotland-Stewart
				- Schedule: Monday, Wednesday, Friday (2:15 PM - 3:20 PM)
				- Format: In-Person

				9. **Asian Philosophy**
				- Instructor: Joshua Stoll
				- Schedule: Monday, Wednesday (4:45 PM - 6:25 PM)
				- Format: In-Person

				10. **Ethics**
					- Instructor: Joshua Carboni
					- Schedule: Tuesday, Thursday (2:40 PM - 4:25 PM)
					- Format: In-Person

				11. **Environmental Ethics**
					- Instructor: Stephen Friesen
					- Schedule: Monday, Wednesday, Friday (9:15 AM - 10:20 AM)
					- Format: In-Person

				12. **Mind, Freedom & Knowledge**
					- Instructor: Jennifer Fisher
					- Schedule: Tuesday, Thursday (4:35 PM - 6:20 PM)
					- Format: In-Person

				13. **Existentialism**
					- Instructor: Brian Pines
					- Schedule: Tuesday, Thursday (9:55 AM - 11:40 AM)
					- Format: In-Person

				14. **Ancient & Medieval Philosophy**
					- Instructor: Thomas Cavanaugh
					- Schedule: Tuesday, Thursday (9:55 AM - 11:40 AM)
					- Format: In-Person

				15. **Logic**
					- Instructor: Nick Leonard
					- Schedule: Monday, Wednesday, Friday (11:45 AM - 12:50 PM)
					- Format: In-Person

				16. **Topics in Contemporary Philosophy**
					- Instructor: David Kim
					- Schedule: Tuesday, Thursday (12:45 PM - 2:30 PM)
					- Format: In-Person

				These courses span a variety of philosophical fields and approaches, from ethics and existentialism to philosophy of science and ancient texts. Please let me know if you need more detailed information about any specific course!
				`,
			wantErr: false,
		},
		{
			name:   "TestBio",
			prompt: "Where does Bioinformatics meet?",
			want: `
				The Bioinformatics course meets on Monday and Wednesday from 4:45 PM to 6:25 PM. There are two sections:

				1. **Section 01** meets in the KA building, room 111.
				2. **Section 02** meets in the HR building, room 136.

				The course is part of the Biotechnology program and the primary instructor for the first section is Naupaka Zimmerman.
				`,
			wantErr: false,
		},
		{
			name:   "TestGuitar",
			prompt: "Can I learn guitar this semester?",
			want: `
				Yes, you can learn guitar this semester by enrolling in "Guitar and Bass Lessons" courses. Here are the available sections:

				1. **Guitar and Bass Lessons - Section 01**
				- Instructor: Christopher Ruscoe
				- Dates: 8/20/24 - 11/28/24
				- Location: To Be Announced (TBA)
				- Email: [cgruscoe@usfca.edu](mailto:cgruscoe@usfca.edu)
				- Current Enrollment: 5 students

				2. **Guitar and Bass Lessons - Section 02**
				- Instructor: Christopher Ruscoe
				- Dates: 8/20/24 - 11/28/24
				- Location: To Be Announced (TBA)
				- Email: [cgruscoe@usfca.edu](mailto:cgruscoe@usfca.edu)
				- Current Enrollment: 2 students

				These classes will provide the opportunity to learn guitar in a structured setting.
				`,
			wantErr: false,
		},
		{
			name:   "TestMultiple",
			prompt: "I would like to take a Rhetoric course from Phil Choong. What can I take?",
			want: `
				Philip Choong is instructing several Rhetoric courses that you can take. Here are the options:

				1. **Public Speaking (RHET 103, Section 05)**
				- **Schedule:** Monday, Wednesday, Friday
				- **Time:** 10:30 AM - 11:35 AM
				- **Location:** LM Building, Room 346A
				- **Dates:** August 20, 2024, to December 4, 2024

				2. **Public Speaking (RHET 103, Section 26)**
				- **Schedule:** Monday, Wednesday, Friday
				- **Time:** 11:45 AM - 12:50 PM
				- **Location:** LM Building, Room 346A
				- **Dates:** August 20, 2024, to December 4, 2024

				3. **FYS: Podcasts: Eloquentia & Aud (RHET 195, Section 02)**
				- **Schedule:** Monday, Wednesday, Friday
				- **Time:** 2:15 PM - 3:20 PM
				- **Location:** LM Building, Room 352
				- **Dates:** August 20, 2024, to December 4, 2024

				4. **Speaking Center Internship (RHET 328, Section 01)**
				- **Schedule:** Tuesday
				- **Time:** 4:35 PM - 6:25 PM
				- **Location:** LM Building, Room 345
				- **Dates:** August 20, 2024, to December 3, 2024

				If you're interested in any course for further details, let me know!
				`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toolCallingAgent(setup, tt.prompt); isSimilar(setup, got, tt.want) {
				t.Errorf("toolCallingAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}

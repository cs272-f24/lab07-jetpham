package main

import (
	"reflect"
	"testing"
)

func Test_getCoursesFromPrompt(t *testing.T) {
	setup, err := newSetup()
	if err != nil {
		t.Errorf("Error during setup: %v", err)
		return
	}
	tests := []struct {
		name    string
		prompt  string
		want    []int
		wantErr bool
	}{
		{
			name:    "TestPhil",
			prompt:  "What courses is Phil Peterson teaching in Fall 2024?",
			want:    []int{40646, 40647, 42343, 42344, 40146, 40166, 40215, 42533},
			wantErr: false,
		},
		{
			name:    "TestPHIL",
			prompt:  "Which philosophy courses are offered this semester?",
			want:    []int{41163, 41164, 41165, 41166, 41166, 41167, 41168, 41170, 41172, 41174, 41175, 41176, 42008, 41177, 41178, 41179, 41182, 41183, 41184, 41185, 41187, 41188, 41188, 41181, 41192, 41196, 41197, 41198, 41200, 41201, 42523, 41202, 41203},
			wantErr: false,
		},
		{
			name:    "TestBio",
			prompt:  "Where does Bioinformatics meet?",
			want:    []int{40548, 42323},
			wantErr: false,
		},
		{
			name:    "TestGuitar",
			prompt:  "Can I learn guitar this semester?",
			want:    []int{41140, 41141},
			wantErr: false,
		},
		{
			name:    "TestMultiple",
			prompt:  "I would like to take a Rhetoric course from Phil Choong. What can I take?",
			want:    []int{40146, 40166, 40215, 42533},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCourses(&setup, tt.prompt)
			CRNs := make([]int, 0, len(got))
			for _, course := range got {
				CRNs = append(CRNs, course.CRN)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getCoursesFromPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(CRNs, tt.want) {
				t.Errorf("getCoursesFromPrompt() got = %v, want %v", CRNs, tt.want)
			}
		})
	}
}

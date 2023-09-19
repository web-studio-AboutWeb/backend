package domain

import (
	"testing"

	"github.com/stretchr/testify/require"

	"web-studio-backend/internal/pkg/strhelp"
)

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		p         *Project
	}{
		{
			name:      "empty structure",
			wantError: true,
			p:         &Project{},
		},
		{
			name:      "empty title",
			wantError: true,
			p:         &Project{Description: "something"},
		},
		{
			name:      "description is too long",
			wantError: true,
			p: &Project{
				Title: "title",
				Description: (func() string {
					s, _ := strhelp.GenerateRandomString(2049)
					return s
				})(),
			},
		},
		{
			name:      "invalid link",
			wantError: true,
			p: &Project{
				Title: "title",
				Link:  "qwefqwef",
			},
		},
		{
			name:      "invalid link #2",
			wantError: true,
			p: &Project{
				Title: "title",
				Link:  "www.mysite.com",
			},
		},
		{
			name:      "valid link #1",
			wantError: false,
			p: &Project{
				Title: "title",
				Link:  "/relative/path/1",
			},
		},
		{
			name:      "valid link #2",
			wantError: false,
			p: &Project{
				Title: "title",
				Link:  "https://something.com/relative/path/2",
			},
		},
		{
			name:      "valid link #3",
			wantError: false,
			p: &Project{
				Title: "title",
				Link:  "http://10.0.0.0:8443",
			},
		},
		{
			name:      "should pass",
			wantError: false,
			p: &Project{
				Title:       "title",
				Description: "some description",
				Link:        "https://google.com",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.p.Validate()
			if tc.wantError {
				require.Error(tt, err)
				return
			}
			require.NoError(tt, err)
		})
	}
}

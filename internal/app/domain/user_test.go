package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUser_ComparePassword(t *testing.T) {
	tests := []struct {
		name     string
		result   bool
		encode   bool
		password string
		u        *User
	}{
		{
			name:     "empty structure",
			result:   false,
			password: "123",
			u:        &User{},
		},
		{
			name:     "empty password",
			result:   false,
			password: "",
			u:        &User{EncodedPassword: "123", Salt: "321"},
		},
		{
			name:     "empty encoded password",
			result:   false,
			password: "123",
			u:        &User{Salt: "321"},
		},
		{
			name:     "empty salt",
			result:   false,
			password: "123",
			u:        &User{EncodedPassword: "123"},
		},
		{
			name:     "passwords match",
			result:   true,
			password: "123",
			encode:   true,
			u:        &User{EncodedPassword: "123"},
		},
		{
			name:     "passwords dont match",
			result:   false,
			password: "1234",
			encode:   true,
			u:        &User{EncodedPassword: "123"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			if tc.encode {
				err := tc.u.EncodePassword()
				require.NoError(tt, err)
			}

			res := tc.u.ComparePassword(tc.password)
			require.Equal(tt, tc.result, res)
		})
	}
}

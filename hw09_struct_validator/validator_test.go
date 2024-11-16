package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Valid User",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    30,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User ID length",
			in: User{
				ID:     "12345",
				Age:    30,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
			},
		},
		{
			name: "Invalid User Age",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    17,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("must be >= 18")},
			},
		},
		{
			name: "Invalid User Email",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    30,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: fmt.Errorf("regexp does not match invalid-email")},
			},
		},
		{
			name: "Invalid User Role",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    30,
				Email:  "user@example.com",
				Role:   "guest",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Role", Err: fmt.Errorf("must be one of admin,stuff")},
			},
		},
		{
			name: "Invalid User Phone length",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Age:    30,
				Email:  "user@example.com",
				Role:   "admin",
				Phones: []string{"12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "Phones element 0", Err: fmt.Errorf("length must be 11")},
			},
		},
		{
			name: "Valid App",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid App Version length",
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length must be 5")},
			},
		},
		{
			name: "Valid Response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Response Code",
			in: Response{
				Code: 201,
				Body: "Created",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("must be one of 200,404,500")},
			},
		},
		{
			name: "Token without validation",
			in: Token{
				Header:    []byte{1, 2, 3},
				Payload:   []byte{4, 5, 6},
				Signature: []byte{7, 8, 9},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			var valErrs ValidationErrors
			require.ErrorAs(t, err, &valErrs)
			if tt.expectedErr == nil {
				require.Nil(t, valErrs)
				return
			}
			for _, e := range strings.Split(tt.expectedErr.Error(), "\n") {
				require.ErrorContains(t, err, e)
			}
		})
	}
}

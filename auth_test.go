package twitter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterInput_Sanitize(t *testing.T) {
	input := RegisterInput{
		Username: "   bob",
		Email:    "bOb@gmail.com",
		Password: "pass1234",
		ConfirmPassword: "pass1234",
	}

	want := RegisterInput{
		Username: "bob",
		Email:    "bob@gmail.com",
		Password: "pass1234",
		ConfirmPassword: "pass1234",
	}

	input.Sanitize()

	require.Equal(t, want, input)
}

func TestResgisterInput_Validate(t *testing.T) {
	testCases := []struct {
		name  string
		input RegisterInput
		err   error
	}{
		{
			name: "valid",
			input: RegisterInput{
				Username: "bob",
				Email:    "bob@gmail.com",
				Password: "pass1234",
				ConfirmPassword: "pass1234",

			},
			err: nil,
		},
		{
			name: "invalid email",
			input: RegisterInput{
				Username: "bob",
				Email:    "bob.com",
				Password: "pass1234",
				ConfirmPassword: "pass1234",
			},
			err: ErrValidation,
		},

		{
			name: "too short username",
			input: RegisterInput{
				Username: "b",
				Email:    "bob@gmail.com",
				Password: "pass1234",
				ConfirmPassword: "pass1234",
			},
			err: ErrValidation,
		},
		{
			name: "too short password",
			input: RegisterInput{
				Username: "b",
				Email:    "bob@gmail.com",
				Password: "pas",
				ConfirmPassword: "pas",
			},
			err: ErrValidation,
		},

		{
			name: "confirm password doesn't match",
			input: RegisterInput{
				Username: "b",
				Email:    "bob@gmail.com",
				Password: "pas",
			},
			err: ErrValidation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			err := tc.input.Validate()

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

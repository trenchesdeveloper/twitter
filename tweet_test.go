package twitter

import (
	"github.com/stretchr/testify/require"
	"github.com/trenchesdeveloper/tweeter/faker"
	"testing"
)

func TestCreateTweetInput_Sanitize(t *testing.T) {
	tests := []struct {
		name string
		in   *CreateTweetInput
	}{
		{
			name: "sanitizes the input",
			in:   &CreateTweetInput{Body: "Hello, World!      "},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.Sanitize()
			require.Equal(t, "Hello, World!", tt.in.Body)
		})
	}
}

func TestCreateTweetInput_Validate(t *testing.T) {
	tests := []struct {
		name string
		in   CreateTweetInput
		want error
	}{
		{
			name: "valid tweet",
			in:   CreateTweetInput{Body: "Hello, World!"},
			want: nil,
		},
		{
			name: "invalid tweet",
			in:   CreateTweetInput{Body: ""},
			want: ErrValidation,
		},
		{
			name: "tweet too long",
			in:   CreateTweetInput{Body: faker.RandString(300)},
			want: ErrValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Validate()
			require.Equal(t, tt.want, got)
		})
	}
}

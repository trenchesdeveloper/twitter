package test_helper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trenchesdeveloper/tweeter/postgres"
)

func TeardownDB(ctx context.Context, t *testing.T, db *postgres.DB) {
	t.Helper()

	err := db.Truncate(ctx)

	require.NoError(t, err)
}

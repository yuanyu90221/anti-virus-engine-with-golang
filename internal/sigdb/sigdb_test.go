package sigdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/avengine/internal/sigdb"
)

type staticLoader []sigdb.Signature

func (sl staticLoader) Load() ([]sigdb.Signature, error) { return []sigdb.Signature(sl), nil }

func TestLookup_Hit(t *testing.T) {
	loader := staticLoader{
		{SHA256: "abc123", Name: "TestVirus", Category: "ransomware", Severity: "high"},
	}
	db, err := sigdb.NewDB(loader)
	require.NoError(t, err)

	r, ok := db.Lookup("abc123")
	require.True(t, ok, "預期命中，但未找到")

	assert.Equal(t, "TestVirus", r.Name)
	assert.Equal(t, "ransomware", r.Category)
	assert.Equal(t, "high", r.Severity)
}

func TestLookup_Miss(t *testing.T) {
	loader := staticLoader{}
	db, err := sigdb.NewDB(loader)
	require.NoError(t, err)

	_, ok := db.Lookup("deadbeef")
	assert.False(t, ok, "空資料庫應回傳未命中")
}

package scopes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Example struct {
	ID           int64  `json:"id"`
	AdminOnlyStr string `json:"admin_only" scope:"admin"`
	UserOnlyStr  string `json:"user_only,omitempty" scope:"user"`
	BothStr      string `json:"both,omitempty" scope:"user,admin"`
	OmitStr      string `json:"omiter,omitempty" scope:"user,admin"`
	Hidden       string `json:"-"`
	All          string `json:"all,omitempty"`
}

var example = &Example{
	ID:           1,
	AdminOnlyStr: "im an admin",
	UserOnlyStr:  "im a user",
	BothStr:      "im on both",
	Hidden:       "cant see this",
	All:          "should always be known",
}

var example2 = &Example{
	ID:           2,
	AdminOnlyStr: "im an admin",
	UserOnlyStr:  "im a user",
	BothStr:      "im on both",
	Hidden:       "cant see this",
	All:          "should always be known",
}

var examples = []*Example{example, example2}

func TestNew(t *testing.T) {
	out := New("admin", example)
	g := decode(out.JSON())
	assert.Equal(t, int64(1), g.ID)
	assert.Equal(t, "im an admin", g.AdminOnlyStr)
	assert.Empty(t, g.UserOnlyStr)
	assert.Equal(t, "im on both", g.BothStr)
	assert.Empty(t, g.Hidden)

	out = New("user", example)
	g = decode(out.JSON())
	assert.Equal(t, int64(1), g.ID)
	assert.Equal(t, "im a user", g.UserOnlyStr)
	assert.Empty(t, g.AdminOnlyStr)
	assert.Equal(t, "im on both", g.BothStr)
	assert.Empty(t, g.Hidden)

	out = New("user", examples)
	assert.Len(t, decodeMulti(out.JSON()), 2)
}

func decode(val []byte) Example {
	var e Example
	json.Unmarshal(val, &e)
	return e
}

func decodeMulti(val []byte) []Example {
	var e []Example
	json.Unmarshal(val, &e)
	return e
}

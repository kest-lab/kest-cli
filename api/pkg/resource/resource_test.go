package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test model
type testUser struct {
	ID       uint
	Username string
	Email    string
	Avatar   string
	IsAdmin  bool
}

// Test transformer
var testUserTransformer Transformer[*testUser] = func(u *testUser) map[string]any {
	return Filter(map[string]any{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"avatar":   WhenNotEmpty(u.Avatar),
		"is_admin": When(u.IsAdmin, true),
	})
}

func TestTransformer(t *testing.T) {
	user := &testUser{
		ID:       1,
		Username: "john",
		Email:    "john@example.com",
		Avatar:   "",
		IsAdmin:  false,
	}

	result := testUserTransformer.Apply(user)

	assert.Equal(t, uint(1), result["id"])
	assert.Equal(t, "john", result["username"])
	assert.Equal(t, "john@example.com", result["email"])
	assert.NotContains(t, result, "avatar")   // Empty, filtered out
	assert.NotContains(t, result, "is_admin") // False, filtered out
}

func TestTransformerWithValues(t *testing.T) {
	user := &testUser{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Avatar:   "avatar.jpg",
		IsAdmin:  true,
	}

	result := testUserTransformer.Apply(user)

	assert.Equal(t, "avatar.jpg", result["avatar"])
	assert.Equal(t, true, result["is_admin"])
}

func TestTransformerApplyAll(t *testing.T) {
	users := []*testUser{
		{ID: 1, Username: "alice", Email: "alice@example.com"},
		{ID: 2, Username: "bob", Email: "bob@example.com"},
	}

	results := testUserTransformer.ApplyAll(users)

	assert.Len(t, results, 2)
	assert.Equal(t, uint(1), results[0]["id"])
	assert.Equal(t, "alice", results[0]["username"])
	assert.Equal(t, uint(2), results[1]["id"])
	assert.Equal(t, "bob", results[1]["username"])
}

func TestSimpleCollection(t *testing.T) {
	users := []*testUser{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
	}

	collection := NewCollection(users, nil, testUserTransformer)
	results := collection.ToArray()

	assert.Len(t, results, 2)
	assert.Nil(t, collection.GetPaginator())
}

func TestWhen(t *testing.T) {
	assert.Equal(t, "value", When(true, "value"))
	assert.Nil(t, When(false, "value"))
}

func TestWhenNotEmpty(t *testing.T) {
	assert.Equal(t, "test", WhenNotEmpty("test"))
	assert.Nil(t, WhenNotEmpty(""))
}

func TestWhenNotNil(t *testing.T) {
	value := "test"
	assert.Equal(t, "test", WhenNotNil(&value))
	assert.Nil(t, WhenNotNil[string](nil))
}

func TestWhenNotZero(t *testing.T) {
	assert.Equal(t, 42, WhenNotZero(42))
	assert.Nil(t, WhenNotZero(0))
}

func TestUnless(t *testing.T) {
	assert.Nil(t, Unless(true, "value"))
	assert.Equal(t, "value", Unless(false, "value"))
}

func TestFilter(t *testing.T) {
	input := map[string]any{
		"a": 1,
		"b": nil,
		"c": "test",
		"d": nil,
	}

	result := Filter(input)

	assert.Equal(t, 1, result["a"])
	assert.Equal(t, "test", result["c"])
	assert.NotContains(t, result, "b")
	assert.NotContains(t, result, "d")
}

func TestMerge(t *testing.T) {
	result := Merge(
		map[string]any{"a": 1},
		map[string]any{"b": 2},
		nil,
		map[string]any{"c": 3},
	)

	assert.Equal(t, 1, result["a"])
	assert.Equal(t, 2, result["b"])
	assert.Equal(t, 3, result["c"])
}

func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3}
	result := Map(numbers, func(n int) int {
		return n * 2
	})

	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestPluck(t *testing.T) {
	users := []*testUser{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
	}

	ids := Pluck(users, func(u *testUser) uint { return u.ID })
	assert.Equal(t, []uint{1, 2}, ids)

	names := Pluck(users, func(u *testUser) string { return u.Username })
	assert.Equal(t, []string{"alice", "bob"}, names)
}

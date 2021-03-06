package storm

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	db, cleanup := createDB(t, Root("a"))
	defer cleanup()

	n1 := db.From("b", "c")
	node1, ok := n1.(*node)
	assert.True(t, ok)
	assert.Equal(t, db, node1.s)
	assert.NotEqual(t, db.root, n1)
	assert.Equal(t, []string{"a"}, db.root.rootBucket)
	assert.Equal(t, []string{"b", "c"}, node1.rootBucket)
	n2 := n1.From("d", "e")
	node2, ok := n2.(*node)
	assert.True(t, ok)
	assert.Equal(t, []string{"b", "c", "d", "e"}, node2.rootBucket)
}

func TestNodeWithTransaction(t *testing.T) {
	db, cleanup := createDB(t)
	defer cleanup()

	var user User
	db.Bolt.Update(func(tx *bolt.Tx) error {
		dbx := db.WithTransaction(tx)
		err := dbx.Save(&User{ID: 10, Name: "John"})
		assert.NoError(t, err)
		err = dbx.One("ID", 10, &user)
		assert.NoError(t, err)
		assert.Equal(t, "John", user.Name)
		return nil
	})

	err := db.One("ID", 10, &user)
	assert.NoError(t, err)

}

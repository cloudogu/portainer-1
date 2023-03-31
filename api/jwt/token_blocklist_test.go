package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewBlocklistTokenMap(t *testing.T) {
	t.Run("Create new blocklist map", func(t *testing.T) {
		b := NewBlocklistTokenMap(100, time.Hour*1)
		assert.NotNil(t, b)
		assert.Equal(t, 0, len(b.blocklist))
	})
}

func TestNewBlocklistTokenMapTicker(t *testing.T) {
	t.Run("Do automatically remove dead entires with ticker", func(t *testing.T) {
		b := NewBlocklistTokenMap(1, time.Second*2)

		b.Put("test1")
		b.Put("test2")

		assert.Equal(t, 2, len(b.blocklist))
		time.Sleep(time.Second * 3)
		assert.Equal(t, 0, len(b.blocklist))
	})
	t.Run("Do automatically remove dead entires with ticker", func(t *testing.T) {
		b := NewBlocklistTokenMap(5, time.Second*2)

		b.Put("test1")
		b.Put("test2")

		assert.Equal(t, 2, len(b.blocklist))
		time.Sleep(time.Second * 3)
		assert.Equal(t, 2, len(b.blocklist))
	})
}

func TestUpdateList(t *testing.T) {
	t.Run("Do manually remove dead elements", func(t *testing.T) {
		b := NewBlocklistTokenMap(1, time.Hour*1)

		b.Put("test1")
		b.Put("test2")

		assert.Equal(t, 2, len(b.blocklist))
		time.Sleep(time.Second * 2)

		b.UpdateList()
		assert.Equal(t, 0, len(b.blocklist))
	})
	t.Run("Do not manually remove dead elements", func(t *testing.T) {
		b := NewBlocklistTokenMap(5, time.Hour*1)
		b.Put("test1")
		b.Put("test2")

		assert.Equal(t, 2, len(b.blocklist))
		time.Sleep(time.Second * 2)

		b.UpdateList()
		assert.Equal(t, 2, len(b.blocklist))
	})
}

func TestBlocklistTokenMap_Put(t *testing.T) {
	b := NewBlocklistTokenMap(5, time.Hour*1)

	b.Put("test1")
	b.Put("test2")
	assert.Equal(t, 2, len(b.blocklist))

	//Overwrite values
	b.Put("test1")
	b.Put("test2")
	assert.Equal(t, 2, len(b.blocklist))

	//Overwrite values
	b.Put("test3")
	b.Put("test4")
	assert.Equal(t, 4, len(b.blocklist))
}

func TestBlocklistTokenMap_IsBlocked(t *testing.T) {
	b := NewBlocklistTokenMap(5, time.Hour*1)

	b.Put("test1")
	b.Put("test2")
	assert.Equal(t, 2, len(b.blocklist))

	assert.False(t, b.IsBlocked("Test1"))
	assert.False(t, b.IsBlocked("A wda"))
	assert.True(t, b.IsBlocked("test1"))
	assert.True(t, b.IsBlocked("test2"))
}

func TestBlocklistTokenMap_IsExpired(t *testing.T) {
	b := NewBlocklistTokenMap(1, time.Hour*1)

	token := "test1"
	b.Put(token)
	assert.False(t, b.IsExpired(b.blocklist[token]))

	time.Sleep(time.Second * 2)
	assert.True(t, b.IsExpired(b.blocklist[token]))
}

func TestBlocklistTokenMap_Remove(t *testing.T) {
	b := NewBlocklistTokenMap(1, time.Hour*1)

	b.Put("Test1")
	b.Put("Test2")
	assert.NotNil(t, b.blocklist["Test1"])
	assert.NotNil(t, b.blocklist["Test2"])

	b.Remove("Test1")
	assert.Nil(t, b.blocklist["Test1"])
	assert.NotNil(t, b.blocklist["Test2"])

	b.Remove("Test2")
	assert.Nil(t, b.blocklist["Test1"])
	assert.Nil(t, b.blocklist["Test2"])
}

package wshub

import (
	wslib "foodtruck/pkg/websocket/server"
	"sync"
	"testing"

	assrt "github.com/blend/go-sdk/assert"
)

func TestWSServer(t *testing.T) {
	var assert = assrt.New(t)

	hub := New()
	conn1 := wslib.New()
	conn2 := wslib.New()

	hub.Add("1", conn1)
	hub.Add("2", conn2)

	assert.Equal(2, len(hub.Conns()), "number of connections should increase after being added to hub")

	hub.Remove("1")
	assert.Equal(1, len(hub.Conns()), "number of connections should decrease after being removed from hub")
}

func TestWSServerMutex(t *testing.T) {
	var assert = assrt.New(t)
	var wg sync.WaitGroup

	hub := New()
	conn1 := wslib.New()
	conn2 := wslib.New()

	wg.Add(2)

	go func() {
		hub.Add("1", conn1)
		wg.Done()
	}()
	go func() {
		hub.Add("2", conn2)
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(2, len(hub.Conns()), "number of connections should increase after being added to hub")

	wg.Add(2)
	go func() {
		hub.Remove("1")
		wg.Done()
	}()
	go func() {
		hub.Remove("2")
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(0, len(hub.Conns()), "number of connections should decrease after being removed from hub")

}

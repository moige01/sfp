package pubsub_test

import (
	"fmt"
	"testing"

	"sugud0r.dev/sfp/internal/pubsub"
)

func Test(t *testing.T) {
	ps := pubsub.NewPubsub()
	ch1 := ps.Subscribe(1)
	ch2 := ps.Subscribe(2)
	ch3 := ps.Subscribe(2)

	listener := func(name string, ch <-chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}

		fmt.Printf("[%s] done\n", name)
	}

	go listener("1", ch1)
	go listener("2", ch2)
	go listener("3", ch3)

	pub := func(topic int, msg string) {
		fmt.Printf("Publishing @%v: %v\n", topic, msg)

		ps.Publish(topic, msg)
	}

	pub(1, "tablets")
	pub(3, "vitamins")
	pub(2, "beaches")
	pub(2, "hiking")

	ps.Close()
}

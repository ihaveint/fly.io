package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type GlobalID struct {
	counter int32

	mu sync.Mutex
}

func main() {
	n := maelstrom.NewNode()

	current_id := GlobalID{
		counter: 0,
	}

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		current_id.mu.Lock()
		defer current_id.mu.Unlock()

		unique_id := fmt.Sprintf("%s-%d", n.ID(), current_id.counter)

		body["type"] = "generate_ok"
		body["id"] = unique_id
		current_id.counter += 1

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}

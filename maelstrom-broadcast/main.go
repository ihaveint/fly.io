package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	var mu sync.Mutex
	msgs := make(map[int]struct{})

	// handle topology: reply topology_ok
	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{"type": "topology_ok"})
	})

	// broadcast: add message to set, reply {"type":"broadcast_ok"}
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		message := int(body["message"].(float64))  // Extract as number
		mu.Lock()
		msgs[message] = struct{}{}  // Store the integer
		mu.Unlock()
		return n.Reply(msg, map[string]any{"type": "broadcast_ok"})
	})

	// read: return list of messages {"type":"read_ok","messages":[...]} 
	n.Handle("read", func(msg maelstrom.Message) error {
		mu.Lock()
		arr := make([]int, 0, len(msgs))
		for k := range msgs {
			arr = append(arr, k)
		}
		mu.Unlock()
		return n.Reply(msg, map[string]any{"type": "read_ok", "messages": arr})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

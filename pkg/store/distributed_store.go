package store

import (
	"encoding/json"

	"github.com/TomStuart92/asfalis/pkg/logger"
)

var log = logger.NewStdoutLogger("store")

// DistributedStore is a key-value store which is designed to propose changes
// and read commits from a pair of channels. These channels are usually backed
// by an implementation of the Raft algorithm.
type DistributedStore struct {
	proposeC chan<- string
	commitC  <-chan string
	store    *LocalStore
}

// keyValue is an internal representation of a key-value pair used to send such
// pairs into and out of the stores channels.
type keyValue struct {
	Key   string
	Value string
}

func (kv *keyValue) Encode() ([]byte, error) {
	return json.Marshal(kv)
}

// NewDistributedStore creates a new instance of a DistributedStore
func NewDistributedStore(proposeC chan<- string, commitC <-chan string) *DistributedStore {
	s := &DistributedStore{
		proposeC,
		commitC,
		NewLocalStore(),
	}
	go s.readCommits()
	return s
}

// Lookup delegates a request through to the underlying store instance
func (s *DistributedStore) Lookup(key string) (string, bool) {
	log.Infof("Looking up key %s", key)
	return s.store.Get(key)
}

// Propose a change to the store through the proposeC channel
func (s *DistributedStore) Propose(key string, value string) error {
	kv := keyValue{
		Key:   key,
		Value: value,
	}
	bytes, err := kv.Encode()
	if err != nil {
		return err
	}

	s.proposeC <- string(bytes)
	log.Infof("Proposed setting %s => %s", key, value)
	return nil
}

// readCommits loops through commits in the commitC channel
//  and applies them as appropriate
func (s *DistributedStore) readCommits() {
	for data := range s.commitC {
		if data == "" {
			continue
		}

		log.Infof("Received data from commit channel: %v", data)

		var kv keyValue
		if err := json.Unmarshal([]byte(data), &kv); err != nil {
			log.Fatalf("Failed to decode message (%v)", err)
		}

		if kv.Value == "" {
			log.Infof("Key %s Deleted", kv.Key)
			s.store.Delete(kv.Key)
		} else {
			log.Infof("Setting %s => %s", kv.Key, kv.Value)
			s.store.Set(kv.Key, kv.Value)
		}
	}
}

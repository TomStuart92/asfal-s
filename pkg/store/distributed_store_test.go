package store

import (
	"testing"
	"time"
)

func replayProposals(proposeC <-chan string, commitC chan<- string) {
	for data := range proposeC {
		commitC <- data
	}
}

func loadDummyRecord(commitC chan string) {
	var s string
	commitC <- s
}

func TestNoRecordOnLookup(t *testing.T) {
	proposeC := make(chan string, 1)
	commitC := make(chan string, 1)
	loadDummyRecord(commitC)
	d := NewDistributedStore(proposeC, commitC)
	_, ok := d.Lookup("key")
	if ok {
		t.Error("Retrieved Invalid Data")
	}
}

func TestPropose(t *testing.T) {
	proposeC := make(chan string, 1)
	commitC := make(chan string, 1)
	loadDummyRecord(commitC)
	d := NewDistributedStore(proposeC, commitC)
	d.Propose("key", "value")
	select {
	case <-proposeC:
		return
	default:
		t.Error("No Proposal Recieved")
	}
}

func TestCommit(t *testing.T) {
	proposeC := make(chan string, 1)
	commitC := make(chan string, 1)
	loadDummyRecord(commitC)
	d := NewDistributedStore(proposeC, commitC)
	// replay proposals directly into commits
	go replayProposals(proposeC, commitC)

	d.Propose("key", "value")

	// allow store time to set key. (eventually consistent)
	time.Sleep(2 * time.Second)
	value, ok := d.Lookup("key")
	if !ok {
		t.Error("Could Not Retrieve Key From Store")
	}
	if value != "value" {
		t.Error("Incorrect Value Retrieved")
	}
}

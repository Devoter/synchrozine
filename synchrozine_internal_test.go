package synchrozine

import (
	"testing"
)

func TestNew(t *testing.T) {
	synchro := New()

	if synchro.message == nil {
		t.Fatalf("message channel should be initialized\n")
	}
	if synchro.receiverCh == nil {
		t.Fatalf("receiver channel should be initialized\n")
	}
	if synchro.counterCh == nil {
		t.Fatalf("counter channel should be initialized\n")
	}
	if synchro.startCounterCh == nil {
		t.Fatalf("start counter channel should be initialized\n")
	}
}

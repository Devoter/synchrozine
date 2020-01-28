package synchrozine

import (
	"testing"
)

func TestNew(t *testing.T) {
	synchro := New()

	if synchro.message == nil {
		t.Fatalf("message channel should be initialized\n")
	}
	if synchro.receivers == nil {
		t.Fatalf("receivers list should be initialized\n")
	}
	if synchro.wg == nil {
		t.Fatalf("wait group should be initialized\n")
	}
	if synchro.startWG == nil {
		t.Fatalf("start wait group should be initialized\n")
	}
}

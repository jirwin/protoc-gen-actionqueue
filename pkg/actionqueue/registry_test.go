package actionqueue

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"google.golang.org/protobuf/proto"
)

// resetRegistry replaces the global registry with a fresh empty map.
func resetRegistry() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]*Definition)
}

func TestRegisterAndGet(t *testing.T) {
	resetRegistry()

	def := &Definition{
		Name:   "test-queue",
		Signal: "test-signal",
	}
	Register(def)

	got := Get("test-queue")
	if got == nil {
		t.Fatal("expected definition, got nil")
	}
	if got.Name != "test-queue" {
		t.Errorf("Name = %q, want %q", got.Name, "test-queue")
	}
	if got.Signal != "test-signal" {
		t.Errorf("Signal = %q, want %q", got.Signal, "test-signal")
	}
}

func TestRegisterPanicsOnEmptyName(t *testing.T) {
	resetRegistry()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty name")
		}
	}()
	Register(&Definition{Name: ""})
}

func TestRegisterPanicsDuplicateName(t *testing.T) {
	resetRegistry()

	Register(&Definition{Name: "dup"})

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for duplicate name")
		}
	}()
	Register(&Definition{Name: "dup"})
}

func TestGetReturnsNilForUnknown(t *testing.T) {
	resetRegistry()

	if got := Get("nonexistent"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestAllReturnsAllDefinitions(t *testing.T) {
	resetRegistry()

	names := []string{"alpha", "bravo", "charlie"}
	for _, n := range names {
		Register(&Definition{Name: n})
	}

	all := All()
	if len(all) != 3 {
		t.Fatalf("All() returned %d definitions, want 3", len(all))
	}

	gotNames := make([]string, len(all))
	for i, d := range all {
		gotNames[i] = d.Name
	}
	sort.Strings(gotNames)

	for i, want := range names {
		if gotNames[i] != want {
			t.Errorf("gotNames[%d] = %q, want %q", i, gotNames[i], want)
		}
	}
}

func TestAllReturnsEmptySliceWhenEmpty(t *testing.T) {
	resetRegistry()

	all := All()
	if all == nil {
		t.Fatal("All() returned nil, want empty slice")
	}
	if len(all) != 0 {
		t.Errorf("All() returned %d definitions, want 0", len(all))
	}
}

func TestConcurrentRegisterAndGet(t *testing.T) {
	resetRegistry()

	var wg sync.WaitGroup

	// 50 writers, each registering a unique name
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			Register(&Definition{Name: fmt.Sprintf("concurrent-%d", n)})
		}(i)
	}

	// 50 readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = Get(fmt.Sprintf("concurrent-%d", n))
			_ = All()
		}(i)
	}

	wg.Wait()

	all := All()
	if len(all) != 50 {
		t.Errorf("All() returned %d definitions, want 50", len(all))
	}
}

func TestRegisterStoresFunctionFields(t *testing.T) {
	resetRegistry()

	wfIDFunc := func(msg proto.Message) string { return "wf-id" }
	tenantIDFunc := func(msg proto.Message) string { return "tenant-id" }
	entityIDsFunc := func(msg proto.Message) []string { return []string{"e1", "e2"} }
	wfIDFromArgs := func(tenantID string, entityIDs []string) string { return "from-args" }

	def := &Definition{
		Name:               "func-test",
		Signal:             "sig",
		WorkflowIDFunc:     wfIDFunc,
		TenantIDFunc:       tenantIDFunc,
		EntityIDsFunc:      entityIDsFunc,
		WorkflowIDFromArgs: wfIDFromArgs,
	}
	Register(def)

	got := Get("func-test")
	if got == nil {
		t.Fatal("expected definition, got nil")
	}
	if got.WorkflowIDFunc(nil) != "wf-id" {
		t.Error("WorkflowIDFunc not preserved")
	}
	if got.TenantIDFunc(nil) != "tenant-id" {
		t.Error("TenantIDFunc not preserved")
	}
	if ids := got.EntityIDsFunc(nil); len(ids) != 2 || ids[0] != "e1" || ids[1] != "e2" {
		t.Error("EntityIDsFunc not preserved")
	}
	if got.WorkflowIDFromArgs("t", []string{"e"}) != "from-args" {
		t.Error("WorkflowIDFromArgs not preserved")
	}
}

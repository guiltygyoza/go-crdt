package crdt

import (
	// "bytes"
	// "encoding/json"
	"testing"
)

func TestInfPhaseSetAdd(t *testing.T) {
	ipset := NewIPSet()
	elemToAdd := "blah"
	elemToSkip := "hey"

	if ipset.Contains(elemToAdd) {
		t.Errorf("set should not contain %q; not yet added", elemToAdd)
	}

	// add
	ipset.Add(elemToAdd)
	if !ipset.Contains(elemToAdd) {
		t.Errorf("set should contain %q; have been added", elemToAdd)
	}
	if ipset.Contains(elemToSkip) {
		t.Errorf("set should not contain %q", elemToSkip)
	}

	// re-add, verifying add's idempotency
	ipset.Add(elemToAdd)
	if !ipset.Contains(elemToAdd) {
		t.Errorf("set should contain %q; have been added", elemToAdd)
	}
}

func TestInfPhaseSetAddRemoveSequence(t *testing.T) {
	for ii, tt := range []struct {
		elem           		interface{}
		fnAddRemoveSequence func(*IPSet, interface{})
		result         		bool
		description			string
	}{
		{
			"blah",
			func(s *IPSet, obj interface{}) {
				s.Add(obj)
				s.Remove(obj)
			},
			false,
			"Add, then remove. Outcome: element should not exist",
		},
		{
			"blah",
			func(s *IPSet, obj interface{}) {
				s.Add(obj)
				s.Add(obj)
				s.Remove(obj)
			},
			false,
			"Add twice, then remove. Outcome: element should not exist",
		},
		{
			"blah",
			func(s *IPSet, obj interface{}) {
				s.Add(obj)
				s.Add(obj)
				s.Remove(obj)
				s.Remove(obj)
				s.Add(obj)
				s.Add(obj)
				s.Add(obj)
			},
			true,
			"Add twice, then remove twice, then add three times. Outcome: element should exist",
		},
		{
			"blah",
			func(s *IPSet, obj interface{}) {
				s.Add(obj)
				s.Remove(obj)
				s.Add(obj)
				s.Remove(obj)
				s.Add(obj)
				s.Remove(obj)
				s.Add(obj)
				s.Remove(obj)
				s.Add(obj)
				s.Remove(obj)
				s.Add(obj)
				s.Add(obj)
			},
			true,
			"Add-remove 5 times, then add two more times. Outcome: element should exist",
		},
	} {
		ipset := NewIPSet()

		if ipset.Contains(tt.elem) {
			t.Errorf("set should not contain elem %q", tt.elem)
		}

		tt.fnAddRemoveSequence(ipset, tt.elem)

		if ipset.Contains(tt.elem) != tt.result {
			t.Errorf("AddRemoveSequence test #%q failed; test description: %s", ii, tt.description)
		}
	}
}

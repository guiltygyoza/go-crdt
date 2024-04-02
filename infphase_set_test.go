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

func TestInfPhaseSetMerge(t *testing.T) {
	for _, tt := range []struct {
		fnMutateSet1 		func(*IPSet)
		fnMutateSet2 		func(*IPSet)
		mergeResult         map[interface{}]*GCounter
		description			string
	}{
		{
			func(s *IPSet) {
				s.Add(1)
				s.Add(1)
				s.Remove(3)
				s.Add(3)
				s.Remove(3)
				s.Add(5)
				s.Remove(5)
				s.Add(5)
			},
			func(s *IPSet) {
				s.Add(1)
				s.Remove(1)
				s.Add(3)
			},
			map[interface{}]*GCounter {
				1:NewGCounterInit(2),
				3:NewGCounterInit(2),
				5:NewGCounterInit(3),
			},
			"Merging {1:1, 3:2, 5:3} and {1:2, 3:1} should yield {1:2, 3:2, 5:3}",
		},
		{
			func(s *IPSet) {
				s.Add('a')
				s.Add('b')
			},
			func(s *IPSet) {
				s.Add(5)
				s.Remove(7)
				s.Add(7)
			},
			map[interface{}]*GCounter {
				'a':NewGCounterInit(1),
				'b':NewGCounterInit(1),
				5:NewGCounterInit(1),
				7:NewGCounterInit(1),
			},
			"Merging {'a':1, 'b':1} and {5:1, 7:1} should yield {'a':1, 'b':1, 5:1, 7:1}",
		},
	} {
		ipset1 := NewIPSet()
		tt.fnMutateSet1(ipset1)

		ipset2 := NewIPSet()
		tt.fnMutateSet2(ipset2)

		ipset1.Merge(ipset2)

		for elem_, counter_ := range tt.mergeResult {
			counter, ok := ipset1.dict[elem_]
			if !ok || counter.Count()!=counter_.Count() {
				t.Errorf("Set's internal dictionary should contain the following entry: %q:%q", elem_, counter_.Count())
			}
		}
	}
}
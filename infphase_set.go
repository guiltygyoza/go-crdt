package crdt

// import "encoding/json"

type gSetGCounter map[interface{}]*GCounter

var (
	// IPSet should implement the set interface.
	_ Set = &IPSet{}
)

//
// State
// - an element exists in the infinite-phase set if dict[element] exists and is an odd number
// - otherwise the element doesn't exist in the infinite-phase set
//
type IPSet struct {
	dict gSetGCounter
}

//
// New
//
func NewIPSet() *IPSet {
	return &IPSet{
		dict: gSetGCounter{},
	}
}

//
// Operation: Add
//
func (s *IPSet) Add(elem interface{}) {
	// if elem is not in s.dict => s.dict[elem] = 1
	// else if s.dict[elem] is even => s.dict[elem] += 1
	_, ok := s.dict[elem]
	if !ok {
		s.dict[elem] = NewGCounterInit(1)
	} else if s.dict[elem].Count() % 2 == 0 {
		s.dict[elem].Inc()
	}
}

//
// Operation: Remove
//
func (s *IPSet) Remove(elem interface{}) {
	// if elem is in s.dict && s.dict[elem] is odd => s.dict[elem] += 1
	counter, ok := s.dict[elem]
	if ok && counter.Count() % 2 == 1 {
		s.dict[elem].Inc()
	}
}

//
// Check element existence
//
func (s *IPSet) Contains(elem interface{}) bool {
	counter, ok := s.dict[elem]
	if !ok {
		return false
	}
	return counter.Count() % 2 == 1
}

//
// Get size of set
//
func (s *IPSet) Len() int {
	len := 0
	for _, counter := range s.dict {
		if counter.Count() % 2 == 1 {
			len += 1
		}
	}
	return len
}

//
// Materialize the set
//
func (s *IPSet) Elems() []interface{} {
	elems := make([]interface{}, 0)
	for elem, counter := range s.dict {
		if counter.Count() % 2 == 1 {
			elems = append(elems, elem)
		}
	}
	return elems
}

//
// Merge with another replica
//
func (s *IPSet) Merge(s_ *IPSet) {
	for elem_, counter_ := range s_.dict {
		// if elem_ is not in local replica => add
		// otherwise, set local counter for elem_ to max(counter, counter_)
		counter, ok := s.dict[elem_]
		if !ok {
			s.dict[elem_] = NewGCounterInit(1)
		} else {
			s.dict[elem_] = NewGCounterInit(max(counter.Count(), counter_.Count()))
		}
	}
}

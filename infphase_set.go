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
		newGCounter := NewGCounter()
		newGCounter.Inc()
		s.dict[elem] = newGCounter
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
// Get all existing elements in the set
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

// type gsetJSON struct {
// 	T string        `json:"type"`
// 	E []interface{} `json:"e"`
// }

// // MarshalJSON will be used to generate a serialized output
// // of a given GSet.
// func (g *GSet) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(&gsetJSON{
// 		T: "g-set",
// 		E: g.Elems(),
// 	})
// }

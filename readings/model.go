package readings

import (
	"sort"

	"github.com/ViBiOh/auth/auth"
)

type tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	user *auth.User
}

type tagsSorter struct {
	arr []*tag
	by  func(p1, p2 *tag) bool
}

func (s *tagsSorter) Len() int {
	return len(s.arr)
}

func (s *tagsSorter) Swap(i, j int) {
	s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
}

func (s *tagsSorter) Less(i, j int) bool {
	return s.by(s.arr[i], s.arr[j])
}

type sortTagBy func(p1, p2 *tag) bool

func (by sortTagBy) Sort(arr []*tag) {
	sort.Sort(&tagsSorter{
		arr: arr,
		by:  by,
	})
}

func sortTagByID(o1, o2 *tag) bool {
	return o1.ID < o2.ID
}

type reading struct {
	ID     int64  `json:"id"`
	URL    string `json:"url"`
	Public bool   `json:"public"`
	Read   bool   `json:"read"`
	Tags   []*tag `json:"tags"`
	user   *auth.User
}

type readingsSorter struct {
	arr []*reading
	by  func(p1, p2 *reading) bool
}

func (s *readingsSorter) Len() int {
	return len(s.arr)
}

func (s *readingsSorter) Swap(i, j int) {
	s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
}

func (s *readingsSorter) Less(i, j int) bool {
	return s.by(s.arr[i], s.arr[j])
}

type sortReadingBy func(p1, p2 *reading) bool

func (by sortReadingBy) Sort(arr []*reading) {
	sort.Sort(&readingsSorter{
		arr: arr,
		by:  by,
	})
}

func sortReadingByID(o1, o2 *reading) bool {
	return o1.ID < o2.ID
}

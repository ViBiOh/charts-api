package readings

type tag struct {
	id   int64
	Name string `json:"name"`
}

type reading struct {
	id     int64
	URL    string `json:"url"`
	Public bool   `json:"public"`
	Read   bool   `json:"read"`
	Tags   []*tag `json:"tags"`
}

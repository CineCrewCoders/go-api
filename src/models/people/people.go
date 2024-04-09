package people

type People struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Birthday time.Time `json:"birthday"`
	Photo string `json:"photo"`
}

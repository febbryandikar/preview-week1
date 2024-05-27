package entity

type Branch struct {
	ID       int    `json:"branch_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

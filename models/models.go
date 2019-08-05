package models

// Author JSON model
type Author struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

// Book JSON model
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

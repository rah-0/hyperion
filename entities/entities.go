package entities

/*
	Define your entities (structs) in this file.
*/

type Sample struct {
	Name    string `json:"-"`
	Surname string `json:"-"`
	//FullName string
}

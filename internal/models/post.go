package models

import "html/template"

type Post struct {
	Status         string `yaml:"status"`
	Title          string `yaml:"title"`
	Date           string `yaml:"date"`
	Summary        string `yaml:"summary"`
	Body           template.HTML
	File           string
	Comments       []Comment
	EnableComments bool `yaml:"enableComments"`
	Ert            string
}

package main

import "github.com/Prateek2593/snippetbox/internal/models"

// define a templateData type to act as the holding structure for any dynamic data that we want to pass to our HTML template. at the moment it only contains one field, but will add more
type templateData struct {
	Snippet *models.Snippet
}

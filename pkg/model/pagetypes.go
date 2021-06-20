package model

type ProjectDetailsPage struct {
	Project Project
	Context string
}

type PageData struct {
	Projects []Project
	Context string
}

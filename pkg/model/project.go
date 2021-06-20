package model

import "html/template"

type Project struct {
	FolderPath          string
	FolderPathLowerCase string
	Title               string
	Technologies        []string
	Summary             template.HTML
	ShortSummary        string
	LearningOutcomes    template.HTML
	ProjectOutcomes     template.HTML
	Images              []string
	HighlightedProject  bool
	PlaceholderImage string
}
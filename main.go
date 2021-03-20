package main

import (
	"encoding/csv"
	"html/template"
	"io/ioutil"
	_ "io/ioutil"
	"os"
	"strings"
)

const (
	PROJECT_FOLDER = "projects"
	TECHFILENAME    = "technologies.csv"
	SUMMARYFILENAME = "summary.txt"
	LEARNINGOUTCOMES = "learningoutcomes.txt"
	PROJECTOUTCOMES = "projectoutcomes.txt"
	IMAGEFOLDERNAME = "images"

	OUTPUTDIR = "build"
)


type project struct{
	Title string
	Technologies []string
	Summary string
	LearningOutcomes string
	ProjectOutcomes string
	Images []string
}

type projectDetailsPage struct{
	Project project
}

func main() {
	projectDirs, err := ioutil.ReadDir(PROJECT_FOLDER)

	if err != nil {
		panic("Cannot load files")
	}

	var projects []project = loadProjects(projectDirs)
	t, err := template.ParseFiles("templates/basetemplate.html")

	if err != nil {
		panic("Cannot parse template" )
	}

	error := t.Execute(os.Stdout, projectDetailsPage{
		Project: projects[0],
	})

	if error != nil {
		panic("Cannot execute template"+ error.Error())
	}
	_ = projects
}

func loadProjects(projectDirs []os.FileInfo) []project{
	var projects []project

	for _, f := range projectDirs {
		if !f.IsDir() {
			continue
		}

		projects = append(projects, loadProject(f))
	}
	return projects
}

func loadProject(f os.FileInfo) project {
	return project{
		Title:        strings.Replace(f.Name(), "_", " ", -1),
		Technologies: loadProjectTechFile(f),
		Summary:      loadProjectSummaryFile(f),
		Images:       loadProjectImagePaths(f),
	}
}

func loadProjectImagePaths(f os.FileInfo) []string {
	imageFiles, err := ioutil.ReadDir(PROJECT_FOLDER + "/" + f.Name() + "/" + IMAGEFOLDERNAME)
	if err != nil {
		return []string{}
	}
	var imagePaths []string

	for _, f := range imageFiles {
		if f.IsDir() {
			continue
		}
		imagePaths = append(imagePaths, f.Name())
	}
	return imagePaths
}

func loadProjectSummaryFile(f os.FileInfo) string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + SUMMARYFILENAME
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	val, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}
	return string(val)
}

func loadProjectTechFile(f os.FileInfo) []string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + TECHFILENAME
	file, err := os.Open(path)

	if err != nil {
		return []string{}
	}

	reader := csv.NewReader(file)
	result, err := reader.ReadAll()
	if err != nil {
		return []string{}
	}

	var flat []string
	for _, v := range result {
		flat = append(flat, v...)
	}

	return flat
}


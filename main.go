package main

import (
	"encoding/csv"
	"fmt"
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
	title string
	technologies []string
	summary string
	learningOutcomes string
	projectOutcomes string
	images []string
}

func main() {
	projectDirs, err := ioutil.ReadDir(PROJECT_FOLDER)

	if err != nil {
		panic("Cannot load files")
	}

	var projects []project = loadProjects(projectDirs)
	fmt.Print(projects)
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
		title:        strings.Replace(f.Name(), "_", " ", -1),
		technologies: loadProjectTechFile(f),
		summary:      loadProjectSummaryFile(f),
		images:       loadProjectImagePaths(f),
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


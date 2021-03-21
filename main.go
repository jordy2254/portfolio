package main

import (
	"encoding/csv"
	"html/template"
	"io"
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
	FolderPath string
	FolderPathLowerCase string
	Title string
	Technologies []string
	Summary string
	LearningOutcomes string
	ProjectOutcomes string
	Images []string
	HighlightedProject bool
}

type projectDetailsPage struct{
	Project project
}

type indexPage struct{
	FeaturedProjects []project
}

func main() {
	projectDirs, err := ioutil.ReadDir(PROJECT_FOLDER)
	//https://www.calhoun.io/intro-to-templates-p2-actions/
	if err != nil {
		panic("Cannot load files")
	}

	var projects []project = loadProjects(projectDirs)
	t, err := template.ParseFiles("templates/ProjectViewTemplate.gohtml", "templates/index.gohtml", "templates/nav.gohtml")

	if _, err := os.Stat(OUTPUTDIR); !os.IsNotExist(err){
		error := os.RemoveAll(OUTPUTDIR + "/*")
		if error != nil {
			panic(error)
		}
	}

	os.Mkdir(OUTPUTDIR, os.ModePerm)
	os.Mkdir(OUTPUTDIR + "/" + PROJECT_FOLDER, os.ModePerm)
	os.Mkdir(OUTPUTDIR + "/" + PROJECT_FOLDER + "/images", os.ModePerm)

	//Generate the html pages for the projects
	for _, v := range projects {
		fileOuput, err := os.Create(OUTPUTDIR + "/" + PROJECT_FOLDER + "/" + v.FolderPathLowerCase + ".html")

		IMAGEPATH := OUTPUTDIR + "/" + PROJECT_FOLDER + "/images/" + v.FolderPathLowerCase

		os.Mkdir(IMAGEPATH, os.ModePerm)
		//copy image resources
		for _, image := range(v.Images){
			reader, err := os.Open(PROJECT_FOLDER + "/" + v.FolderPath + "/images/" + image)
			if err != nil {
				panic(err)
			}
			writer, err := os.Create(IMAGEPATH + "/" + image)
			if err != nil {
				panic(err)
			}
			io.Copy(writer, reader)
		}
		defer fileOuput.Close()

		if err != nil {
			panic(err)
		}
		error := t.ExecuteTemplate(fileOuput, "ProjectViewTemplate.gohtml", projectDetailsPage{
			Project: v,
		})
		if error != nil {
			panic(error)
		}
	}

	fileOuput, err := os.Create(OUTPUTDIR + "/index.html")

	var featuredProjects []project
	for _, v := range projects{
		if(!v.HighlightedProject){
			continue
		}
		featuredProjects = append(featuredProjects, v)
	}

	error := t.ExecuteTemplate(fileOuput, "index.gohtml", indexPage{
		FeaturedProjects: featuredProjects,
	})
	if error != nil {
		panic(error)
	}
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
		FolderPath:          f.Name(),
		FolderPathLowerCase: strings.ToLower(f.Name()),
		Title:               strings.Replace(f.Name(), "_", " ", -1),
		Technologies:        loadProjectTechFile(f),
		Summary:             loadProjectSummaryFile(f),
		LearningOutcomes:    "",
		ProjectOutcomes:     "",
		Images:              loadProjectImagePaths(f),
		HighlightedProject:  isHighlightedProject(f),
	}
}

func isHighlightedProject(f os.FileInfo) bool {
	_, err := os.Stat(PROJECT_FOLDER + "/" + f.Name() + "/highlighted.txt")
	return !os.IsNotExist(err)
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


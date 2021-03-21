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
	SHORTSUMMARYFILENAME = "shortsummary.txt"
	LEARNINGOUTCOMES = "learningoutcomes.txt"
	PROJECTOUTCOMES = "projectoutcomes.txt"
	IMAGEFOLDERNAME = "images"
	OUTPUTDIR = "build"
	RESOURCES = "resources"
	STATICPAGES = "staticpages"
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

type projectsPage struct{
	Projects []project
}

func main() {
	projectDirs, err := ioutil.ReadDir(PROJECT_FOLDER)

	singlePages := make(map[string]string)
	singlePages["index.gohtml"] = "index.html"
	singlePages["projects.gohtml"] = "projects.html"


	if err != nil {
		panic("Cannot load files")
	}

	var projects = loadProjects(projectDirs)
	t, err := template.ParseFiles("templates/ProjectViewTemplate.gohtml",
		"templates/index.gohtml",
		"templates/projects.gohtml",
		"templates/nav.gohtml")

	if _, err := os.Stat(OUTPUTDIR); !os.IsNotExist(err){
		error := os.RemoveAll(OUTPUTDIR)
		if error != nil {
			panic(error)
		}
	}

	os.Mkdir(OUTPUTDIR, os.ModePerm)
	os.Mkdir(OUTPUTDIR + "/" + PROJECT_FOLDER, os.ModePerm)
	os.Mkdir(OUTPUTDIR + "/" + PROJECT_FOLDER + "/images", os.ModePerm)
	os.Mkdir(OUTPUTDIR + "/" + RESOURCES, os.ModePerm)

	copyDirectoryRecursive(RESOURCES, OUTPUTDIR + "/" + RESOURCES)
	copyDirectoryRecursive(STATICPAGES, OUTPUTDIR)

	generateProjectPages(projects, t)
	generateSinglePages(singlePages, projects, t)
}

func generateSinglePages(pages map[string]string, projects []project, t *template.Template) {
	for k, v := range pages{
		fileOuput, _ := os.Create(OUTPUTDIR + "/" + v)
		error := t.ExecuteTemplate(fileOuput, k, projectsPage{
			Projects: projects,
		})
		if error != nil {
			panic(error)
		}
	}


}

func generateProjectPages(projects []project, t *template.Template) {
	for _, v := range projects {
		fileOuput, err := os.Create(OUTPUTDIR + "/" + PROJECT_FOLDER + "/" + v.FolderPathLowerCase + ".html")

		imagePath := OUTPUTDIR + "/" + PROJECT_FOLDER + "/images/" + v.FolderPathLowerCase

		os.Mkdir(imagePath, os.ModePerm)
		//copy image resources
		for _, image := range v.Images {
			copyFile(PROJECT_FOLDER+"/"+v.FolderPath+"/images/"+image, imagePath+"/"+image)
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
}

func copyDirectoryRecursive(src, dst string){
	files, err := ioutil.ReadDir(src)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(dst); os.IsNotExist(err){
		os.Mkdir(dst, os.ModePerm)
	}
	if err != nil {
		panic(err)
	}
	for _, v := range files{
		if(v.IsDir()){
			copyDirectoryRecursive(src + "/" + v.Name(), dst + "/" + v.Name())
			continue
		}
 		copyFile(src + "/" + v.Name(), dst + "/" + v.Name())
	}
}

func copyFile(src, dst string){
	reader, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	writer, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	io.Copy(writer, reader)
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
		LearningOutcomes:    loadLearningOutcomesFile(f),
		ProjectOutcomes:     loadProjectOutcomesFile(f),
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
func loadProjectOutcomesFile(f os.FileInfo) string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + PROJECTOUTCOMES
	return loadTextFile(path)
}
func loadLearningOutcomesFile(f os.FileInfo) string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + LEARNINGOUTCOMES
	return loadTextFile(path)
}

func loadProjectSummaryFile(f os.FileInfo) string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + SUMMARYFILENAME
	return loadTextFile(path)
}

func loadTextFile(src string) string{
	file, err := os.Open(src)
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

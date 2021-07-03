package main

import (
	"encoding/csv"
	"flag"
	"github.com/jordy2254/portfolio/pkg/model"
	"github.com/op/go-logging"
	"html/template"
	"io"
	"io/ioutil"
	_ "io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	PROJECT_FOLDER       = "projects"
	TECHFILENAME         = "technologies.csv"
	SUMMARYFILENAME      = "summary.txt"
	SHORTSUMMARYFILENAME = "shortsummary.txt"
	LEARNINGOUTCOMES     = "learningoutcomes.txt"
	PROJECTOUTCOMES      = "projectoutcomes.txt"
	IMAGEFOLDERNAME      = "images"
	RESOURCES            = "resources"
	STATICPAGES          = "staticpages"
)

var(
	context = "/"
	outputPath = "build"

	logger  = logging.MustGetLogger("example")
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
	)
)

func createTemplateForPage(page string) (*template.Template, error){
	return template.ParseFiles(page, "templates/masterTemplate.gohtml", "templates/nav.gohtml")
}

func main() {
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0),format))

	bs := flag.Bool("build-site", false, "Builds the entire site given the params")
	flag.StringVar(&context, "context", context, "Sets the context for the website route when building")
	flag.Parse()

	if *bs {
		logger.Info("Building Site...")
		buildSite()
		logger.Info("Site built successfully")
	}
}

func buildSite(){

	projectDirs, err := ioutil.ReadDir(PROJECT_FOLDER)

	if err != nil {
		logger.Error("Failed to load project Folder build failed")
		os.Exit(1)
	}

	projects := loadProjects(projectDirs)

	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		error := os.RemoveAll(outputPath)
		if error != nil {
			logger.Error("Failed to delete prior build, build failed")
			os.Exit(1)
		}
	}

	os.Mkdir(outputPath, os.ModePerm)
	os.Mkdir(outputPath+"/"+PROJECT_FOLDER, os.ModePerm)
	os.Mkdir(outputPath+"/"+PROJECT_FOLDER+"/images", os.ModePerm)
	os.Mkdir(outputPath+"/"+RESOURCES, os.ModePerm)

	logger.Debug("Copying resources")
	copyDirectoryRecursive(RESOURCES, outputPath+"/"+RESOURCES)
	logger.Debug("Copying static pages")
	copyDirectoryRecursive(STATICPAGES, outputPath)

	logger.Debug("Creating project pages")
	generateProjectPages(projects)

	logger.Debug("Creating single pages")
	generateSinglePages(projects)
}

func generateSinglePages(projects []model.Project) {
	singlePages, err := ioutil.ReadDir("templates/singlePages")
	if(err != nil){
		logger.Warning("Single pages directory not found")
		return
	}

	for _, v := range singlePages {
		if v.IsDir() || !strings.HasSuffix(v.Name(), "gohtml") {
			continue
		}

		templateName := v.Name()
		outputName := strings.Replace(v.Name(), "gohtml", "html", -1)

		templ, err := createTemplateForPage("templates/singlePages/" + templateName)
		if(err != nil){
			logger.Warningf("Failed to parse template %s, page will not be generated", templateName)
		}
		fileOuput, _ := os.Create(outputPath + "/" + outputName)
		error := templ.ExecuteTemplate(fileOuput, "base", model.PageData{
			Projects: projects,
			Context:  context,
		})
		if error != nil {
			logger.Warningf("Failed to execute template %s, page will not be generated", templateName)
		}
		fileOuput.Close()
	}
}

func generateProjectPages(projects []model.Project) {
	templ, err := createTemplateForPage("templates/ProjectViewTemplate.gohtml")
	if(err != nil){
		logger.Warning("Failed to parse template for projects, project pages will not be generated")
	}
	for _, v := range projects {
		generateProject(v, templ)
	}
}

func generateProject(project model.Project, templ *template.Template){
	fileOuput, err := os.Create(outputPath + "/" + PROJECT_FOLDER + "/" + project.FolderPathLowerCase + ".html")
	if err != nil {
		logger.Warningf("Failed to create output file for project %s", project.FolderPath)

	}

	imagePath := outputPath + "/" + PROJECT_FOLDER + "/images/" + project.FolderPathLowerCase

	os.Mkdir(imagePath, os.ModePerm)
	//copy image resources
	for _, image := range project.Images {
		copyFile(PROJECT_FOLDER+"/"+project.FolderPath+"/images/"+image, imagePath+"/"+image)
	}

	error := templ.ExecuteTemplate(fileOuput, "base", model.ProjectDetailsPage{
		Project: project,
		Context: context,
	})
	if error != nil {
		logger.Warningf("Failed to execute project template for for project %s", project.FolderPath)

	}
	logger.Debugf("Generated project html for %s", project.FolderPath)
	fileOuput.Close()
}

func copyDirectoryRecursive(src, dst string) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		logger.Warningf("Failed to read directory for source %s", src)
		return
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.Mkdir(dst, os.ModePerm)
	}

	for _, v := range files {
		if v.IsDir() {
			copyDirectoryRecursive(src+"/"+v.Name(), dst+"/"+v.Name())
			continue
		}
		copyFile(src+"/"+v.Name(), dst+"/"+v.Name())
	}
}

func copyFile(src, dst string) {
	reader, err := os.Open(src)
	if err != nil {
		logger.Warningf("Failed to open file from %s", dst)
		return
	}
	writer, err := os.Create(dst)
	if err != nil {
		logger.Warningf("Failed to write file to %s", dst)
		return
	}
	io.Copy(writer, reader)
	logger.Debugf("Copied file from %s to %s", src, dst)
}

func loadProjects(projectDirs []os.FileInfo) []model.Project {
	logger.Info("Loading projects")
	var projects []model.Project

	for _, f := range projectDirs {
		if !f.IsDir() {
			continue
		}

		projects = append(projects, loadProject(f))
	}
	logger.Info("Successfully loaded projects")
	return projects
}

func loadProject(f os.FileInfo) model.Project {
	project := model.Project{
		FolderPath:          f.Name(),
		FolderPathLowerCase: strings.ToLower(f.Name()),
		Title:               strings.Replace(f.Name(), "_", " ", -1),
		Technologies:        loadProjectTechFile(f),
		Summary:             loadProjectSummaryFile(f),
		ShortSummary:        loadProjectShortSummaryFile(f),
		LearningOutcomes:    loadLearningOutcomesFile(f),
		ProjectOutcomes:     loadProjectOutcomesFile(f),
		Images:              loadProjectImagePaths(f),
		HighlightedProject:  isHighlightedProject(f),
	}
	placeHolderPath := ""

	for _, image := range project.Images {
		matched, _ := regexp.Match("preview\\.[A-Za-z0-9]", []byte(image))
		if matched {
			placeHolderPath = "projects/images/" + project.FolderPathLowerCase + "/" + image
		}
	}

	if placeHolderPath == "" {
		placeHolderPath = "resources/img/placeholder-profile-sq.jpg"
	}
	project.PlaceholderImage = context + placeHolderPath
	return project
}

func isHighlightedProject(f os.FileInfo) bool {
	_, err := os.Stat(PROJECT_FOLDER + "/" + f.Name() + "/highlighted.txt")
	return !os.IsNotExist(err)
}

func loadProjectImagePaths(f os.FileInfo) []string {
	imageFiles, err := ioutil.ReadDir(PROJECT_FOLDER + "/" + f.Name() + "/" + IMAGEFOLDERNAME)
	if err != nil {
		logger.Warningf("Image directory for project %s is missing", f.Name())
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
func loadProjectOutcomesFile(f os.FileInfo) template.HTML {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + PROJECTOUTCOMES
	return template.HTML(loadTextFile(path))
}
func loadLearningOutcomesFile(f os.FileInfo) template.HTML {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + LEARNINGOUTCOMES
	return template.HTML(loadTextFile(path))
}

func loadProjectSummaryFile(f os.FileInfo) template.HTML {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + SUMMARYFILENAME
	return template.HTML(loadTextFile(path))
}

func loadProjectShortSummaryFile(f os.FileInfo) string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + SHORTSUMMARYFILENAME
	return loadTextFile(path)
}

func loadTextFile(src string) string {
	file, err := os.Open(src)
	if err != nil {
		logger.Warningf("Failed to open text file %s", src)
		return ""
	}
	val, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Warningf("Failed to read text file %s", src)
		return ""
	}
	return string(val)
}

func loadProjectTechFile(f os.FileInfo) []string {
	path := PROJECT_FOLDER + "/" + f.Name() + "/" + TECHFILENAME
	file, err := os.Open(path)

	if err != nil {
		logger.Warningf("Failed to open csv file %s", path)
		return []string{}
	}

	reader := csv.NewReader(file)
	result, err := reader.ReadAll()
	if err != nil {
		logger.Warningf("Failed to parse csv file %s", path)
		return []string{}
	}

	var flat []string
	for _, v := range result {
		flat = append(flat, v...)
	}

	return flat
}

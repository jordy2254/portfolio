# Copyright
This file and the containing repository are covered by copy right, and are public for reference only.

&copy; 2020 Jordan Hellier All rights reserved

# Portfolio Site generator
This application is a static site generator for my portfolio page. Projects are created with a specific 
structure of files/folders which are then interpreted and built into project pages. The reason for doing
this was to allow easy editing of pages (IE: updating the nav, adding/editing a project) without the need 
to make changes in multiple locations.

Not only did this project serve creating an easy to maintain portfolio structure it allowed me to apply
and reenforce knowledge about html templates within go.

# Project Structure
A series of files, and folders make up the data for the website, below is the explanations of each of these
areas.

## Project storage
Project are stored within a folder called 'projects' in the same file as the main go file (or exe) this folder
contains snake case project titles for each project. The folder name is used to populate the project title. 
Inside the project folder is a folder for any images along with several text files, below is an explanation of these.

| File name | Description |
| --- | ----------- |
| hightlighted.txt | An empty file turned into a boolean, used to determine if the project should be featured on the homepage |
| learningoutcomes.txt | What was learnt from the project |
| projectoutcomes.txt | The final outcome of the project |
| shortsummary.txt | Used within previews of the project |
| summary.txt | An overall summary of the project |
| technologies.txt | Rough CSV formatting of the different technologies used, Intended for future use to allow filtering of projects |
> **_NOTE:_**  In some events not all files are used

## Templates
There is a master template containing all the basic information for the page with a template of content which needs to be defined.
In addition to this there is also a nav template (nav.gohtml) which holds the navigation bar and is included in the same page. 
For any single pages IE: Contact, index etc. There is a folder featuring this data called singlePages which are generated
and outputted to the root of the build folder. Finally, there is a project view template which puts together the information
about the project from the project folder structure defined above.

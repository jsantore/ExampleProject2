package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	app2 "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
	"log"
	"math/rand"
	"reflect"
	"time"
)

var AllJobs []JobInfo
var MainWindow JobWindow

type JobWindow struct {
	DataDisplay        *widget.List
	CompanyDisplay     *widget.Entry
	PostingDateDisplay *widget.Entry
	LocationDisplay    *widget.Entry
	MaxSalaryDisplay   *widget.Entry
	MinSalaryDisplay   *widget.Entry
	SalaryTypeDisplay  *widget.RadioGroup
	JobTitleDisplay    *widget.Entry
	CurrentSelection   int
}

type JobInfo struct {
	Company     string
	PostingDate string
	JobID       string
	Country     string
	Location    string
	LinuxDate   string
	SalaryMax   string
	SalaryMin   string
	SalaryType  string
	JobTitle    string
}

func main() {
	excelData := GetData("Project2Data.xlsx")
	AllJobs = make([]JobInfo, len(excelData), 1200)
	processJobs(excelData, &AllJobs)
	MainWindow = JobWindow{}
	app := app2.New()
	fyneWindow := app.NewWindow("Your Next Job?")
	makeJobWindow(&MainWindow, fyneWindow)
	fyneWindow.ShowAndRun()
	saveOnExit()
}

func processJobs(exceldata [][]string, jobData *[]JobInfo) {
	for line, excelLine := range exceldata {
		if line < 1 {
			continue //skip the headers
		}
		thisjob := JobInfo{
			Company:     excelLine[0],
			PostingDate: excelLine[1],
			JobID:       excelLine[2],
			Country:     excelLine[3],
			Location:    excelLine[4],
			LinuxDate:   excelLine[5],
			SalaryMax:   excelLine[6],
			SalaryMin:   excelLine[7],
			SalaryType:  excelLine[8],
			JobTitle:    excelLine[9],
		}
		(*jobData)[line] = thisjob
	}
}

func GetData(fileName string) [][]string {
	excelFile, err := excelize.OpenFile(fileName)
	defer excelFile.Close()
	if err != nil {
		log.Fatalln("couldn't open file", err)
	}
	all_rows, err := excelFile.GetRows("Comp490 Jobs")
	if err != nil {
		log.Fatalln(err)
	}
	return all_rows
}

func saveJob() {

	newJob := JobInfo{
		Company:     MainWindow.CompanyDisplay.Text,
		PostingDate: time.Now().Format("2024-03-30"),
		JobID:       RandStringRunes(20),
		Country:     "us",
		Location:    MainWindow.LocationDisplay.Text,
		LinuxDate:   fmt.Sprintf("%d", time.Now().Unix()),
		SalaryMax:   MainWindow.MaxSalaryDisplay.Text,
		SalaryMin:   MainWindow.MinSalaryDisplay.Text,
		SalaryType:  MainWindow.SalaryTypeDisplay.Selected,
		JobTitle:    MainWindow.JobTitleDisplay.Text,
	}
	AllJobs = append(AllJobs, newJob)
}

// from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//end from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

func deleteJob() {
	AllJobs = append(AllJobs[:MainWindow.CurrentSelection], AllJobs[MainWindow.CurrentSelection+1:]...)
	MainWindow.CurrentSelection = 0
	MainWindow.DataDisplay.Select(MainWindow.CurrentSelection)
	MainWindow.DataDisplay.Refresh()
	MainWindow.CompanyDisplay.Text = ""
	MainWindow.CompanyDisplay.Refresh()
	MainWindow.JobTitleDisplay.Text = ""
	MainWindow.JobTitleDisplay.Refresh()
	MainWindow.LocationDisplay.Text = ""
	MainWindow.LocationDisplay.Refresh()
	MainWindow.MinSalaryDisplay.Text = ""
	MainWindow.MinSalaryDisplay.Refresh()
	MainWindow.MaxSalaryDisplay.Text = ""
	MainWindow.MaxSalaryDisplay.Refresh()
	MainWindow.SalaryTypeDisplay.SetSelected("N/A")

}

func updateJob() {
	AllJobs[MainWindow.CurrentSelection].JobTitle = MainWindow.JobTitleDisplay.Text
	AllJobs[MainWindow.CurrentSelection].SalaryMax = MainWindow.MaxSalaryDisplay.Text
	AllJobs[MainWindow.CurrentSelection].SalaryMin = MainWindow.MinSalaryDisplay.Text
	AllJobs[MainWindow.CurrentSelection].Company = MainWindow.CompanyDisplay.Text
	AllJobs[MainWindow.CurrentSelection].Location = MainWindow.LocationDisplay.Text
	AllJobs[MainWindow.CurrentSelection].SalaryType = MainWindow.SalaryTypeDisplay.Selected
}

func saveOnExit() {
	outputFile := excelize.NewFile()
	outputFile.NewSheet("Comp490 Jobs")
	//from https://www.kelche.co/blog/go/excel/
	headers := []string{"Company Name", "Posting Age", "Job Id", "Country", "Location", "Publication Date", "Salary Max",
		"Salary Min", "Salary Type", "Job Title"}
	for i, header := range headers {
		outputFile.SetCellValue("Comp490 Jobs", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)

	}
	for arrayLoc, Job := range AllJobs {
		line := arrayLoc + 1
		if line < 2 {
			continue //zero based and move to after the headings
		}
		jobValue := reflect.ValueOf(Job) //cheating here since I haven't covered reflection yet
		//jobType := jobValue.Type()
		for fieldNum := 0; fieldNum < jobValue.NumField(); fieldNum++ {
			outputFile.SetCellValue("Comp490 Jobs", fmt.Sprintf("%s%d", string(rune(65+fieldNum)), line),
				jobValue.Field(fieldNum))
		}

	}
	if err := outputFile.SaveAs("Project2Data.xlsx"); err != nil {
		log.Fatalln("Error Saving", err)
	}
}

func makeJobWindow(jobDisplay *JobWindow, window fyne.Window) {
	jobDisplay.DataDisplay = widget.NewList(GetNumJobs, CreateListItem, UpdateListItem)

	right_pane := container.NewGridWithColumns(2)
	nextLabel := widget.NewLabel("Company:")
	jobDisplay.CompanyDisplay = widget.NewEntry()
	right_pane.Add(nextLabel)
	right_pane.Add(jobDisplay.CompanyDisplay)
	nextLabel = widget.NewLabel("Job Title:")
	jobDisplay.JobTitleDisplay = widget.NewEntry()
	right_pane.Add(nextLabel)
	right_pane.Add(jobDisplay.JobTitleDisplay)
	nextLabel = widget.NewLabel("Job Location:")
	jobDisplay.LocationDisplay = widget.NewEntry()
	right_pane.Add(nextLabel)
	right_pane.Add(jobDisplay.LocationDisplay)
	nextLabel = widget.NewLabel("Min Salary")
	jobDisplay.MinSalaryDisplay = widget.NewEntry()
	right_pane.Add(nextLabel)
	right_pane.Add(jobDisplay.MinSalaryDisplay)
	nextLabel = widget.NewLabel("Max Salary")
	jobDisplay.MaxSalaryDisplay = widget.NewEntry()
	right_pane.Add(nextLabel)
	right_pane.Add(jobDisplay.MaxSalaryDisplay)
	nextLabel = widget.NewLabel("Salary type")
	validSalaryTypes := []string{"hourly", "yearly", "N/A"}
	jobDisplay.SalaryTypeDisplay = widget.NewRadioGroup(validSalaryTypes, func(s string) {
		//do nothing when changed
	})
	right_pane.Add(nextLabel)
	right_pane.Add(jobDisplay.SalaryTypeDisplay)
	buttonZone := container.NewGridWithColumns(3)
	saveButton := widget.NewButton("Save", saveJob)
	deleteButton := widget.NewButton("Delete", deleteJob)
	UpdateButton := widget.NewButton("Update", updateJob)
	buttonZone.Add(saveButton)
	buttonZone.Add(deleteButton)
	buttonZone.Add(UpdateButton)
	bigRight := container.NewVSplit(right_pane, buttonZone)
	contentPane := container.NewHSplit(jobDisplay.DataDisplay, bigRight)
	window.SetContent(contentPane)
	window.Resize(fyne.NewSize(1200, 900))
}

func UpdateListItem(id widget.ListItemID, object fyne.CanvasObject) {
	listButton := object.(*widget.Button)
	job := AllJobs[id]
	listButton.SetText(job.Company + " : " + job.JobTitle)
	jobSelected := func() {
		MainWindow.CompanyDisplay.Text = job.Company
		MainWindow.CompanyDisplay.Refresh()
		MainWindow.JobTitleDisplay.Text = job.JobTitle
		MainWindow.JobTitleDisplay.Refresh()
		MainWindow.LocationDisplay.Text = job.Location
		MainWindow.LocationDisplay.Refresh()
		MainWindow.MinSalaryDisplay.Text = job.SalaryMin
		MainWindow.MinSalaryDisplay.Refresh()
		MainWindow.MaxSalaryDisplay.Text = job.SalaryMax
		MainWindow.MaxSalaryDisplay.Refresh()
		MainWindow.SalaryTypeDisplay.SetSelected(job.SalaryType)
		MainWindow.CurrentSelection = id
	}
	listButton.OnTapped = jobSelected

}

func CreateListItem() fyne.CanvasObject {
	return widget.NewButton("Someday I will be a Job", func() {
		return
	})
}

func GetNumJobs() int {
	return len(AllJobs)
}

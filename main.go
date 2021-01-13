package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

const (
	baseUrl = ""
)

var (
	username = ""
	password = ""
	token    = ""
)

type App struct {
	Client *http.Client
	Token  string
}

type Project struct {
	Name string
}

type JobCollection struct {
	JobsList []JobData
}

type JobData struct {
	Status    string
	Dataset   string
	User      string
	StartTime time.Time
	Duration  time.Duration
	EndTime   time.Time
}

func (app *App) setToken(resp *http.Response) {

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln("Error loading HTTP response body. ", err)
	}

	token = document.Find("token").Text()

}

func (app *App) login() {
	client := app.Client

	loginURL := baseUrl + "/login?redirect=%2F"

	data := url.Values{
		"userName": {username},
		"password": {password},
	}

	response, err := client.PostForm(loginURL, data)

	if err != nil {
		log.Fatalln(err)
	}

	app.setToken(response)

	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *App) getJobs() []JobData {
	projectUrl := baseUrl + "/jobs?filters=%7B\"qt\"%3A%5B\"UI\"%2C\"EXTERNAL\"%5D%7D&order=DESCENDING&sort=st#20020a1f-8730-04c3-9981-e24af2714600"
	client := app.Client

	req, err := http.NewRequest("GET", projectUrl, nil)

	req.Header.Add("authorization", token)
	response, err := client.Do(req)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatalln("Error loading HTTP response body. ", err)
	}

	var jobs []JobData

	document.Find(".jobs-table-tr").Each(func(i int, selection *goquery.Selection) {
		job := JobData{}

		status, ok := selection.Children().Children().Attr("TextHighlight")
		if ok {
			println("status found: ", status)
			job.Status = status
		}

		dataset, ok := selection.Attr("TextHighlight")
		if ok {
			job.Dataset = dataset
		}

	})

	return jobs
}

func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	app.login()
	jobs := app.getJobs()

	for index, job := range jobs {
		fmt.Println(index+1, job.Status)
	}
}

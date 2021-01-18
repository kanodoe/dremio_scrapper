package main

import (
	"bytes"
	"encoding/json"
	"github.com/andybrewer/mack"
	"github.com/gen2brain/beeep"
	"github.com/prprprus/scheduler"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"runtime"
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
}

type DremioResponse struct {
	Token                     string `json:"token"`
	UserName                  string `json:"userName"`
	FirstName                 string `json:"firstName"`
	LastName                  string `json:"lastName"`
	Expires                   int64  `json:"expires"`
	Email                     string `json:"email"`
	UserID                    string `json:"userId"`
	Admin                     bool   `json:"admin"`
	ClusterID                 string `json:"clusterId"`
	ClusterCreatedAt          int64  `json:"clusterCreatedAt"`
	ShowUserAndUserProperties bool   `json:"showUserAndUserProperties"`
	Version                   string `json:"version"`
	Permissions               struct {
		CanUploadProfiles   bool `json:"canUploadProfiles"`
		CanDownloadProfiles bool `json:"canDownloadProfiles"`
		CanEmailForSupport  bool `json:"canEmailForSupport"`
		CanChatForSupport   bool `json:"canChatForSupport"`
	} `json:"permissions"`
	UserCreatedAt int64 `json:"userCreatedAt"`
}

type Jobs struct {
	Jobs []struct {
		ID          string `json:"id"`
		State       string `json:"state"`
		FailureInfo struct {
			Errors []interface{} `json:"errors"`
			Type   string        `json:"type"`
		} `json:"failureInfo"`
		User                 string   `json:"user"`
		StartTime            int64    `json:"startTime"`
		EndTime              int64    `json:"endTime"`
		Description          string   `json:"description"`
		DatasetPathList      []string `json:"datasetPathList"`
		DatasetType          string   `json:"datasetType"`
		RequestType          string   `json:"requestType"`
		Accelerated          bool     `json:"accelerated"`
		DatasetVersion       string   `json:"datasetVersion"`
		SnowflakeAccelerated bool     `json:"snowflakeAccelerated"`
		Spilled              bool     `json:"spilled"`
		OutputRecords        int      `json:"outputRecords"`
		OutputLimited        bool     `json:"outputLimited"`
		IsComplete           bool     `json:"isComplete"`
	} `json:"jobs"`
	Next string `json:"next"`
}

func (app *App) login() {
	client := app.Client

	loginURL := baseUrl + "/apiv2/login"

	data, _ := json.Marshal(map[string]string{
		"userName": username,
		"password": password,
	})

	responseBody := bytes.NewBuffer(data)

	response, err := client.Post(loginURL, "application/json", responseBody)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var dremioResponse = new(DremioResponse)

	body, err := ioutil.ReadAll(response.Body)
	errJson := json.Unmarshal(body, &dremioResponse)
	if errJson != nil {
		log.Fatalln(errJson)
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	token = dremioResponse.Token
}

func (app *App) getJobs() *Jobs {
	projectUrl := baseUrl + "/apiv2/jobs/?sort=st&order=DESCENDING&filter=(qt%3D%3D%22UI%22%2Cqt%3D%3D%22EXTERNAL%22)"
	client := app.Client

	req, err := http.NewRequest("GET", projectUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Authorization", "_dremio"+token)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var jobs = new(Jobs)

	body, err := ioutil.ReadAll(response.Body)

	errJson := json.Unmarshal(body, jobs)
	if errJson != nil {
		log.Fatalln(errJson)
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	app.filterJobs(jobs)

	return jobs
}

func (app *App) filterJobs(jobs *Jobs) {

	for _, v := range jobs.Jobs {
		if v.State == "FAILED" {

			println("encontre un error en: ", v.ID)

			if runtime.GOOS == "darwin" {
				errSay := mack.Say("Alerta  ha ocurrido un error de reflexion en dremio!")
				if errSay != nil {
					log.Fatalln(errSay)
				}
				_, errAlert := mack.Alert("Alerta", "Ha ocurrido un error de la reflexion "+v.FailureInfo.Type, "critical")
				if errAlert != nil {
					log.Fatalln(errAlert)
				}

				errNotify := mack.Notify("Favor de ver reflexiones en dremio", "WIC", "Alerta", "Ping")
				if errNotify != nil {
					log.Fatalln(errNotify)
				}
			} else if runtime.GOOS == "windows" {
				err := beeep.Alert("Alerta", "Ha ocurrido un error de la reflexion "+v.FailureInfo.Type, "assets/warning.png")
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		log.Fatalln(err)
	}

	if len(token) < 1 {
		app.login()
	}

	s.Delay().Minute(10).Do(app.getJobs())

}

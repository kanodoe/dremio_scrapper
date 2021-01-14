package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

const (
	baseUrl = "[base.url.here]"
)

var (
	username = "[your.username]"
	password = "[your.password]"
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

	defer response.Body.Close()

	var dremioResponse = new(DremioResponse)

	body, err := ioutil.ReadAll(response.Body)
	errJson := json.Unmarshal(body, &dremioResponse)
	if errJson != nil {
		log.Fatalln(errJson)
	}

	log.Printf("s = %v", dremioResponse)

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

	req.Header.Add("Authorization", "_dremio"+token)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	var jobs = new(Jobs)

	body, err := ioutil.ReadAll(response.Body)
	log.Println("Body response: ", string(body))
	errJson := json.Unmarshal(body, jobs)
	if errJson != nil {
		log.Fatalln(errJson)
	}

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return jobs

}

func (jobs *Jobs) FromJson(response []byte) error {
	var data = &jobs.Jobs
	return json.Unmarshal(response, data)
}

func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	app.login()
	jobs := app.getJobs()

	log.Printf("s = %v", jobs)
}

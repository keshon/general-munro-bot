// Package kitsu provides methods for Kitsu task management software
package kitsu

import (
	"bot/src/controllers/config"
	"bot/src/utils/debug"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Task struct {
	Assignees       []string    `json:"assignees"`
	ID              string      `json:"id"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
	Name            string      `json:"name"`
	Description     interface{} `json:"description"`
	Priority        int         `json:"priority"`
	Duration        int         `json:"duration"`
	Estimation      int         `json:"estimation"`
	CompletionRate  int         `json:"completion_rate"`
	RetakeCount     int         `json:"retake_count"`
	SortOrder       int         `json:"sort_order"`
	StartDate       interface{} `json:"start_date"`
	EndDate         interface{} `json:"end_date"`
	DueDate         interface{} `json:"due_date"`
	RealStartDate   interface{} `json:"real_start_date"`
	LastCommentDate interface{} `json:"last_comment_date"`
	Data            interface{} `json:"data"`
	ShotgunID       interface{} `json:"shotgun_id"`
	ProjectID       string      `json:"project_id"`
	TaskTypeID      string      `json:"task_type_id"`
	TaskStatusID    string      `json:"task_status_id"`
	EntityID        string      `json:"entity_id"`
	AssignerID      string      `json:"assigner_id"`
	Type            string      `json:"type"`
}

type Person struct {
	ID                        string `json:"id"`
	CreatedAt                 string `json:"created_at"`
	UpdatedAt                 string `json:"updated_at"`
	FirstName                 string `json:"first_name"`
	LastName                  string `json:"last_name"`
	Email                     string `json:"email"`
	Phone                     string `json:"phone"`
	Active                    bool   `json:"active"`
	LastPresence              string `json:"last_presence"`
	DesktopLogin              string `json:"desktop_login"`
	ShotgunID                 string `json:"shotgun_id"`
	Timezone                  string `json:"timezone"`
	Locale                    string `json:"locale"`
	Data                      string `json:"data"`
	Role                      string `json:"role"`
	HasAvatar                 bool   `json:"has_avatar"`
	NotificationsEnabled      bool   `json:"notifications_enabled"`
	NotificationsSlackEnabled bool   `json:"notifications_slack_enabled"`
	NotificationsSlackUserid  string `json:"notifications_slack_userid"`
	Type                      string `json:"type"`
	FullName                  string `json:"full_name"`
}

type Entity struct {
	EntitiesOut     []interface{} `json:"entities_out"`
	InstanceCasting []interface{} `json:"instance_casting"`
	CreatedAt       string        `json:"created_at"`
	UpdatedAt       string        `json:"updated_at"`
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Code            interface{}   `json:"code"`
	Description     interface{}   `json:"description"`
	ShotgunID       interface{}   `json:"shotgun_id"`
	Canceled        bool          `json:"canceled"`
	NbFrames        interface{}   `json:"nb_frames"`
	ProjectID       string        `json:"project_id"`
	EntityTypeID    string        `json:"entity_type_id"`
	ParentID        string        `json:"parent_id"`
	SourceID        interface{}   `json:"source_id"`
	PreviewFileID   interface{}   `json:"preview_file_id"`
	Data            interface{}   `json:"data"`
	EntitiesIn      []interface{} `json:"entities_in"`
	Type            string        `json:"type"`
}

type TaskStatus struct {
	ID              string      `json:"id"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
	Name            string      `json:"name"`
	ShortName       string      `json:"short_name"`
	Color           string      `json:"color"`
	IsDone          bool        `json:"is_done"`
	IsArtistAllowed bool        `json:"is_artist_allowed"`
	IsClientAllowed bool        `json:"is_client_allowed"`
	IsRetake        bool        `json:"is_retake"`
	ShotgunID       interface{} `json:"shotgun_id"`
	IsReviewable    bool        `json:"is_reviewable"`
	Type            string      `json:"type"`
}

func GetTask(taskID string) Task {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/tasks/"+taskID, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(os.Getenv("JWTToken"))
	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Task

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		log.Fatalln(err)
	}

	return typBody
}

func GetPerson(personID string) Person {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/persons/"+personID, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(respBody))

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Person

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		log.Fatalln(err)
	}

	return typBody
}

func GetEntity(EntityID string) Entity {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/entities/"+EntityID, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Entity

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		log.Fatalln(err)
	}

	return typBody
}

func GetTaskStatus(TaskStatusID string) TaskStatus {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/task-status/"+TaskStatusID, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody TaskStatus

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		log.Fatalln(err)
	}

	return typBody
}

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
	Assignees       []string    `json:"assignees,omitempty"`
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	Description     interface{} `json:"description,omitempty"`
	Priority        int         `json:"priority,omitempty"`
	Duration        int         `json:"duration,omitempty"`
	Estimation      int         `json:"estimation,omitempty"`
	CompletionRate  int         `json:"completion_rate,omitempty"`
	RetakeCount     int         `json:"retake_count,omitempty"`
	SortOrder       int         `json:"sort_order,omitempty"`
	StartDate       interface{} `json:"start_date,omitempty"`
	EndDate         interface{} `json:"end_date,omitempty"`
	DueDate         interface{} `json:"due_date,omitempty"`
	RealStartDate   interface{} `json:"real_start_date,omitempty"`
	LastCommentDate string      `json:"last_comment_date,omitempty"`
	Data            interface{} `json:"data,omitempty"`
	ShotgunID       interface{} `json:"shotgun_id,omitempty"`
	ProjectID       string      `json:"project_id,omitempty"`
	TaskTypeID      string      `json:"task_type_id,omitempty"`
	TaskStatusID    string      `json:"task_status_id,omitempty"`
	EntityID        string      `json:"entity_id,omitempty"`
	AssignerID      string      `json:"assigner_id,omitempty"`
	Type            string      `json:"type,omitempty"`
}
type Tasks struct {
	Each []Task
}

type Person struct {
	ID                        string `json:"id,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	UpdatedAt                 string `json:"updated_at,omitempty"`
	FirstName                 string `json:"first_name,omitempty"`
	LastName                  string `json:"last_name,omitempty"`
	Email                     string `json:"email,omitempty"`
	Phone                     string `json:"phone,omitempty"`
	Active                    bool   `json:"active,omitempty"`
	LastPresence              string `json:"last_presence,omitempty"`
	DesktopLogin              string `json:"desktop_login,omitempty"`
	ShotgunID                 string `json:"shotgun_id,omitempty"`
	Timezone                  string `json:"timezone,omitempty"`
	Locale                    string `json:"locale,omitempty"`
	Data                      string `json:"data,omitempty"`
	Role                      string `json:"role,omitempty"`
	HasAvatar                 bool   `json:"has_avatar,omitempty"`
	NotificationsEnabled      bool   `json:"notifications_enabled,omitempty"`
	NotificationsSlackEnabled bool   `json:"notifications_slack_enabled,omitempty"`
	NotificationsSlackUserid  string `json:"notifications_slack_userid,omitempty"`
	Type                      string `json:"type,omitempty"`
	FullName                  string `json:"full_name,omitempty"`
}

type Entity struct {
	EntitiesOut     []interface{} `json:"entities_out,omitempty"`
	InstanceCasting []interface{} `json:"instance_casting,omitempty"`
	CreatedAt       string        `json:"created_at,omitempty"`
	UpdatedAt       string        `json:"updated_at,omitempty"`
	ID              string        `json:"id,omitempty"`
	Name            string        `json:"name,omitempty"`
	Code            interface{}   `json:"code,omitempty"`
	Description     interface{}   `json:"description,omitempty"`
	ShotgunID       interface{}   `json:"shotgun_id,omitempty"`
	Canceled        bool          `json:"canceled,omitempty"`
	NbFrames        interface{}   `json:"nb_frames,omitempty"`
	ProjectID       string        `json:"project_id,omitempty"`
	EntityTypeID    string        `json:"entity_type_id,omitempty"`
	ParentID        string        `json:"parent_id,omitempty"`
	SourceID        interface{}   `json:"source_id,omitempty"`
	PreviewFileID   interface{}   `json:"preview_file_id,omitempty"`
	Data            interface{}   `json:"data,omitempty"`
	EntitiesIn      []interface{} `json:"entities_in,omitempty"`
	Type            string        `json:"type,omitempty"`
}

type Entities struct {
	Each []Entity
}

type TaskStatus struct {
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	ShortName       string      `json:"short_name,omitempty"`
	Color           string      `json:"color,omitempty"`
	IsDone          bool        `json:"is_done,omitempty"`
	IsArtistAllowed bool        `json:"is_artist_allowed,omitempty"`
	IsClientAllowed bool        `json:"is_client_allowed,omitempty"`
	IsRetake        bool        `json:"is_retake,omitempty"`
	ShotgunID       interface{} `json:"shotgun_id,omitempty"`
	IsReviewable    bool        `json:"is_reviewable,omitempty"`
	Type            string      `json:"type,omitempty"`
}

type Comment struct {
	ID        string      `json:"id,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	ShotgunID interface{} `json:"shotgun_id,omitempty"`
	ObjectID  string      `json:"object_id,omitempty"`
	PersonID  string      `json:"person_id,omitempty"`
	Text      string      `json:"text,omitempty"`
}

type Comments struct {
	Each []Comment
}

func GetComment(objectID string) Comments {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/comments?object_id="+objectID, nil)
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
	var typBody Comments

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
	if err != nil {
		log.Fatalln(err)
	}

	return typBody
}

func GetTasks() Tasks {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/tasks/", nil)
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
	var typBody Tasks

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
	if err != nil {
		log.Fatalln(err)
	}

	return typBody
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

func GetEntities(EntityID string) Entities {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/entities/", nil)
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
	var typBody Entities

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
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

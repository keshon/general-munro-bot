// Package kitsu provides methods for Kitsu task management software
package kitsu

import (
	"bot/src/controllers/config"
	"bot/src/utils/debug"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Task struct {
	Assignees       []string    `json:"assignees,omitempty"`
	ID              string      `json:"id,omitempty"`
	CreatedAt       string      `json:"created_at,omitempty"`
	UpdatedAt       string      `json:"updated_at,omitempty"`
	Name            string      `json:"name,omitempty"`
	LastCommentDate string      `json:"last_comment_date,omitempty"`
	Data            interface{} `json:"data,omitempty"`
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

type Attachment struct {
	ID        string `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Name      string `json:"name,omitempty"`
	Size      int    `json:"size,omitempty"`
	Extension string `json:"extension,omitempty"`
	Mimetype  string `json:"mimetype,omitempty"`
	CommentID string `json:"comment_id,omitempty"`
	Comment   struct {
		ObjectID   string `json:"object_id,omitempty"`
		ObjectType string `json:"object_type,omitempty"`
	}
}

type Attachments struct {
	Each []Attachment
}

type TaskType struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortName string `json:"short_name,omitempty"`
}

type TaskTypes struct {
	Each []TaskType
}

type Project struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	ProjectStatusID string `json:"project_status_id,omitempty"`
}

type Projects struct {
	Each []Project
}

type ProjectStatus struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type ProjectStatuses struct {
	Each []ProjectStatus
}

func GetComment(objectID string) Comments {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/comments?object_id="+objectID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Comments

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
	if err != nil {
		//return Comments{}
		panic(err)
	}

	return typBody
}

func GetTasks() Tasks {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/tasks/", nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Tasks

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
	if err != nil {
		//return Tasks{}
		panic(err)
	}

	return typBody
}
func GetTask(taskID string) Task {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/tasks/"+taskID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Task

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return Task{}
		panic(err)
	}

	return typBody
}

func GetPerson(personID string) Person {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/persons/"+personID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(respBody))

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Person

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return Person{}
		panic(err)
	}

	return typBody
}

func GetEntities(EntityID string) Entities {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/entities/", nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Entities

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
	if err != nil {
		//return Entities{}
		panic(err)
	}

	return typBody
}

func GetEntity(EntityID string) Entity {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/entities/"+EntityID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Entity

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return Entity{}
		panic(err)
	}

	return typBody
}

func GetTaskStatus(TaskStatusID string) TaskStatus {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/task-status/"+TaskStatusID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody TaskStatus

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return TaskStatus{}
		panic(err)
	}

	return typBody
}

func GetAttachments() Attachments {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/attachment-files/", nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Attachments

	err = json.Unmarshal([]byte(strBody), &typBody.Each)
	if err != nil {
		//return Attachments{}
		panic(err)
	}

	return typBody
}

func GetAttachment(AttachmentID string) Attachment {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/attachment-files/"+AttachmentID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Attachment

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return Attachment{}
		panic(err)
	}

	return typBody
}

func DownloadAttachment(localPath, id, filename string, conf config.Config) int64 {
	// Create client
	client := &http.Client{}

	// Create dir
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		err := os.Mkdir(localPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	// Create the file
	out, err := os.Create(localPath + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/attachment-files/"+id+"/file/"+filename, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		//return fmt.Errorf("bad status: %s", resp.Status)
		//panic("bad status:" + resp.Status)
		return 0
	}

	// Writer the body to file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	return size
}

func GetTaskType(taskID string) TaskType {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/task-types/"+taskID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody TaskType

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return TaskType{}
		panic(err)
	}

	return typBody
}

func GetProject(projectID string) Project {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/projects/"+projectID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody Project

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return Project{}
		panic(err)
	}

	return typBody
}

func GetProjectStatus(projectStatusID string) ProjectStatus {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodGet, config.Read().Kitsu.Hostname+"api/data/project-status/"+projectStatusID, nil)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("JWTToken"))

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	debug.Info(resp, respBody)

	// Unmarshal
	strBody := string(respBody)
	var typBody ProjectStatus

	err = json.Unmarshal([]byte(strBody), &typBody)
	if err != nil {
		//return ProjectStatus{}
		panic(err)
	}

	return typBody
}

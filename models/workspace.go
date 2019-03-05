package models

import (
	"github.com/abdullahi/codice/sandbox"
	u "github.com/abdullahi/codice/utils"
	"github.com/satori/go.uuid"
	"os"
)

type WorkSpaceType string

const (
	NodeJS WorkSpaceType = "nodejs"
	PHP    WorkSpaceType = "php"
	Python WorkSpaceType = "python"
	Ruby   WorkSpaceType = "ruby"
)

func (e *WorkSpaceType) Scan(value interface{}) error {
	*e = WorkSpaceType(value.([]byte))
	return nil
}

func (e WorkSpaceType) Value() (string, error) {
	return string(e), nil
}

type Workspace struct {
	GormModel
	Name        string        `json:"name"`
	Path        string        `json:"-"`
	UserID      string        `json:"-"`
	User        User          `json:"owner"`
	ContainerID string        `json:"containerId"`
	Type        WorkSpaceType `json:"type" sql:"type:ENUM('nodejs','php','python','ruby')"`
}

func (workspace *Workspace) Create() map[string]interface{} {
	workspace.ID = uuid.Must(uuid.NewV4()).String()
	workspace.Path = os.Getenv("TEMP_PATH") + "/" + workspace.GetRootDir()

	srcLang, langErr := sandbox.GetLang(string(workspace.Type))
	if langErr != nil {
		return u.Message(false, "Invalid Language, Could not find "+string(workspace.Type))
	}

	sandbox.CreatePayload(workspace.Path, srcLang)
	workspace.ContainerID = sandbox.CreateContainer(workspace.Path)

	if err := GetDB().Create(&workspace).Related(&workspace.User).Error; err != nil {
		return u.Message(false, "Connection error. Please retry")
	}

	response := u.Message(true, "New workspace Created")

	response["workspace"] = workspace

	return response
}

func GetWorkspace(id string) (*Workspace, error) {
	workspace := &Workspace{}
	if err := GetDB().Table("workspaces").Where("id = ?", id).First(&workspace).Error; err != nil {
		return nil, err
	}

	return workspace, nil
}

func (workspace *Workspace) GetRootDir() string {
	return workspace.Name + workspace.ID
}

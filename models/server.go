package models

import (
	"os"
	"time"
	u "webmine/utils"
)

const (
	SERVER_DATA_DIR_PATH = "./data"
	SERVER_DATA_ZIP_NAME = "server_data.zip"
	SERVER_DATA_ZIP_PATH = "./tmp/" + SERVER_DATA_ZIP_NAME

	CMD_DOCKER = "docker"

	STATE_RUNNING    = "Running"
	STATE_RESTARTING = "Restarting"
	STATE_WARNING    = "Not running and restarting"
)

var (
	SERVER_IMG_NAME = "itzg/minecraft-server:java8"

	// 	CMD_SERVER_UPTIME = [...]string{"inspect", "-f", "{{.State.StartedAt}}"}
	// 	CMD_SERVER_STATE  = [...]string{"container", "inspect", "-f", "{{.State.%s}}"}
	// 	CMD_SERVER_RELOAD = [...]string{"restart"}
)

type Server struct {
	ID     string    `json:"-"`
	State  string    `json:"state"`
	Uptime time.Time `json:"uptime"`
}

var server Server

func init() {
	server.ID = os.Getenv("CONTAINER_ID")

	if tmp := os.Getenv("SERVER_IMG_NAME"); tmp != "" {
		SERVER_IMG_NAME = tmp
	} else if tmp = u.DefaultEnv["SERVER_IMG_NAME"]; tmp != "" {
		SERVER_IMG_NAME = tmp
	}
}

func GetServer() *Server {
	return &server
}

// Perhaps you should make a backup before reload
// err := u.RecursiveZip(SERVER_DATA_DIR_PATH, ...)

func (server *Server) Reload() error {
	_, err := u.GetResultOfExecCmd(CMD_DOCKER, "restart", server.ID)
	if err != nil {
		return err
	}
	server.UpdateUptime()
	server.UpdateState()
	return nil
}

func (server *Server) UpdateUptime() error {
	res, err := u.GetResultOfExecCmd(CMD_DOCKER, "inspect", "-f", "{{.State.StartedAt}}", server.ID)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, res)
	server.Uptime = t
	return err
}

func (server *Server) UpdateState() error {
	res, err := u.GetResultOfExecCmd(CMD_DOCKER, "container", "inspect", "-f", "{{.State.Health.Status}}", server.ID)
	if err != nil {
		return err
	}
	if res == "starting" {
		server.State = STATE_RESTARTING
	} else {
		res, err = u.GetResultOfExecCmd(CMD_DOCKER, "container", "inspect", "-f", "{{.State.Health.Status}}", server.ID)
		if err != nil {
			return err
		}
		if res == "healthy" {
			server.State = STATE_RUNNING
		} else {
			server.State = STATE_WARNING
		}
	}
	return nil
}

func (server *Server) UpdateId() error {
	id := os.Getenv("CONTAINER_ID")
	if id == "" {
		res, err := u.GetResultOfExecCmd(CMD_DOCKER, "ps", "-aqf", "ancestor="+SERVER_IMG_NAME)
		server.ID = res
		return err
	}
	server.ID = id
	return nil
}

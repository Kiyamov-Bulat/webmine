package models

import (
	"log"
	"path/filepath"
	"time"
	u "webmine/utils"
)

const (
	SERVER_DATA_DIR_PATH = "./data"
	SERVER_DATA_ZIP_NAME = "server_data.zip"
	SERVER_DATA_ZIP_PATH = "./tmp/" + SERVER_DATA_ZIP_NAME
	SERVER_IMG_NAME      = "itzg/minecraft-server:java8"

	CMD_DOCKER = "docker"

	STATE_RUNNING    = "Running"
	STATE_RESTARTING = "Restarting"
	STATE_WARNING    = "Not running and restarting"
)

// var (
// 	CMD_SERVER_UPTIME = [...]string{"inspect", "-f", "'{{ .Created }}'"}
// 	CMD_SERVER_STATE  = "docker container inspect -f '{{.State.%s}}'"
// 	CMD_SERVER_RELOAD = "docker restart"
// )

type Server struct {
	ID     string    `json:"-"`
	State  string    `json:"state"`
	Uptime time.Time `json:"uptime"`
}

var server Server

func GetServer() *Server {
	return &server
}

func (server *Server) Reload() error {
	err := u.RecursiveZip(SERVER_DATA_DIR_PATH, filepath.Join(SERVER_DATA_DIR_PATH, SERVER_DATA_ZIP_NAME))
	if err != nil {
		log.Println(err.Error())
	}
	_, err = u.GetResultOfExecCmd(CMD_DOCKER, "restart", server.ID)
	if err != nil {
		return err
	}
	server.UpdateUptime()
	server.UpdateState()
	return nil
}

func (server *Server) UpdateUptime() error {
	res, err := u.GetResultOfExecCmd(CMD_DOCKER, "inspect", "-f", "{{ .Created }}", server.ID)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, res)
	server.Uptime = t
	return err
}

func (server *Server) UpdateState() error {
	res, err := u.GetResultOfExecCmd(CMD_DOCKER, "container", "inspect", "-f", "{{.State.Running}}", server.ID)
	if err != nil {
		return err
	}
	if res == "true" {
		server.State = STATE_RUNNING
	} else {
		res, err = u.GetResultOfExecCmd(CMD_DOCKER, "container", "inspect", "-f", "{{.State.Restarting}}", server.ID)
		if err != nil {
			return err
		}
		if res == "true" {
			server.State = STATE_RESTARTING
		} else {
			server.State = STATE_WARNING
		}
	}
	return nil
}

func (server *Server) UpdateId() error {
	res, err := u.GetResultOfExecCmd(CMD_DOCKER, "ps", "-aqf", "ancestor="+SERVER_IMG_NAME)
	server.ID = res
	return err
}

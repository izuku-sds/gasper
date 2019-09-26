package php

import (
	"github.com/sdslabs/SWS/configs"
	"github.com/sdslabs/SWS/lib/api"
	"github.com/sdslabs/SWS/lib/commons"
	"github.com/sdslabs/SWS/lib/docker"
	"github.com/sdslabs/SWS/lib/types"
)

// installPackages installs dependancies for the specific microservice
func installPackages(path string, appEnv *types.ApplicationEnv) (string, types.ResponseError) {
	cmd := []string{"bash", "-c", `composer install -d ` + path}
	execID, err := docker.ExecDetachedProcess(appEnv.Context, appEnv.Client, appEnv.ContainerID, cmd)
	if err != nil {
		return "", types.NewResErr(500, "Failed to perform composer install in the container", err)
	}
	return execID, nil
}

// Pipeline is the application creation pipeline
func Pipeline(data map[string]interface{}) types.ResponseError {
	appConf := &types.ApplicationConfig{
		DockerImage:  configs.ImageConfig["php"].(string),
		ConfFunction: configs.CreatePHPContainerConfig,
	}

	appEnv, resErr := api.SetupApplication(appConf, data)
	if resErr != nil {
		return resErr
	}

	context := data["context"].(map[string]interface{})

	if context["rcFile"].(bool) {
		return nil
	}

	// Perform composer install in the container
	if data["composer"] != nil {
		if data["composer"].(bool) {
			var composerPath string
			if data["composerPath"] != nil {
				composerPath = data["composerPath"].(string)
			} else {
				composerPath = "."
			}
			execID, resErr := installPackages(composerPath, appEnv)
			if resErr != nil {
				go commons.AppFullCleanup(data["name"].(string))
				return resErr
			}
			data["execID"] = execID
		}
	}

	return nil
}
package helpers

import (
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"os"
)

const (
	logsDir = "./logs"
)

func DumpContainerLogs() error {
	client := GetDockerClient()
	//project := fmt.Sprintf("com.docker.compose.project=%s", os.Getenv("COMPOSE_PROJECT_NAME"))
  //logrus.Debugf("PROJECT: %v", project)

	queryAll := docker.ListContainersOptions{All: true} //,
    //Filters: map[string][]string{"label": []string{project}}}

	containers, err := client.ListContainers(queryAll)
	if err != nil {
		logrus.Error("could not query containers to dump logs for", err)
		return err
	}
	logrus.Debugf("found %d container(s) to dump logs for", len(containers))

	if err := createLogsDirectory(); err != nil {
		logrus.Error("could not create directory for logs", err)
		return err
	}

	for _, container := range containers {
		logFileName := fmt.Sprintf("%v/%v.log", logsDir, container.Names[0])
    //logrus.Debugf("dumping log for %v", logFileName)
    logrus.Debugf("container: %v", container)
		if f, err := os.Create(logFileName); err != nil {
			logrus.Error("Unable to create "+logFileName, err)
		} else {
			w := bufio.NewWriter(f)
			client.Logs(docker.LogsOptions{Container: container.ID, Tail: "all", OutputStream: w})
      f.Close()
		}
	}

	return nil
}

func createLogsDirectory() error {
	_, err := os.Stat(logsDir)
	if err == nil {
		return nil
	}

	return os.Mkdir(logsDir, 0755)
}

func getContainerIdsAndNames(contaners *[]docker.APIContainers) []string {
	ids := []string{}
	for _, container := range *contaners {
		ids = append(ids, container.ID)
	}

	return ids
}

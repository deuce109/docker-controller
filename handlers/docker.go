package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/samber/lo"
)

func getDockerClient() (cli *client.Client, err error) {
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return
}

func getContainerNamesAndIds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cli, err := getDockerClient()

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("%s", err)
	} else {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})

		if err != nil {
			w.WriteHeader(500)
			fmt.Printf("%s", err)
		} else {

			result := make(map[string]interface{})

			containerInfo := lo.Map[types.Container, map[string]interface{}](containers, func(item types.Container, index int) map[string]interface{} {

				info := make(map[string]interface{})
				info["Names"] = item.Names
				info["ID"] = item.ID
				return info

			})

			result["result"] = containerInfo

			resultJson, err := json.Marshal(result)

			if err != nil {
				w.WriteHeader(500)
				fmt.Printf("%s", err)
			}

			w.Write(resultJson)

		}

	}

}

func inspectContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	containerId := params["containerId"]

	cli, err := getDockerClient()

	if err != nil {
		fmt.Printf("%s", err)
	} else {
		container, err := cli.ContainerInspect(r.Context(), containerId)

		if err != nil {
			fmt.Printf("%s", err)
			w.WriteHeader(500)
			w.Write([]byte(""))
		}

		containerJson, err := json.Marshal(container)

		if err != nil {
			fmt.Printf("%s", err)
			w.WriteHeader(500)
			w.Write([]byte(""))
		}

		w.Write(containerJson)

	}
}

func startContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	containerId := params["containerId"]

	cli, err := getDockerClient()

	if err != nil {
		fmt.Printf("%s", err)
	} else {
		err := cli.ContainerStart(r.Context(), containerId, types.ContainerStartOptions{})

		if err != nil {
			fmt.Printf("%s", err)
			w.WriteHeader(500)
			w.Write([]byte(""))
		} else {
			w.Write([]byte(fmt.Sprintf("Successfully started container %s", containerId)))
		}

	}
}

func stopContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	containerId := params["containerId"]

	cli, err := getDockerClient()

	if err != nil {
		fmt.Printf("%s", err)
	} else {
		err := cli.ContainerStop(r.Context(), containerId, container.StopOptions{})

		if err != nil {
			fmt.Printf("%s", err)
			w.WriteHeader(500)
			w.Write([]byte(""))
		} else {
			w.Write([]byte(fmt.Sprintf("Successfully stopped container %s", containerId)))

		}

	}
}

func restartContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	containerId := params["containerId"]

	cli, err := getDockerClient()

	if err != nil {
		fmt.Printf("%s", err)
	} else {
		err := cli.ContainerRestart(r.Context(), containerId, container.StopOptions{})

		if err != nil {
			fmt.Printf("%s", err)
			w.WriteHeader(500)
			w.Write([]byte(""))
		} else {
			w.Write([]byte(fmt.Sprintf("Successfully restarted container %s", containerId)))

		}
	}
}

func SetContainerRoutes(r *mux.Router) *mux.Router {

	r.HandleFunc("/containers", getContainerNamesAndIds)
	r.HandleFunc("/containers/{containerId}", inspectContainer)
	r.HandleFunc("/containers/{containerId}/start", startContainer)
	r.HandleFunc("/containers/{containerId}/stop", stopContainer)
	r.HandleFunc("/containers/{containerId}/restart", restartContainer)

	return r
}

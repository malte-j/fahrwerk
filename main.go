package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SkyGuardian42/fahrwerk/cli"
	"github.com/SkyGuardian42/fahrwerk/hetzner"
	"github.com/SkyGuardian42/fahrwerk/k3s"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/theckman/yacspin"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := yacspin.Config{
		Frequency:  100 * time.Millisecond,
		CharSet:    yacspin.CharSets[14],
		Suffix:     " ",
		Colors:     []string{"fgGreen"},
		StopColors: []string{"fgGreen"},
	}
	spinner, _ := yacspin.New(cfg)

	hcloudToken, _ := cli.GetHcloudToken()

	spinner.Message("Connecting to HCloud")
	spinner.Start()
	client := hcloud.NewClient(hcloud.WithToken(hcloudToken))
	spinner.Stop()

	spinner.Message("Loading Server Types")
	spinner.Start()
	types := getAllServerTypes(client)
	spinner.Stop()
	serverType := cli.SelectServerType(types)

	spinner.Message("Loading Keys")
	spinner.Start()
	sshKeys := getAllSSHKeys(client)
	spinner.Stop()
	sshKey := cli.SelectSSHKey(sshKeys)

	cli.SelectNodeAmount()

	spinner.Message("Loading Locations")
	spinner.Start()
	locations := getAllLocations(client)
	spinner.Stop()
	location := cli.SelectLocation(locations)

	spinner.Message("Loading Images")
	spinner.Start()
	images := getAllImages(client)
	spinner.Stop()
	image, _ := cli.SelectImage(images)

	createClusterConfirmation := cli.ConfirmClusterCreation()
	if !createClusterConfirmation {
		fmt.Println("Aborting...")
		return
	}

	spinner.Message("Creating Cluster")
	spinner.Start()
	time.Sleep(3 * time.Second)
	network, res, err := hetzner.CreateNetwork(client)

	if err != nil {
		fmt.Print(res)
		return
	}

	server, _, _ := hetzner.CreateServer(client, hetzner.CreateServerConfig{
		ServerType:  serverType,
		Image:       image,
		Location:    location,
		SSHKey:      sshKey,
		Network:     network,
		ClusterRole: k3s.Master,
	})

	successMessage := "Successfully created server with IP:" + server.Server.PublicNet.IPv4.IP.String()

	spinner.StopMessage(successMessage)
	spinner.Stop()
}

type CreateClusterConfig struct {
	ServerType *hcloud.ServerType
	Image      *hcloud.Image
	Location   *hcloud.Location
	SSHKey     *hcloud.SSHKey
}

func getAllImages(client *hcloud.Client) []*hcloud.Image {
	allImages, _ := client.Image.All(context.Background())

	var filteredImages []*hcloud.Image

	for _, image := range allImages {
		if strings.HasPrefix(image.Name, "ubuntu") {
			filteredImages = append(filteredImages, image)
		}
	}

	return filteredImages
}

func getAllLocations(client *hcloud.Client) []*hcloud.Location {
	allLocations, _ := client.Location.All(context.Background())

	return allLocations
}

func getAllServerTypes(client *hcloud.Client) []*hcloud.ServerType {
	allServerTypes, _ := client.ServerType.All(context.Background())
	return allServerTypes
}

func getAllSSHKeys(client *hcloud.Client) []*hcloud.SSHKey {
	allSSHKeys, _ := client.SSHKey.All(context.Background())

	return allSSHKeys
}

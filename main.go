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

	spinner.Message("Loading Locations")
	spinner.Start()
	locations, _ := getAllLocations(client)
	spinner.Stop()
	location, err := cli.SelectLocation(locations)
	if err != nil {
		fmt.Println("Aborting...")
		return
	}

	spinner.Message("Loading Server Types")
	spinner.Start()
	types := getServerTypes(client, location)
	spinner.Stop()
	serverType, err := cli.SelectServerType(types)
	if err != nil {
		fmt.Println("Aborting...")
		return
	}

	spinner.Message("Loading Keys")
	spinner.Start()
	sshKeys := getAllSSHKeys(client)
	spinner.Stop()
	sshKey, err := cli.SelectSSHKey(sshKeys)
	if err != nil {
		fmt.Println("Aborting...")
		return
	}

	_, err = cli.SelectNodeAmount()
	if err != nil {
		fmt.Println("Aborting...")
		return
	}

	spinner.Message("Loading Images")
	spinner.Start()
	images := getAllImages(client)
	spinner.Stop()
	image, err := cli.SelectImage(images)
	if err != nil {
		fmt.Println("Aborting...")
		return
	}

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

func getAllLocations(client *hcloud.Client) ([]*hcloud.Location, error) {
	return client.Location.All(context.Background())
}

func getServerTypes(client *hcloud.Client, selectedLocation *hcloud.Location) []*hcloud.ServerType {
	allServerTypes, _ := client.ServerType.All(context.Background())

	var filteredServerTypes []*hcloud.ServerType

	for _, serverType := range allServerTypes {
		for _, pricing := range serverType.Pricings {
			location := pricing.Location
			if location.Name == selectedLocation.Name {
				filteredServerTypes = append(filteredServerTypes, serverType)
				break
			}
		}
	}

	return filteredServerTypes
}

func getAllSSHKeys(client *hcloud.Client) []*hcloud.SSHKey {
	allSSHKeys, _ := client.SSHKey.All(context.Background())

	return allSSHKeys
}

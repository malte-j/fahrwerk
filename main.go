package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SkyGuardian42/fahrwerk/cli"
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
	createCluster(client, CreateClusterConfig{
		Image:      image,
		Location:   location,
		SSHKey:     sshKey,
		ServerType: serverType,
	})
	spinner.StopMessage("Successfully created the cluster!")
	spinner.Stop()
}

type CreateClusterConfig struct {
	ServerType *hcloud.ServerType
	Image      *hcloud.Image
	Location   *hcloud.Location
	SSHKey     *hcloud.SSHKey
}

func createCluster(client *hcloud.Client, config CreateClusterConfig) {
	if "no" == "no" {
		return
	}
	createdServer, _, err := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       "main-1",
		Image:      config.Image,
		ServerType: config.ServerType,
		Location:   config.Location,
		SSHKeys:    []*hcloud.SSHKey{config.SSHKey},
		Labels:     map[string]string{"kubernetes": "true"},
	})
	if err != nil {
		log.Fatalf("error retrieving server: %s\n", err)
	}

	fmt.Println("created server with ip:", createdServer.Server.PublicNet.IPv4.IP)
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

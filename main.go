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
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[14],
		Suffix: " ",
		Colors: []string{"fgGreen"},
		StopColors:      []string{"fgGreen"},
	}
	spinner, _ := yacspin.New(cfg)


	hcloudToken, _ := cli.GetHcloudToken()


	spinner.Message("Connecting to HCloud")
	spinner.Start()
	client := hcloud.NewClient(hcloud.WithToken(hcloudToken))
	spinner.Stop()

	spinner.Message("Loading Keys")
	spinner.Start()
	sshKeys := getAllSSHKeys(client)
	spinner.Stop()
	cli.SelectSSHKey(sshKeys)

	cli.SelectNodeAmount()

	spinner.Message("Loading Locations")
	spinner.Start()
	locations := getAllLocations(client)
	spinner.Stop()
	cli.SelectLocation(locations)


	spinner.Message("Loading Images")
	spinner.Start()
	images := getAllImages(client)
	spinner.Stop()
	cli.SelectImage(images)
}

func createServer(client *hcloud.Client) {
	cx11ServerType, _, _ := client.ServerType.Get(context.Background(), "cx11")
	ubuntuImage, _, _ := client.Image.GetByName(context.Background(), "ubuntu-20.04")
	nurnbergLocation, _, _ := client.Location.GetByName(context.Background(), "nbg1")
	sshKey, _, _ := client.SSHKey.GetByName(context.Background(), "malte@maltspad")

	createdServer, _, err := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       "main-1",
		Image:      ubuntuImage,
		ServerType: cx11ServerType,
		Location:   nurnbergLocation,
		SSHKeys:    []*hcloud.SSHKey{sshKey},
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

func getAllSSHKeys(client *hcloud.Client) []*hcloud.SSHKey {
	allSSHKeys, _ := client.SSHKey.All(context.Background())

	return allSSHKeys
}

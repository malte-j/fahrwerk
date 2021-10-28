package cli

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/manifoldco/promptui"
)

func SelectSSHKey(existingSSHKeys []*hcloud.SSHKey) *hcloud.SSHKey {
	templates := &promptui.SelectTemplates{
		Active:   "► {{ .Name }}",
		Inactive: "  {{ .Name }}",
		Selected: `{{ "SSH Key:" | faint }} {{ .Name }}`,
	}

	prompt := promptui.Select{
		Label:     `SSH Key`,
		Items:     existingSSHKeys,
		Templates: templates,
	}

	i, _, _ := prompt.Run()
	return existingSSHKeys[i]
}

func SelectLocation(locations []*hcloud.Location) *hcloud.Location {
	templates := &promptui.SelectTemplates{
		Active:   "► {{ .City }} ({{ .Name | faint }})",
		Inactive: "  {{ .City }} ({{ .Name | faint }})",
		Selected: `{{ "Location:" | faint }} {{ .City }}`,
	}

	prompt := promptui.Select{
		Label:     `Which Location?`,
		Items:     locations,
		Templates: templates,
	}

	i, _, _ := prompt.Run()
	return locations[i]
}

func SelectNodeAmount() int {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("must be numeric")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Success: `{{"Worker Nodes: " | faint}}`,
	}

	prompt := promptui.Prompt{
		Label:     `Amount of Worker Nodes`,
		Default:   "2",
		Validate:  validate,
		Templates: templates,
		Pointer: promptui.PipeCursor,
	}

	result, _ := prompt.Run()
	resultInt, _ := strconv.Atoi(result)

	return resultInt
}

func SelectImage(images []*hcloud.Image) (selectedImage *hcloud.Image, err error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "► {{ .Description }} ({{ .Name | faint }})",
		Inactive: "  {{ .Description }} ({{ .Name | faint }})",
		Selected: `{{ "Image:" | faint}} {{ .Description }}`,
	}

	searcher := func(input string, index int) bool {
		image := images[index]
		name := strings.Replace(strings.ToLower(image.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Images",
		Items:     images,
		Templates: templates,
		Size:      4,
		
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()
	return images[i], err
}

func GetHcloudToken() (result string, err error) {
	envToken, envTokenExists := os.LookupEnv("HCLOUD_TOKEN")
	if(envTokenExists) {
		return envToken, nil
	}

	prompt := promptui.Prompt{
		Label:       "HCloud Token",
		HideEntered: false,
		Mask:        '*',
	}

	return prompt.Run()
}

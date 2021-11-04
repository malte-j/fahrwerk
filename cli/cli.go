package cli

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/manifoldco/promptui"
)

func SelectSSHKey(existingSSHKeys []*hcloud.SSHKey) (*hcloud.SSHKey, error) {
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

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return existingSSHKeys[i], nil
}

func SelectLocation(locations []*hcloud.Location) (*hcloud.Location, error) {
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

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return locations[i], nil
}

func SelectServerType(serverTypes []*hcloud.ServerType) (*hcloud.ServerType, error) {
	templates := &promptui.SelectTemplates{
		Active:   "► {{ .Description }} ({{ .Name | faint }})",
		Inactive: "  {{ .Description }} ({{ .Name | faint }})",
		Selected: `{{ "Server Type:" | faint }} {{ .Description }}`,
	}

	prompt := promptui.Select{
		Label:     `Select a Server Type`,
		Items:     serverTypes,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err 
	}
	return serverTypes[i], nil
}

func SelectNodeAmount() (int, error) {
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
		Pointer:   promptui.PipeCursor,
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	resultInt, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}

	return resultInt, nil
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

		Searcher: searcher,
	}

	i, _, err := prompt.Run()
	return images[i], err
}

func GetHcloudToken() (result string, err error) {
	envToken, envTokenExists := os.LookupEnv("HCLOUD_TOKEN")
	if envTokenExists {
		return envToken, nil
	}

	prompt := promptui.Prompt{
		Label:       "HCloud Token",
		HideEntered: false,
		Mask:        '*',
	}

	return prompt.Run()
}

func ConfirmClusterCreation() bool {
	prompt := promptui.Prompt{
		Label:     "Do you really want to create the cluster",
		IsConfirm: true,
	}

	confirmationResult, _ := prompt.Run()
	confirmationResultLowercase := strings.ToLower(confirmationResult)

	switch confirmationResultLowercase {
	case "y":
		return true
	case "yes":
		return true
	}

	return false
}

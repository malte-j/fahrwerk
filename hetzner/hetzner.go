package hetzner

import (
	"context"
	"net"

	"github.com/SkyGuardian42/fahrwerk/k3s"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func CreateNetwork(client *hcloud.Client) (*hcloud.Network, *hcloud.Response, error) {
	return client.Network.Create(context.Background(), hcloud.NetworkCreateOpts{
		Name: "cluster-network",
		IPRange: &net.IPNet{
			IP:   net.ParseIP("10.0.0.0"),
			Mask: net.CIDRMask(16, 32),
		},
		Subnets: []hcloud.NetworkSubnet{
			{
				Type:        hcloud.NetworkSubnetTypeCloud,
				NetworkZone: hcloud.NetworkZoneEUCentral,
				IPRange: &net.IPNet{
					IP:   net.ParseIP("10.0.0.0"),
					Mask: net.CIDRMask(16, 32),
				},
			},
		},
		Labels: map[string]string{
			"cluster":      "",
			"cluster-role": "network",
		},
	})
}

type CreateServerConfig struct {
	ServerType  *hcloud.ServerType
	Image       *hcloud.Image
	Location    *hcloud.Location
	SSHKey      *hcloud.SSHKey
	Network     *hcloud.Network
	ClusterRole k3s.ClusterRole
}

func CreateServer(client *hcloud.Client, config CreateServerConfig) (hcloud.ServerCreateResult, *hcloud.Response, error) {
	return client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       "main-1",
		Image:      config.Image,
		ServerType: config.ServerType,
		Location:   config.Location,
		Networks: []*hcloud.Network{config.Network},
		SSHKeys:    []*hcloud.SSHKey{config.SSHKey},
		Labels: map[string]string{
			"cluster":      "",
			"cluster-role": string(config.ClusterRole),
		},
	})
}

package k3s

import (
	"bytes"
	"math/rand"
	"strings"
	"text/template"
	"time"
)

type ClusterRole string

const (
	Master ClusterRole = "Master"
	Worker ClusterRole = "Worker"
)

type MasterScriptConfig struct {
	K3sToken   string
}

func GenerateMasterScript(config MasterScriptConfig) string {
	t, err := template.New("masterConfig").Parse(`
		curl -sfL https://get.k3s.io | K3S_TOKEN="{{ .K3sToken }}" INSTALL_K3S_EXEC="server \
			--disable-cloud-controller \
			--disable servicelb \
			--disable traefik \
			--disable local-storage \
			--disable metrics-server \
			--write-kubeconfig-mode=644 \
			--node-name="$(hostname -f)" \
			--cluster-cidr=10.244.0.0/16 \
			--etcd-expose-metrics=true \
			--kube-controller-manager-arg="address=0.0.0.0" \
			--kube-controller-manager-arg="bind-address=0.0.0.0" \
			--kube-proxy-arg="metrics-bind-address=0.0.0.0" \
			--kube-scheduler-arg="address=0.0.0.0" \
			--kube-scheduler-arg="bind-address=0.0.0.0" \
			--kubelet-arg="cloud-provider=external" \
			--advertise-address=$(hostname -I | awk '{print $2}') \
			--node-ip=$(hostname -I | awk '{print $2}') \
			--node-external-ip=$(hostname -I | awk '{print $1}') \
			--flannel-iface=#{flannel_interface} \
			--cluster-init" sh -
	`)

	// 			#{taint} \

	if err != nil {
		panic(err)
	}

	var masterScriptBuffer bytes.Buffer
	err = t.Execute(&masterScriptBuffer, config)
	if err != nil {
		panic(err)
	}

	return masterScriptBuffer.String()
}

func GenerateK3sToken() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"
	var src = rand.NewSource(time.Now().UnixNano())

	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	n := 20
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

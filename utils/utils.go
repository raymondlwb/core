package utils

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	engineapi "github.com/docker/engine-api/client"
	"github.com/docker/go-connections/tlsconfig"
	"gitlab.ricebook.net/platform/core/types"
	"golang.org/x/net/context"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortenLength = 7
)

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	r := make([]byte, n)
	for i := 0; i < n; i++ {
		r[i] = letters[rand.Intn(len(letters))]
	}
	return string(r)
}

func TruncateID(id string) string {
	if len(id) > shortenLength {
		return id[:shortenLength]
	}
	return id
}

func Tail(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func MakeDockerClient(endpoint string, config types.Config) (*engineapi.Client, error) {
	if !strings.HasPrefix(endpoint, "tcp://") {
		endpoint = "tcp://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}

	var cli *http.Client
	// if no cert path is set
	// then just use normal http client without tls
	if config.Docker.CertPath != "" {
		dockerCertPath := filepath.Join(config.Docker.CertPath, host)
		options := tlsconfig.Options{
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: false,
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		cli = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
		}
	}

	return engineapi.NewClient(endpoint, config.Docker.APIVersion, cli, nil)
}

func GetGitRepoName(url string) (string, error) {
	if !strings.HasPrefix(url, "git@") || !strings.HasSuffix(url, ".git") {
		return "", fmt.Errorf("Bad git url format %q", url)
	}

	x := strings.SplitN(url, ":", 2)
	if len(x) != 2 {
		return "", fmt.Errorf("Bad git url format %q", url)
	}

	y := strings.SplitN(x[1], "/", 2)
	if len(y) != 2 {
		return "", fmt.Errorf("Bad git url format %q", url)
	}
	return strings.TrimSuffix(y[1], ".git"), nil
}

func GetVersion(image string) string {
	if !strings.Contains(image, ":") {
		return "unknown"
	}

	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return "unknown"
	}

	return parts[1]
}

// Bind a docker engine client to context
func ToDockerContext(client *engineapi.Client) context.Context {
	return context.WithValue(context.Background(), "engine", client)
}

// Get a docker engine client from a context
func FromDockerContext(ctx context.Context) (*engineapi.Client, bool) {
	client, ok := ctx.Value("engine").(*engineapi.Client)
	return client, ok
}
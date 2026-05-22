package docker

import (
	"strings"
)

// ServiceMap maps a compose service name to its running containers.
type ServiceMap map[string][]ContainerInfo

// BuildServiceMap groups containers by their compose service label.
// Containers without the label are grouped under the empty string key.
func BuildServiceMap(containers []ContainerInfo) ServiceMap {
	sm := make(ServiceMap)
	for _, c := range containers {
		service := c.Labels["com.docker.compose.service"]
		sm[service] = append(sm[service], c)
	}
	return sm
}

// ServiceNames returns a sorted-stable slice of service names present in the map.
func (sm ServiceMap) ServiceNames() []string {
	names := make([]string, 0, len(sm))
	for k := range sm {
		if k != "" {
			names = append(names, k)
		}
	}
	return names
}

// HasService reports whether a compose service name exists in the map.
func (sm ServiceMap) HasService(name string) bool {
	_, ok := sm[name]
	return ok
}

// ImageForService returns the image of the first container for a service.
// Returns an empty string if the service is not found.
func (sm ServiceMap) ImageForService(name string) string {
	containers, ok := sm[name]
	if !ok || len(containers) == 0 {
		return ""
	}
	return containers[0].Image
}

// NormalizeImageRef strips the implicit 'docker.io/library/' prefix so that
// 'nginx:latest' and 'docker.io/library/nginx:latest' compare as equal.
func NormalizeImageRef(image string) string {
	image = strings.TrimPrefix(image, "docker.io/library/")
	image = strings.TrimPrefix(image, "docker.io/")
	if !strings.Contains(image, ":") {
		image += ":latest"
	}
	return image
}

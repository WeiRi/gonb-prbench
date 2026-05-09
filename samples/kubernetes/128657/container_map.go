package containermap

// containerMap is a stripped reproduction of pkg/kubelet/cm/containermap/container_map.go.
// BUG (pre-PR #128657): map operations without a mutex => concurrent map access.
type containerMap map[string]struct {
	podUID        string
	containerName string
}

// ContainerMap is what NewContainerMap returns. The PR introduced a sync.Mutex;
// the BUG state intentionally omits it to surface the race.
type ContainerMap struct {
	cm containerMap
	// no sync.Mutex: BUG state
}

func NewContainerMap() *ContainerMap {
	return &ContainerMap{cm: make(containerMap)}
}

func (c *ContainerMap) Add(podUID, containerName, containerID string) {
	c.cm[containerID] = struct {
		podUID        string
		containerName string
	}{podUID, containerName}
}

func (c *ContainerMap) RemoveByContainerID(containerID string) {
	delete(c.cm, containerID)
}

func (c *ContainerMap) GetContainerRef(containerID string) (string, string, bool) {
	v, ok := c.cm[containerID]
	if !ok {
		return "", "", false
	}
	return v.podUID, v.containerName, true
}

func (c *ContainerMap) GetContainerID(podUID, containerName string) (string, bool) {
	for cid, v := range c.cm {
		if v.podUID == podUID && v.containerName == containerName {
			return cid, true
		}
	}
	return "", false
}

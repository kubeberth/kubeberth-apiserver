package berth

type AttachedArchive struct {
	Name string `json:"name"`
}

type AttachedCloudInit struct {
	Name string `json:"name"`
}

type AttachedDisk struct {
	Name string `json:"name"`
}

type AttachedSource struct {
	Archive *AttachedArchive `json:"archive,omitempty"`
	Disk    *AttachedDisk    `json:"disk,omitempty"`
}

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type VersionList struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	RTUVersion  string `json:"rtuVersion"`
	ReleaseDate string `json:"releaseDate"`
}

func GetAllVersions() (*VersionList, error) {
	vl := &VersionList{}

	versionsPath, _ := filepath.Abs("../common/versions.json")
	content, err := ioutil.ReadFile(versionsPath)
	if err != nil {
		return vl, fmt.Errorf("ERROR: unable to read version.json: %s", err)
	}

	err = json.Unmarshal(content, vl)
	if err != nil {
		return vl, fmt.Errorf("ERROR: unable to parse version.json content: %s", err)
	}

	return vl, nil
}

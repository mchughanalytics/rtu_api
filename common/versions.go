package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type VersionList struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	RTUVersion string `json:"rtuVersion"`
}

func GetAllVersions() (*VersionList, error) {
	vl := &VersionList{}

	content, err := ioutil.ReadFile("versions.json")
	if err != nil {
		return vl, fmt.Errorf("ERROR: unable to read version.json: %s", err)
	}

	err = json.Unmarshal(content, vl)
	if err != nil {
		return vl, fmt.Errorf("ERROR: unable to parse version.json content: %s", err)
	}

	return vl, nil
}

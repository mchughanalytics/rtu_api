package rtuapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type VersionList struct {
	Versions []*Version `json:"versions"`
}

type Version struct {
	RTUVersion           string   `json:"rtuVersion"`
	ReleaseDate          string   `json:"releaseDate"`
	DownloadLink         string   `json:"downloadLink"`
	Platform             string   `json:"platform"`
	CompatibleRMVersions []string `json:"compatibleRMVersions"`
}

func (v *Version) GetMajor() int {
	_split := strings.Split(v.RTUVersion, ".")
	_val, _ := strconv.Atoi(_split[0])
	return _val
}

func (v *Version) GetMinor() int {
	_split := strings.Split(v.RTUVersion, ".")
	_val, _ := strconv.Atoi(_split[1])
	return _val
}

func (v *Version) GetHotFix() int {
	_split := strings.Split(v.RTUVersion, ".")
	_val, _ := strconv.Atoi(_split[2])
	return _val
}

func GetAllVersions() (*VersionList, error) {
	vl := &VersionList{}

	dir, err := os.Getwd()
	if err != nil {
		return vl, fmt.Errorf("ERROR: unable to read version.json: %s", err)
	}

	versionsPath := fmt.Sprintf("%s/versions.json", dir)
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

func (vl *VersionList) GetLatest() *Version {

	latest := &Version{}

	for index, v := range vl.Versions {
		if index == 0 {
			latest = v
		} else {
			if v.GetMajor() > latest.GetMajor() {
				latest = v
			}
			if v.GetMajor() == latest.GetMajor() && v.GetMinor() > latest.GetMinor() {
				latest = v
			}
			if v.GetMajor() == latest.GetMajor() && v.GetMinor() == latest.GetMinor() && v.GetHotFix() > latest.GetHotFix() {
				latest = v
			}
		}
	}

	return latest
}

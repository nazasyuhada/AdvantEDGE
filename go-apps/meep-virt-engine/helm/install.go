/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helm

import (
	"errors"
	"os/exec"
	"strings"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
)

func installCharts(charts []Chart) error {
	err := ensureReleases(charts)
	if err != nil {
		return err
	}

	for _, chart := range charts {
		err := install(chart)
		if err != nil {
			log.Info("Cleaning installed releases")
			cleanReleases(charts)
			return err
		}
	}

	return nil
}

func ensureReleases(charts []Chart) error {
	// ensure that releases do not already exist
	releases, _ := GetReleasesName()
	for _, c := range charts {
		for _, r := range releases {
			if c.ReleaseName == r.Name {
				err := errors.New("Release [" + c.ReleaseName + "] already exists")
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

func install(chart Chart) error {
	log.Debug("Installing chart: " + chart.ReleaseName)
	var cmd *exec.Cmd
	if strings.Trim(chart.ValuesFile, " ") == "" {
		cmd = exec.Command("helm", "install",
			"--name", chart.ReleaseName,
			"--namespace", chart.Namespace,
			"--set", "nameOverride="+chart.Name,
			"--set", "fullnameOverride="+chart.Name,
			chart.Location, "--replace")
	} else {
		cmd = exec.Command("helm", "install",
			"--name", chart.ReleaseName,
			"--namespace", chart.Namespace,
			"--set", "nameOverride="+chart.Name,
			"--set", "fullnameOverride="+chart.Name,
			"-f", chart.ValuesFile,
			chart.Location, "--replace")
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to install Release [" + chart.ReleaseName + "] at " + chart.Location)
		log.Error("Error(", err.Error(), "): ", string(out))
		return err
	}
	return nil
}

func cleanReleases(charts []Chart) {
	var toClean []Chart
	var cnt int
	releases, _ := GetReleasesName()

	for _, c := range charts {
		for _, r := range releases {
			if c.ReleaseName == r.Name {
				toClean = append(toClean, c)
				cnt++
			}
		}
	}

	if cnt > 0 {
		_ = deleteReleases(toClean)
	}
}

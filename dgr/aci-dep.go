package main

import (
	"github.com/n0rad/go-erlog/data"
	"github.com/n0rad/go-erlog/errs"
	"strconv"
	"github.com/blablacar/dgr/dgr/common"
	"github.com/n0rad/go-erlog/logs"
)

func GetDependencyDgrVersion(acName common.ACFullname) (int, error) {
	depFields := data.WithField("dependency", acName.String())

	im, err := Home.Rkt.GetManifest(acName.String())
	if err != nil {
		return 0, errs.WithEF(err, depFields, "Dependency not found")
	}

	version, ok := im.Annotations.Get(common.ManifestDrgVersion)
	var val int
	if ok {
		val, err = strconv.Atoi(version)
		if err != nil {
			return 0, errs.WithEF(err, depFields.WithField("version", version), "Failed to parse "+common.ManifestDrgVersion+" from manifest")
		}
	}
	return val, nil
}

func CheckLatestVersion(deps []common.ACFullname, warnText string) {
	for _, dep := range deps {
		if dep.Version() == "" {
			continue
		}
		version, _ := dep.LatestVersion()
		if version != "" && common.Version(dep.Version()).LessThan(common.Version(version)) {
			logs.WithField("newer", dep.Name()+":"+version).
				WithField("current", dep.String()).
				Warn("Newer " + warnText + " version")
		}
	}
}

func (aci *Aci) checkCompatibilityVersions() {
	defer aci.checkWg.Done()
	for _, dep := range aci.manifest.Aci.Dependencies {
		depFields := aci.fields.WithField("dependency", dep.String())

		logs.WithF(aci.fields).WithField("dependency", dep.String()).Info("Fetching dependency")
		Home.Rkt.Fetch(dep.String())
		version, err := GetDependencyDgrVersion(dep)
		if err != nil {
			logs.WithEF(err, depFields).Error("Failed to check compatibility version of dependency")
		} else {
			if version < 55 {
				logs.WithF(aci.fields).
					WithField("dependency", dep).
					WithField("require", ">=55").
					Error("dependency was not build with a compatible version of dgr")
			}
		}

	}
}
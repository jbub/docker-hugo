// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	dockerImage      = "jbub/docker-hugo"
	dockerBaseImage  = "alpine:latest"
	dockerMaintainer = "Juraj Bubniak <juraj.bubniak@gmail.com>"
)

type versionInfo struct {
	Name       string
	BaseImage  string
	Version    string
	Maintainer string
}

func (v versionInfo) tag(latest bool) string {
	if latest {
		return fmt.Sprintf("%v:latest", dockerImage)
	}
	return fmt.Sprintf("%v:%v", dockerImage, v.Name)
}

var versions = []versionInfo{
	{Name: "0.28", Version: "0.28", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.29", Version: "0.29", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.30", Version: "0.30.2", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.31", Version: "0.31.1", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.32", Version: "0.32.4", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.33", Version: "0.33", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.34", Version: "0.34", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.35", Version: "0.35", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "0.36", Version: "0.36.1", BaseImage: dockerBaseImage, Maintainer: dockerMaintainer},
}

var dockerfileTmplString = `FROM {{ .BaseImage }}
MAINTAINER {{ .Maintainer }}

ENV HUGO_VERSION={{ .Version }}

RUN apk --no-cache add wget ca-certificates && \
  cd /tmp && \
  wget https://github.com/spf13/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz && \
  tar xzf hugo_${HUGO_VERSION}_Linux-64bit.tar.gz && \
  rm hugo_${HUGO_VERSION}_Linux-64bit.tar.gz && \
  mv hugo /usr/bin/hugo && \
  apk del wget ca-certificates

ENTRYPOINT ["/usr/bin/hugo"]`

var dockerfileTmpl = template.Must(template.New("dockerfile").Parse(dockerfileTmplString))

var docker = sh.RunCmd("docker")

// Generate generates dockerfiles for all versions.
func Generate() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get wd: %v", err)
	}

	for _, info := range versions {
		dir := filepath.Join(wd, info.Name)
		if err := ensureDir(dir); err != nil {
			return err
		}

		if err := genDockerfile(dir, info); err != nil {
			return err
		}
	}

	latest := versions[len(versions)-1]
	if err := genDockerfile(wd, latest); err != nil {
		return err
	}

	return nil
}

// Docker builds and runs docker images for all versions.
func Docker() error {
	mg.Deps(Generate)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get wd: %v", err)
	}

	for _, info := range versions {
		dir := filepath.Join(wd, info.Name)
		if err := buildAndRunDocker(dir, info, false); err != nil {
			return fmt.Errorf("could not build/run docker: %v", err)
		}
	}

	latest := versions[len(versions)-1]
	if err := buildAndRunDocker(wd, latest, true); err != nil {
		return fmt.Errorf("could not build/run docker: %v", err)
	}

	return nil
}

// Push pushes built docker images do docker hub.
func Push() error {
	mg.Deps(Docker)

	for _, info := range versions {
		if err := pushDocker(info, false); err != nil {
			return fmt.Errorf("could not push docker image: %v", err)
		}
	}

	latest := versions[len(versions)-1]
	if err := pushDocker(latest, true); err != nil {
		return fmt.Errorf("could not push docker image: %v", err)
	}

	return nil
}

func pushDocker(info versionInfo, latest bool) error {
	return docker("push", info.tag(latest))
}

func buildAndRunDocker(dir string, info versionInfo, latest bool) error {
	if err := docker("build", "-t", info.tag(latest), dir); err != nil {
		return err
	}
	if err := docker("run", "--interactive", "--tty", "--rm", info.tag(latest), "version"); err != nil {
		return err
	}
	return nil
}

func genDockerfile(dir string, info versionInfo) error {
	fpath := filepath.Join(dir, "Dockerfile")
	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer fp.Close()

	if err := dockerfileTmpl.Execute(fp, info); err != nil {
		return fmt.Errorf("could not execute template: %v", err)
	}
	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not stat path: %v", err)
		}

		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return fmt.Errorf("could not create dir: %v", err)
		}
	}
	return nil
}

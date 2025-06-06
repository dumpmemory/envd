// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/types"
)

func Base(image string, dev bool) error {
	g := DefaultGraph.(*generalGraph)

	if image != "" {
		g.Image = image
	}
	g.Dev = dev
	return nil
}

func Python(version string) error {
	if strings.HasPrefix(version, "2") {
		logrus.Debugf("envd doesn't support Python2: %s", version)
		return errors.New("envd doesn't support this Python version")
	}
	g := DefaultGraph.(*generalGraph)

	g.Languages = append(g.Languages, ir.Language{
		Name:    "python",
		Version: &version,
	})
	return nil
}

func Conda(mamba bool) {
	g := DefaultGraph.(*generalGraph)

	g.CondaConfig = &ir.CondaConfig{
		UseMicroMamba: mamba,
	}
}

func Pixi(usePixiMirror bool, pypiIndex string) {
	g := DefaultGraph.(*generalGraph)

	g.PixiConfig = &ir.PixiConfig{
		UsePixiMirror: usePixiMirror,
		PyPIIndex:     nil,
	}
	if len(pypiIndex) != 0 {
		g.PixiConfig.PyPIIndex = &pypiIndex
	}
}

func UV(pythonVersion string) {
	g := DefaultGraph.(*generalGraph)

	g.UVConfig = &ir.UVConfig{
		PythonVersion: pythonVersion,
	}
}

func RLang() {
	g := DefaultGraph.(*generalGraph)

	g.Languages = append(g.Languages, ir.Language{
		Name: "r",
	})
}

func Julia() {
	g := DefaultGraph.(*generalGraph)

	g.Languages = append(g.Languages, ir.Language{
		Name: "julia",
	})
}

func PyPIPackage(deps []string, requirementsFile string, wheels []string) error {
	g := DefaultGraph.(*generalGraph)

	if len(deps) > 0 {
		g.PyPIPackages = append(g.PyPIPackages, deps)
	}
	g.PythonWheels = append(g.PythonWheels, wheels...)

	if requirementsFile != "" {
		g.RequirementsFile = &requirementsFile
	}

	return nil
}

func RPackage(deps []string) error {

	if len(deps) == 0 {
		return errors.New("Can not install empty R package")
	}

	g := DefaultGraph.(*generalGraph)

	g.RPackages = append(g.RPackages, deps)

	return nil
}

func JuliaPackage(deps []string) error {

	if len(deps) == 0 {
		return errors.New("Can not install empty Julia package")
	}

	g := DefaultGraph.(*generalGraph)

	g.JuliaPackages = append(g.JuliaPackages, deps)

	return nil
}

func SystemPackage(deps []string) {
	g := DefaultGraph.(*generalGraph)

	g.SystemPackages = append(g.SystemPackages, deps...)
}

func ShmSize(shmSize int) {
	g := DefaultGraph.(*generalGraph)

	g.ShmSize = shmSize
}

func GPU(numGPUs int) {
	g := DefaultGraph.(*generalGraph)

	g.NumGPUs = numGPUs
}

func CUDA(version, cudnn string) {
	g := DefaultGraph.(*generalGraph)

	g.CUDA = &version
	if len(cudnn) > 0 {
		g.CUDNN = cudnn
	}
}

func VSCodePlugins(plugins []string) error {
	g := DefaultGraph.(*generalGraph)

	for _, p := range plugins {
		plugin, err := vscode.ParsePlugin(p)
		if err != nil {
			return err
		}
		g.VSCodePlugins = append(g.VSCodePlugins, *plugin)
	}
	return nil
}

// UbuntuAPT updates the Ubuntu apt source.list in the image.
func UbuntuAPT(source string) error {
	if source == "" {
		return errors.New("source is required")
	}
	g := DefaultGraph.(*generalGraph)

	g.UbuntuAPTSource = &source
	return nil
}

func PyPIIndex(url, extraURL string, trust bool) error {
	if url == "" {
		return errors.New("url is required")
	}
	g := DefaultGraph.(*generalGraph)

	g.PyPIIndexURL = &url
	if len(extraURL) > 0 {
		g.PyPIExtraIndexURL = &extraURL
	}
	g.PyPITrust = trust
	return nil
}

func CRANMirror(url string) error {
	g := DefaultGraph.(*generalGraph)

	g.CRANMirrorURL = &url
	return nil
}

func JuliaPackageServer(url string) error {
	g := DefaultGraph.(*generalGraph)

	g.JuliaPackageServer = &url
	return nil
}

func Shell(shell string) error {
	g := DefaultGraph.(*generalGraph)

	g.Shell = strings.ToLower(shell)
	return nil
}

func Jupyter(pwd string, port int64) error {
	g := DefaultGraph.(*generalGraph)

	g.JupyterConfig = &ir.JupyterConfig{
		Token: pwd,
		Port:  port,
	}
	return nil
}

func RStudioServer() error {
	g := DefaultGraph.(*generalGraph)

	g.RStudioServerConfig = &ir.RStudioServerConfig{}
	return nil
}

func Run(commands []string, mount bool) error {
	g := DefaultGraph.(*generalGraph)

	g.Exec = append(g.Exec, ir.RunBuildCommand{
		Commands:  commands,
		MountHost: mount,
	})
	return nil
}

func Git(name, email, editor string) error {
	g := DefaultGraph.(*generalGraph)

	g.GitConfig = &ir.GitConfig{
		Name:   name,
		Email:  email,
		Editor: editor,
	}
	return nil
}

func CondaChannel(channel string) error {
	g := DefaultGraph.(*generalGraph)

	if g.CondaConfig == nil {
		return errors.New("cannot config conda when conda is not installed")
	}
	g.CondaConfig.CondaChannel = &channel
	return nil
}

func CondaPackage(deps []string, channel []string, envFile string) error {
	g := DefaultGraph.(*generalGraph)

	if g.CondaConfig == nil {
		return errors.New("cannot install conda packages when conda is not installed")
	}
	g.CondaConfig.CondaPackages = append(
		g.CondaConfig.CondaPackages, deps...)

	g.CondaConfig.CondaEnvFileName = envFile

	if len(channel) != 0 {
		g.CondaConfig.AdditionalChannels = append(
			g.CondaConfig.AdditionalChannels, channel...)
	}
	return nil
}

func Copy(src, dest, image string) {
	g := DefaultGraph.(*generalGraph)

	g.Copy = append(g.Copy, ir.CopyInfo{
		Source:      src,
		Destination: dest,
		Image:       image,
	})
}

func Mount(src, dest string) {
	g := DefaultGraph.(*generalGraph)

	g.Mount = append(g.Mount, ir.MountInfo{
		Source:      src,
		Destination: dest,
	})
}

func HTTP(url, checksum, filename string) error {
	g := DefaultGraph.(*generalGraph)

	info := ir.HTTPInfo{
		URL:      url,
		Filename: filename,
	}
	if len(checksum) > 0 {
		d, err := digest.Parse(checksum)
		if err != nil {
			return err
		}
		info.Checksum = d
	}
	g.HTTP = append(g.HTTP, info)
	return nil
}

func Entrypoint(args []string) {
	g := DefaultGraph.(*generalGraph)

	g.Entrypoint = append(g.Entrypoint, args...)
}

func RuntimeCommands(commands map[string]string) {
	g := DefaultGraph.(*generalGraph)

	for k, v := range commands {
		g.RuntimeCommands[k] = v
	}
}

func RuntimeDaemon(commands [][]string) {
	g := DefaultGraph.(*generalGraph)

	g.RuntimeDaemon = append(g.RuntimeDaemon, commands...)
}

func RuntimeExpose(envdPort, hostPort int, serviceName string, listeningAddr string) error {
	g := DefaultGraph.(*generalGraph)

	g.RuntimeExpose = append(g.RuntimeExpose, ir.ExposeItem{
		EnvdPort:      envdPort,
		HostPort:      hostPort,
		ServiceName:   serviceName,
		ListeningAddr: listeningAddr,
	})
	return nil
}

func RuntimeEnviron(env map[string]string, path []string) {
	g := DefaultGraph.(*generalGraph)

	for k, v := range env {
		g.RuntimeEnviron[k] = v
	}
	g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, path...)
}

func RuntimeInitScript(commands []string) {
	g := DefaultGraph.(*generalGraph)

	g.RuntimeInitScript = append(g.RuntimeInitScript, commands)
}

func Repo(url, description string) {
	g := DefaultGraph.(*generalGraph)

	g.Repo = types.RepoInfo{
		Description: description,
		URL:         url,
	}
}

func Owner(uid, gid int) {
	g := DefaultGraph.(*generalGraph)
	g.uid = uid
	g.gid = gid
}

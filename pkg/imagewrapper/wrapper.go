package imagewrapper

import (
	"errors"
	"strings"

	"github.com/JitenPalaparthi/imagelinter/pkg/cmdhelper"
)

type Wrapper struct {
	Image     string
	Container string
	CmdHelper *cmdhelper.CmdHelper
}

func New(image, contaier string, cmdhelper *cmdhelper.CmdHelper) (*Wrapper, error) {
	if cmdhelper == nil {
		return nil, errors.New("command helper cannot be nil")
	}
	if image == "" || contaier == "" {
		return nil, errors.New("image or Container parameters cannot be empty")
	}
	return &Wrapper{Image: image, Container: contaier, CmdHelper: cmdhelper}, nil
}

func (w *Wrapper) PullImage() (string, error) {
	result, err := w.CmdHelper.CliRunner("docker", nil, []string{"pull", w.Image}...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (w *Wrapper) CreateContainer() (string, error) {
	result, err := w.CmdHelper.CliRunner("docker", nil, []string{"run", "-d", "--name", w.Container, w.Image}...)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (w *Wrapper) RunCommand(args ...string) (string, error) {
	result, err := w.CmdHelper.CliRunner("docker", nil, args...)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (w *Wrapper) IsContainerExists() bool {
	result, _ := w.CmdHelper.CliRunner("docker", nil, []string{"ps", "-a", "--format", `table {{.Names}}`}...)
	return strings.Contains(result, w.Container)
}

// src: /etc/os-release	 dst: ./
func (w *Wrapper) ContainerCP(src, dst string) (string, error) {
	result, err := w.CmdHelper.CliRunner("docker", nil, []string{"cp", w.Container + ":" + src, dst}...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (w *Wrapper) DeleteContainer() (string, error) {
	result, err := w.CmdHelper.CliRunner("docker", nil, []string{"rm", "-f", w.Container}...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (w *Wrapper) Validate(validators []string) (bool, error) {
	history, err := w.CmdHelper.CliRunner("docker", nil, []string{"history", w.Image, "--no-trunc"}...)
	if err != nil {
		return false, err
	}
	for _, item := range validators {
		if strings.Contains(history, item) {
			return true, nil
		}
	}
	return false, nil
}

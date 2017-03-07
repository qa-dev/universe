package container

import (
	"io/ioutil"
	"plugin"
)

type ContainerService struct{}

func (s *ContainerService) LoadPlugins() error {
	files, err := ioutil.ReadDir("../modules")
	if err != nil {
		return err
	}
	for _, file := range files {
		plug, err := plugin.Open(file.Name())
		if err != nil {
			return err
		}
		_, err = plug.Lookup("LoadPlugin")
		if err != nil {
			return err
		}
	}
	return nil
}

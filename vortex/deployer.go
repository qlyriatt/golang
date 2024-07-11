package main

import "errors"

type Deployer interface {
	CreatePod(name string) error
	DeletePod(name string) error
	GetPodList() ([]string, error)
}

// пример сервиса, подходящего под интерфейс Deployer
type SimpleDeployer struct {
	pods map[string]struct{}
}

func NewDeployer() Deployer {
	sd := SimpleDeployer{}
	sd.pods = make(map[string]struct{})
	return &sd
}

func (sd *SimpleDeployer) CreatePod(name string) error {
	if _, ok := sd.pods[name]; ok {
		return errors.New("Pod" + name + "already exists")
	}

	sd.pods[name] = struct{}{}
	return nil
}

func (sd *SimpleDeployer) DeletePod(name string) error {
	if _, ok := sd.pods[name]; !ok {
		return errors.New("Pod" + name + "does not exist")
	}

	delete(sd.pods, name)
	return nil
}

func (sd *SimpleDeployer) GetPodList() ([]string, error) {
	var names []string
	for name := range sd.pods {
		names = append(names, name)
	}

	return names, nil
}

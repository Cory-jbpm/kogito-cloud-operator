// Copyright 2020 Red Hat, Inc. and/or its affiliates
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

package steps

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/kiegroup/kogito-cloud-operator/test/framework"
)

const sourceLocation = "src/main/resources"

func registerKogitoDeployFilesSteps(s *godog.Suite, data *Data) {
	// Deploy steps
	s.Step(`^Deploy file "([^"]*)" from example service "([^"]*)"$`, data.deployFileFromExampleService)
	s.Step(`^Deploy folder from example service "([^"]*)"$`, data.deployFolderFromExampleService)
}

// Deploy steps

func (data *Data) deployFileFromExampleService(file, serviceName string) error {
	sourceFilePath := fmt.Sprintf(`%s/%s/%s/%s`, data.KogitoExamplesLocation, serviceName, sourceLocation, file)
	return deploySourceFilesFromPath(data.Namespace, serviceName, sourceFilePath)
}

func (data *Data) deployFolderFromExampleService(serviceName string) error {
	sourceFolderPath := fmt.Sprintf(`%s/%s/%s`, data.KogitoExamplesLocation, serviceName, sourceLocation)
	return deploySourceFilesFromPath(data.Namespace, serviceName, sourceFolderPath)
}

func deploySourceFilesFromPath(namespace, serviceName, path string) error {
	framework.GetLogger(namespace).Infof("Deploy example %s with source files in path %s", serviceName, path)

	kogitoAppHolder, err := getKogitoAppHolder(namespace, "quarkus", serviceName, &messages.PickleStepArgument_PickleTable{})
	if err != nil {
		return err
	}
	kogitoAppHolder.Spec.Build.GitSource.URI = path

	return framework.DeployService(namespace, framework.CLIInstallerType, kogitoAppHolder.KogitoApp)
}

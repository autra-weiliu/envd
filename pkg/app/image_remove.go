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

package app

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/driver/docker"
)

var CommandRemoveImage = &cli.Command{
	Name:    "remove",
	Aliases: []string{"r", "rm"},
	Usage:   "Remove an envd image",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "image",
			Usage:   "Specify the image name to be removed",
			Aliases: []string{"i"},
		},
		&cli.StringFlag{
			Name:        "tag",
			Usage:       "Remove the image with a specific tag",
			Aliases:     []string{"t"},
			DefaultText: "dev",
		},
	},
	Action: removeImage,
}

func removeImage(clicontext *cli.Context) error {
	imageName := clicontext.String("image")
	if imageName == "" {
		return errors.New("image name is required, find images by `envd images list`")
	}
	tag := clicontext.String("tag")
	if tag == "" {
		logrus.Debug("tag not specified, using default tag: `dev`")
		tag = "dev"
	}
	imageNameWithTag := fmt.Sprintf("%s:%s", imageName, tag)

	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return err
	}
	if err := dockerClient.RemoveImage(clicontext.Context, imageNameWithTag); err != nil {
		return errors.Errorf("remove image %s failed: %w", imageNameWithTag, err)
	}
	logrus.Infof("image(%s) has removed", imageNameWithTag)
	return nil
}

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
	"time"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandPrune = &cli.Command{
	Name:     "prune",
	Category: CategorySettings,
	Usage:    "Clean up the build cache",
	Flags: []cli.Flag{
		&cli.DurationFlag{
			Name:  "keep-duration",
			Usage: "Keep data newer than this limit",
		},
		&cli.Float64Flag{
			Name:  "keep-storage",
			Usage: "Keep data below this limit (in MB)",
		},
		&cli.StringSliceFlag{
			Name:   "filter, f",
			Usage:  "Filter records",
			Hidden: true,
		},
		&cli.BoolFlag{
			Name:  "all",
			Usage: "Include internal caches (oh-my-zsh, vscode extensions and other envd caches)",
		},
		&cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Verbose output",
		},
	},

	Action: prune,
}

func prune(clicontext *cli.Context) error {
	cleanAll := clicontext.Bool("all")
	if cleanAll {
		if err := home.GetManager().CleanCache(); err != nil {
			return errors.Wrap(err, "failed to clean internal cache")
		}
	}
	defer func(start time.Time) {
		telemetry.GetReporter().Telemetry(
			"prune", telemetry.AddField("duration", time.Since(start).Seconds()))
	}(time.Now())

	keepDuration := clicontext.Duration("keep-duration")
	maxUsed := clicontext.Float64("max-used-mb")
	minFree := clicontext.Float64("min-free-mb")
	reserved := clicontext.Float64("reserved-mb")
	filter := clicontext.StringSlice("filter")
	verbose := clicontext.Bool("verbose")

	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	var bkClient buildkitd.Client
	if c.Builder == types.BuilderTypeMoby {
		bkClient, err = buildkitd.NewMobyClient(clicontext.Context,
			c.Builder, c.BuilderAddress, nil)
		if err != nil {
			return errors.Wrap(err, "failed to create moby buildkit client")
		}
	} else {
		bkClient, err = buildkitd.NewClient(clicontext.Context,
			c.Builder, c.BuilderAddress, nil)
		if err != nil {
			return errors.Wrap(err, "failed to create buildkit client")
		}
	}
	if err := bkClient.Prune(clicontext.Context,
		keepDuration, maxUsed, minFree, reserved, filter, verbose, cleanAll); err != nil {
		return errors.Wrap(err, "failed to prune buildkit cache")
	}
	return nil
}

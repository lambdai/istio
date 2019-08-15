// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mesh

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"istio.io/operator/pkg/helm"
	"istio.io/operator/pkg/util"
	"istio.io/pkg/log"
)

func profileDiffCmd(rootArgs *rootArgs) *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Diffs two Istio configuration profiles.",
		Long:  "The diff subcommand is used to display the difference between two Istio configuration profiles.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profileDiff(rootArgs, args)
		}}

}

// profileDiff compare two profile files.
func profileDiff(rootArgs *rootArgs, args []string) {
	checkLogsOrExit(rootArgs)

	a, err := helm.ReadValuesYAML(args[0])
	if err != nil {
		log.Errorf("could not read the profile values from %s: %s", args[0], err)
		os.Exit(1)
	}

	b, err := helm.ReadValuesYAML(args[1])
	if err != nil {
		log.Errorf("could not read the profile values from %s: %s", args[1], err)
		os.Exit(1)
	}

	diff := util.YAMLDiff(a, b)
	if diff == "" {
		fmt.Println("Profiles are identical")
	} else {
		fmt.Printf("Difference of profiles are:\n%s", diff)
		os.Exit(1)
	}
}

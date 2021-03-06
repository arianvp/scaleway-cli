// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// RmiArgs are flags for the `RunRmi` function
type RmiArgs struct {
	Images []string
}

// RunRmi is the handler for 'scw rmi'
func RunRmi(ctx CommandContext, args RmiArgs) error {
	hasError := false
	for _, needle := range args.Images {
		image := ctx.API.GetImageID(needle, true)
		err := ctx.API.DeleteImage(image)
		if err != nil {
			logrus.Errorf("failed to delete image %s: %s", image, err)
			hasError = true
		} else {
			fmt.Fprintln(ctx.Stdout, needle)
		}
	}
	if hasError {
		return fmt.Errorf("at least 1 image failed to be removed")
	}
	return nil
}

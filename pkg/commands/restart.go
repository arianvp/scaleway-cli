// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

// RestartArgs are flags for the `RunRestart` function
type RestartArgs struct {
	Wait    bool
	Timeout float64
	Servers []string
}

// restartIdentifiers resolves server IDs, restarts, and waits for them to be ready (-w)
func restartIdentifiers(ctx CommandContext, wait bool, servers []string, cr chan string) {
	var wg sync.WaitGroup
	for _, needle := range servers {
		wg.Add(1)
		go func(needle string) {
			defer wg.Done()
			server := ctx.API.GetServerID(needle)
			res := server
			err := ctx.API.PostServerAction(server, "reboot")
			if err != nil {
				if err.Error() != "server is being stopped or rebooted" {
					logrus.Errorf("failed to restart server %s: %s", server, err)
				}
				res = ""
			} else {
				if wait {
					// FIXME: handle gateway
					api.WaitForServerReady(ctx.API, server, "")
				}
			}
			cr <- res
		}(needle)
	}
	wg.Wait()
	close(cr)
}

// RunRestart is the handler for 'scw restart'
func RunRestart(ctx CommandContext, args RestartArgs) error {
	if args.Wait && args.Timeout > 0 {
		go func() {
			time.Sleep(time.Duration(args.Timeout*1000) * time.Millisecond)
			// FIXME: avoid use of fatalf
			logrus.Fatalf("Operation timed out")
		}()
	}

	cr := make(chan string)
	go restartIdentifiers(ctx, args.Wait, args.Servers, cr)
	done := false
	hasError := false

	for !done {
		select {
		case uuid, more := <-cr:
			if !more {
				done = true
				break
			}
			if len(uuid) > 0 {
				fmt.Fprintln(ctx.Stdout, uuid)
			} else {
				hasError = true
			}
		}
	}

	if hasError {
		return fmt.Errorf("at least 1 server failed to restart")
	}
	return nil
}

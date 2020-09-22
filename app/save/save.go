package save

import (
	"encoding/json"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
	"github.com/stalker-loki/app/slog"
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/pod"
)

// Pod saves the configuration to the configured location
func Pod(c *pod.Config) (success bool) {
	slog.Trace("saving configuration to", *c.ConfigFile)
	var uac cli.StringSlice
	if len(*c.UserAgentComments) > 0 {
		uac = make(cli.StringSlice, len(*c.UserAgentComments))
		copy(uac, *c.UserAgentComments)
		*c.UserAgentComments = uac[1:]
	}
	// don't save pipe log setting as we want it to only be active from a
	// flag or environment variable
	pipeLogOn := *c.PipeLog
	*c.PipeLog = false
	if yp, e := json.MarshalIndent(c, "", "  "); e == nil {
		apputil.EnsureDir(*c.ConfigFile)
		if e := ioutil.WriteFile(*c.ConfigFile, yp, 0600); e != nil {
			slog.Error(e)
			success = false
		}
		success = true
	}
	*c.UserAgentComments = uac
	*c.PipeLog = pipeLogOn
	return
}

// Filters saves the logger per-package logging configuration
func Filters(dataDir string) func(pkgs Pk.Package) (success bool) {
	return func(pkgs Pk.Package) (success bool) {
		if filterJSON, e := json.MarshalIndent(pkgs, "", "  "); e == nil {
			slog.Trace("Saving log filter:\n```", string(filterJSON), "\n```")
			apputil.EnsureDir(dataDir)
			if e := ioutil.WriteFile(filepath.Join(dataDir, "log-filter.json"), filterJSON,
				0600); slog.Check(e) {
				success = false
			}
			success = true
		}
		return
	}
}

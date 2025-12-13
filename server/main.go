package main

import (
	"embed"

	"github.com/WuKongIM/WuKongIM/cmd"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"github.com/WuKongIM/WuKongIM/version"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

//go:embed web/dist
var webFS embed.FS

//go:embed kefu
var kefuFS embed.FS

// go ldflags
var Version string    // version
var Commit string     // git commit id
var CommitDate string // git commit date
var TreeState string  // git tree state

func main() {

	version.Version = Version
	version.Commit = Commit
	version.CommitDate = CommitDate
	version.TreeState = TreeState
	version.WebFs = webFS
	version.KefuFS = kefuFS

	undo, err := maxprocs.Set()
	defer undo()
	if err != nil {
		wklog.Warn("maxprocs set error", zap.Error(err))
	}

	cmd.Execute()

}

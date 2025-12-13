package main

import (
	"embed"

	"github.com/WuKongIM/WuKongIM/cmd"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"github.com/WuKongIM/WuKongIM/version"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

//go:embed admin/web/im/dist
var imWebFS embed.FS

//go:embed admin/web/kefu/dist
var kefuWebFS embed.FS

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
	version.ImWebFs = imWebFS
	version.KefuWebFS = kefuWebFS

	undo, err := maxprocs.Set()
	defer undo()
	if err != nil {
		wklog.Warn("maxprocs set error", zap.Error(err))
	}

	cmd.Execute()

}

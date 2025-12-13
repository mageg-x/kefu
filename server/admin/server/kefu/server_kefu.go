package kefu

import (
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/WuKongIM/WuKongIM/internal/options"
	"github.com/WuKongIM/WuKongIM/internal/server"
	"github.com/WuKongIM/WuKongIM/pkg/wkhttp"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"github.com/WuKongIM/WuKongIM/version"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type KefuServer struct {
	r    *wkhttp.WKHttp
	addr string
	s    *server.Server
	wklog.Log
}

// NewKefuServer new一个kefu server
func NewKefuServer(s *server.Server) *KefuServer {
	// r := wkhttp.New()
	log := wklog.NewWKLog("KefuServer")
	r := wkhttp.NewWithLogger(wkhttp.LoggerWithWklog(log))
	r.Use(wkhttp.CORSMiddleware())

	ds := &KefuServer{
		r:    r,
		addr: options.G.Kefu.Addr,
		s:    s,
		Log:  log,
	}
	return ds
}

// Start 开始
func (s *KefuServer) Start() {
	s.r.GetGinRoute().Use(gzip.Gzip(gzip.DefaultCompression))

	st, _ := fs.Sub(version.KefuWebFS, "admin/web/kefu/dist")
	s.r.GetGinRoute().NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/admin/kefu") {
			c.FileFromFS("./", http.FS(st))
			return
		}
	})

	s.r.GetGinRoute().StaticFS("/admin/kefu", http.FS(st))

	s.setRoutes()
	go func() {
		err := s.r.Run(s.addr) // listen and serve
		if err != nil {
			panic(err)
		}
	}()
	s.Info("Kefu server started", zap.String("addr", s.addr))

	_, port := parseAddr(s.addr)
	s.Info(fmt.Sprintf("Kefu web address： http://localhost:%d/admin/kefu", port))
}

// Stop 停止服务
func (s *KefuServer) Stop() {
}

func (s *KefuServer) setRoutes() {

}

func parseAddr(addr string) (string, int64) {
	addrPairs := strings.Split(addr, ":")
	if len(addrPairs) < 2 {
		return "", 0
	}
	portInt64, _ := strconv.ParseInt(addrPairs[len(addrPairs)-1], 10, 64)
	return addrPairs[0], portInt64
}

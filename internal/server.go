package internal

import (
	"log/slog"

	"lkgiovani_go_rinha_2025/internal/http"
	"lkgiovani_go_rinha_2025/internal/router"

	"github.com/panjf2000/gnet/v2"
)

type HttpServer struct {
	gnet.BuiltinEventEngine
	Addr      string
	Multicore bool
	Eng       gnet.Engine
}

func (hs *HttpServer) OnBoot(eng gnet.Engine) gnet.Action {
	hs.Eng = eng
	slog.Info("the server has started", "multi-core", hs.Multicore, "addr", hs.Addr)
	return gnet.None
}

func (hs *HttpServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(http.NewHttpCodec())
	return nil, gnet.None
}

func (hs *HttpServer) OnTraffic(c gnet.Conn) gnet.Action {
	hc := c.Context().(*http.HttpCodec)
	buf, _ := c.Next(-1)

pipeline:
	nextOffset, err := hc.Parse(buf, c)
	if err != nil {
		goto response
	}

	hc.ResetParser()
	router.HandleRequest(hc)
	buf = buf[nextOffset:]
	if len(buf) > 0 {
		goto pipeline
	}

response:
	c.Write(hc.GetBuffer())
	hc.Reset()
	return gnet.None
}

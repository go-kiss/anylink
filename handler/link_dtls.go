package handler

import (
	"net"
	"time"

	"github.com/go-kiss/anylink/base"
	"github.com/go-kiss/anylink/sessdata"
)

func LinkDtls(conn net.Conn, cSess *sessdata.ConnSession) {
	dSess := cSess.NewDtlsConn()
	if dSess == nil {
		// 创建失败，直接关闭链接
		_ = conn.Close()
		return
	}

	defer func() {
		base.Debug("LinkDtls return", cSess.IpAddr)
		_ = conn.Close()
		dSess.Close()
	}()

	var (
		dead = time.Duration(cSess.CstpDpd+5) * time.Second
	)

	go dtlsWrite(conn, dSess, cSess)

	now := time.Now()

	for {

		if time.Now().Sub(now) > time.Second*30 {
			// return
		}

		err := conn.SetReadDeadline(time.Now().Add(dead))
		if err != nil {
			base.Error("SetDeadline: ", err)
			return
		}
		hdata := make([]byte, BufferSize)
		n, err := conn.Read(hdata)
		if err != nil {
			base.Error("read hdata: ", err)
			return
		}

		// 限流设置
		err = cSess.RateLimit(n, true)
		if err != nil {
			base.Error(err)
		}

		switch hdata[0] {
		case sessdata.TypeKeepAlive:
			// do nothing
			base.Debug("recv keepalive", cSess.IpAddr)
		case sessdata.TypeDsiconnect:
			base.Debug("DISCONNECT", cSess.IpAddr)
			return
		case sessdata.TypeDpdReq: // DPD-REQ
			// base.Debug("recv DPD-REQ", cSess.IpAddr)
			payload := &sessdata.Payload{
				Type: sessdata.TypeDpdResp,
				Data: nil,
			}

			select {
			case cSess.PayloadOutDtls <- payload:
			case <-dSess.CloseChan:
				return
			}
		case sessdata.TypeDpdResp:
			// base.Debug("recv DPD-RESP", cSess.IpAddr)
		case sessdata.TypeData:
			if payloadIn(cSess, sessdata.TypeData, hdata[1:n]) {
				return
			}
		}
	}
}

func dtlsWrite(conn net.Conn, dSess *sessdata.DtlsSession, cSess *sessdata.ConnSession) {
	defer func() {
		base.Debug("dtlsWrite return", cSess.IpAddr)
		_ = conn.Close()
		dSess.Close()
	}()

	var (
		header  []byte
		payload *sessdata.Payload
	)

	for {
		// dtls优先推送数据
		select {
		case payload = <-cSess.PayloadOutDtls:
		case <-dSess.CloseChan:
			return
		}

		header = []byte{payload.Type}
		header = append(header, payload.Data...)
		n, err := conn.Write(header)
		if err != nil {
			base.Error("write err", err)
			return
		}

		// 限流设置
		err = cSess.RateLimit(n, false)
		if err != nil {
			base.Error(err)
		}
	}
}

package handler

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/go-kiss/anylink/base"
	"github.com/go-kiss/anylink/sessdata"
)

func LinkCstp(conn net.Conn, cSess *sessdata.ConnSession) {
	defer func() {
		base.Debug("LinkCstp return", cSess.IpAddr)
		_ = conn.Close()
		cSess.Close()
	}()

	var (
		err     error
		n       int
		dataLen uint16
		dead    = time.Duration(cSess.CstpDpd+5) * time.Second
	)

	go cstpWrite(conn, cSess)

	for {

		// 设置超时限制
		err = conn.SetReadDeadline(time.Now().Add(dead))
		if err != nil {
			base.Error("SetDeadline: ", err)
			return
		}
		hdata := make([]byte, BufferSize)
		n, err = conn.Read(hdata)
		if err != nil {
			base.Error("read hdata: ", err)
			return
		}

		// 限流设置
		err = cSess.RateLimit(n, true)
		if err != nil {
			base.Error(err)
		}

		switch hdata[6] {
		case sessdata.TypeKeepAlive:
			// do nothing
			// base.Debug("recv keepalive", cSess.IpAddr)
		case sessdata.TypeDsiconnect:
			base.Debug("DISCONNECT", cSess.IpAddr)
			return
		case sessdata.TypeDpdReq:
			// base.Debug("recv DPD-REQ", cSess.IpAddr)
			if payloadOutCstp(cSess, 0x04, nil) {
				return
			}
		case sessdata.TypeDpdResp:
			// base.Debug("recv DPD-RESP")
		case sessdata.TypeData:
			dataLen = binary.BigEndian.Uint16(hdata[4:6]) // 4,5
			if payloadIn(cSess, 0x00, hdata[8:8+dataLen]) {
				return
			}

		}
	}
}

func cstpWrite(conn net.Conn, cSess *sessdata.ConnSession) {
	defer func() {
		base.Debug("cstpWrite return", cSess.IpAddr)
		_ = conn.Close()
		cSess.Close()
	}()

	var (
		err     error
		n       int
		header  []byte
		payload *sessdata.Payload
	)

	for {
		select {
		case payload = <-cSess.PayloadOutCstp:
		case <-cSess.CloseChan:
			return
		}

		header = []byte{'S', 'T', 'F', 0x01, 0x00, 0x00, payload.Type, 0x00}
		if payload.Type == sessdata.TypeData {
			binary.BigEndian.PutUint16(header[4:6], uint16(len(payload.Data)))
			header = append(header, payload.Data...)
		}
		n, err = conn.Write(header)
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

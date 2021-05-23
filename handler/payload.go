package handler

import (
	"github.com/go-kiss/anylink/dbdata"
	"github.com/go-kiss/anylink/sessdata"
	"github.com/songgao/water/waterutil"
)

func payloadIn(cSess *sessdata.ConnSession, pType byte, data []byte) bool {
	payload := &sessdata.Payload{
		PType: pType,
		Data:  data,
	}

	return payloadInData(cSess, payload)
}

func payloadInData(cSess *sessdata.ConnSession, payload *sessdata.Payload) bool {
	// 进行Acl规则判断
	check := checkLinkAcl(cSess.Group, payload)
	if !check {
		// 校验不通过直接丢弃
		return false
	}

	closed := false
	select {
	case cSess.PayloadIn <- payload:
	case <-cSess.CloseChan:
		closed = true
	}

	return closed
}

func payloadOut(cSess *sessdata.ConnSession, pType byte, data []byte) bool {
	dSess := cSess.GetDtlsSession()
	if dSess == nil {
		return payloadOutCstp(cSess, pType, data)
	} else {
		return payloadOutDtls(dSess, pType, data)
	}
}

func payloadOutCstp(cSess *sessdata.ConnSession, pType byte, data []byte) bool {
	payload := &sessdata.Payload{
		PType: pType,
		Data:  data,
	}

	closed := false

	select {
	case cSess.PayloadOutCstp <- payload:
	case <-cSess.CloseChan:
		closed = true
	}

	return closed
}

func payloadOutDtls(dSess *sessdata.DtlsSession, pType byte, data []byte) bool {
	payload := &sessdata.Payload{
		PType: pType,
		Data:  data,
	}

	select {
	case dSess.CSess.PayloadOutDtls <- payload:
	case <-dSess.CloseChan:
	}

	return false
}

// Acl规则校验
func checkLinkAcl(group *dbdata.Group, payload *sessdata.Payload) bool {
	if payload.PType != 0x00 || len(group.LinkAcl) == 0 {
		return true
	}

	ip_dst := waterutil.IPv4Destination(payload.Data)
	ip_port := waterutil.IPv4DestinationPort(payload.Data)

	// 优先放行dns端口
	for _, v := range group.ClientDns {
		if v.Val == ip_dst.String() && ip_port == 53 {
			return true
		}
	}

	for _, v := range group.LinkAcl {
		// 循环判断ip和端口
		if v.IpNet.Contains(ip_dst) {
			if v.Port == ip_port || v.Port == 0 {
				if v.Action == dbdata.Allow {
					return true
				} else {
					return false
				}
			}
		}
	}

	return false
}

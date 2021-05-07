package nebula

import "github.com/slackhq/nebula/header"

func HandleIncomingHandshake(f *Interface, addr *udpAddr, packet []byte, h *header.H, hostinfo *HostInfo) {
	if !f.lightHouse.remoteAllowList.Allow(addr.IP) {
		f.l.WithField("udpAddr", addr).Debug("lighthouse.remote_allow_list denied incoming handshake")
		return
	}

	switch h.Subtype {
	case header.HandshakeIXPSK0:
		switch h.MessageCounter {
		case 1:
			ixHandshakeStage1(f, addr, packet, h)
		case 2:
			newHostinfo, _ := f.handshakeManager.QueryIndex(h.RemoteIndex)
			tearDown := ixHandshakeStage2(f, addr, newHostinfo, packet, h)
			if tearDown && newHostinfo != nil {
				f.handshakeManager.DeleteHostInfo(newHostinfo)
			}
		}
	}

}

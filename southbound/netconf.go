package southbound

import (
	"encoding/xml"

	"github.com/Juniper/go-netconf/netconf"
)

func execRpc[T any](rpcMethod string, host string, s *SouthboundImpl) (T, error) {
	var t T

	sshConfig := netconf.SSHConfigPassword(s.LoginUser, s.LoginPassword)
	session, err := netconf.DialSSH(host, sshConfig)

	if err != nil {
		return t, err
	}

	defer session.Close()

	reply, err := session.Exec(netconf.RawMethod(rpcMethod))
	if err != nil {
		return t, err
	}

	err = parseXmlResponse(reply, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

func parseXmlResponse[T any](reply *netconf.RPCReply, t *T) error {
	err := xml.Unmarshal([]byte(reply.RawReply), t)
	if err != nil {
		return err
	}

	return nil
}

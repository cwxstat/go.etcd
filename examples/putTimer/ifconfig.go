package main

import (
	"fmt"
	"net"
	"strings"
)

func Ifconfig() string {
	s := fmt.Sprintln("=== interfaces ===")

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		s += fmt.Sprintln("net.Interface:", iface)

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			addrStr := addr.String()
			s += fmt.Sprintln("    net.Addr: ", addr.Network(), addrStr)

			// Must drop the stuff after the slash in order to convert it to an IP instance
			split := strings.Split(addrStr, "/")
			addrStr0 := split[0]

			// Parse the string to an IP instance
			ip := net.ParseIP(addrStr0)
			if ip.To4() != nil {
				s += fmt.Sprintln("       ", addrStr0, "is ipv4")
			} else {
				s += fmt.Sprintln("       ", addrStr0, "is ipv6")
			}
			s += fmt.Sprintln("       ", addrStr0, "is interface-local multicast :", ip.IsInterfaceLocalMulticast())
			s += fmt.Sprintln("       ", addrStr0, "is link-local multicast      :", ip.IsLinkLocalMulticast())
			s += fmt.Sprintln("       ", addrStr0, "is link-local unicast        :", ip.IsLinkLocalUnicast())
			s += fmt.Sprintln("       ", addrStr0, "is global unicast            :", ip.IsGlobalUnicast())
			s += fmt.Sprintln("       ", addrStr0, "is multicast                 :", ip.IsMulticast())
			s += fmt.Sprintln("       ", addrStr0, "is loopback                  :", ip.IsLoopback())

		}
	}
	return s
}

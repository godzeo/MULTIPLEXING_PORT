package protocols

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// RunServer runs protoplex
func RunServer(bind string, p []*Protocol, logger zerolog.Logger) {
	//logger = logger.With().Str("module", "listener").Logger()

	if len(p) == 0 {
		fmt.Println("[Err] Shunt function startup failed！")
	} else {
		fmt.Println("[Info] Shunt function startup success！")
		for _, proto := range p {
			if proto.Name == "SOCKS5" {
				fmt.Println("[Info] SOCKS5_forward_addr:", proto.Target)
			} else if proto.Name == "HTTP" {
				fmt.Println("[Info] HTTP_forward_addr:", proto.Target)
			}
			//logger.Info().Str("protocol", proto.Name).Str("target", proto.Target).Msgf("- %s @ %s", proto.Name, proto.Target)
		}
	}

	listener, err := net.Listen("tcp", bind)
	if err != nil {
		fmt.Println("[Err] Shunt port listening failed！")
		//logger.Fatal().Str("bind", bind).Err(err).Msg("Unable to create listener.")
		os.Exit(1)
	}
	defer listener.Close()
	//logger.Info().Str("bind", listener.Addr().String()).Msg("Listening...")
	fmt.Println("[Info] Shunt port listening succeeded, address:", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			//logger.Debug().Err(err).Msg("Error while accepting connection.")
			fmt.Println("[Err] Send error, error message:", err)
		}
		go ConnectionHandler(conn, p, logger.With().Str("module", "handler").Str("ip", conn.RemoteAddr().String()).Logger())
	}
}

// ConnectionHandler connects a net.Conn with a proxy target given a list of protocols
func ConnectionHandler(conn net.Conn, p []*Protocol, logger zerolog.Logger) {
	defer conn.Close() // the connection must close after this goroutine exits

	identifyBuffer := make([]byte, 1024) // at max 1KB buffer to identify payload

	// read the handshake into our buffer
	_ = conn.SetReadDeadline(time.Now().Add(15 * time.Second)) // 15-second timeout to identify
	n, err := conn.Read(identifyBuffer)
	if err != nil {
		logger.Debug().Err(err).Msg("Identify read error. Connection closed.")
		return
	}
	_ = conn.SetReadDeadline(time.Time{}) // reset our timeout

	// determine the protocol
	protocol := DetermineProtocol(identifyBuffer[:n], p)
	if protocol == nil { // unsuccessful protocol identify, close and forget
		logger.Debug().Msg("Protocol unrecognized. Connection closed.")
		return
	}
	//logger = logger.With().Str("protocol", protocol.Name).Str("target", protocol.Target).Logger()
	//logger.Debug().Msg("Protocol recognized.")
	if protocol.Name == "HTTP" {
		fmt.Println("[Info] The HTTP request is received. Procedure Address requested:", conn.RemoteAddr().String())
	} else if protocol.Name == "SOCKS5" {
		fmt.Println("[Info] The SOCKS5 request is received. Procedure Address requested:", conn.RemoteAddr().String())
	} else if protocol.Name == "SSH" {
		fmt.Println("[Info] The SSH request is received. Procedure Address requested:", conn.RemoteAddr().String())
	} else if protocol.Name == "RDP" {
		fmt.Println("[Info] The RDP request is received. Procedure Address requested:", conn.RemoteAddr().String())
	}

	// establish our connection with the target
	targetConn, err := net.Dial("tcp", protocol.Target)
	if err != nil {
		logger.Debug().Err(err).Msg("Remote connection unsuccessful.")
		return // we were unable to establish the connection with the proxy target
	}
	defer targetConn.Close()
	_, err = targetConn.Write(identifyBuffer[:n]) // tell them everything they just told us
	if err != nil {
		logger.Debug().Err(err).Msg("Remote disconnected us during identify.")
		return // remote rejected us?? okay.
	}

	closed := make(chan bool, 2)
	go proxy(conn, targetConn, closed)
	go proxy(targetConn, conn, closed)

	// wait for any connection to close
	<-closed
	//logger.Debug().Msg("Connection closed.")
	if protocol.Name == "HTTP" {
		fmt.Println("[Info] The HTTP request was closed. Address requested:", conn.RemoteAddr().String())
	} else if protocol.Name == "SOCKS5" {
		fmt.Println("[Info] The SOCKS5 request was closed. Address requested:", conn.RemoteAddr().String())
	} else if protocol.Name == "SSH" {
		fmt.Println("[Info] The SSH request was closed. Address requested:", conn.RemoteAddr().String())
	} else if protocol.Name == "RDP" {
		fmt.Println("[Info] The RDP request was closed. Address requested:", conn.RemoteAddr().String())
	}
}

// DetermineProtocol determines a Protocol based on a given handshake
func DetermineProtocol(data []byte, p []*Protocol) *Protocol {
	dataLength := len(data)
	for _, protocol := range p {
		// since every protocol is different, let's limit the way we match things
		if (protocol.NoComparisonBeforeBytes != 0 && dataLength < protocol.NoComparisonBeforeBytes) ||
			(protocol.NoComparisonAfterBytes != 0 && dataLength > protocol.NoComparisonAfterBytes) {
			continue // avoids unnecessary comparisons
		}

		// compare against bytestrings first for efficiency
		// first "contains" (due to ALPNs we can't match against TLS start bytes first)
		for _, byteSlice := range protocol.MatchBytes {
			byteSliceLength := len(byteSlice)
			if dataLength < byteSliceLength {
				continue
			}
			if bytes.Contains(data, byteSlice) {
				return protocol
			}
		}
		// then against prefixes
		for _, byteSlice := range protocol.MatchStartBytes {
			byteSliceLength := len(byteSlice)
			if dataLength < byteSliceLength {
				continue
			}
			if bytes.Equal(byteSlice, data[:byteSliceLength]) {
				return protocol
			}
		}

		// let's use regex matching as a last resort
		for _, regex := range protocol.MatchRegexes {
			if regex.Match(data) {
				return protocol
			}
		}
	}
	return nil
}

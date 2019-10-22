package inpipe

import (
	"io"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/covrom/hls-streamer/manifestgenerator"
)

func InTCP(serveTCP string, readBufferSize int, mg *manifestgenerator.ManifestGenerator, log *logrus.Logger) {
	buf := make([]byte, 0, readBufferSize)
	listener, err := net.Listen("tcp", serveTCP)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("TCP server listening %s", serveTCP)

	chconn := make(chan net.Conn)

	go func() {
		for conn := range chconn {
			for {
				n, err := conn.Read(buf[:cap(buf)])

				if n == 0 && err == io.EOF {
					log.Info("Closing process detected EOF")
					break
				}

				if err != nil && err != io.EOF {
					// Error reading pipe
					log.Fatal(err)
					os.Exit(1)
				}

				// process buf
				log.Debug("Sent to process: ", n, " bytes")
				mg.AddData(buf[:n])
			}
			log.Debug("EOF source from: ", conn.RemoteAddr())
			conn.Close()
		}
	}()

	for {
		log.Println("Wait for connecting source...")
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		log.Debug("Connected source from: ", conn.RemoteAddr())
		select {
		case chconn <- conn:
		default:
			log.Debug("Decline source from: ", conn.RemoteAddr())
			conn.Close()
		}
	}

	// listener.Close()
}

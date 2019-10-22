package inpipe

import (
	"bufio"
	"io"
	"os"

	"github.com/covrom/hls-streamer/manifestgenerator"

	"github.com/sirupsen/logrus"
)

func InPipe(readBufferSize int, mg *manifestgenerator.ManifestGenerator, log *logrus.Logger) {
	buf := make([]byte, 0, readBufferSize)
	r := bufio.NewReader(os.Stdin)
	// Buffer
	for {
		n, err := r.Read(buf[:cap(buf)])
		if n == 0 && err == io.EOF {
			// Detected EOF
			// Closing
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
}

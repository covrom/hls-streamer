package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/covrom/hls-streamer/hls"
	"github.com/covrom/hls-streamer/httpserver"
	"github.com/covrom/hls-streamer/inpipe"
	"github.com/covrom/hls-streamer/logger"
	"github.com/covrom/hls-streamer/manifestgenerator"
	"github.com/covrom/hls-streamer/mediachunk"
)

const (
	readBufferSize = 128
)

var (
	verbose            = flag.Bool("v", false, "enable to get verbose logging")
	baseOutPath        = flag.String("p", "./results", "Output path")
	chunkBaseFilename  = flag.String("f", "chunk_", "Chunks base filename")
	chunkListFilename  = flag.String("cf", "chunklist.m3u8", "Chunklist filename")
	targetSegmentDurS  = flag.Float64("t", 4.0, "Target chunk duration in seconds")
	liveWindowSize     = flag.Int("w", 3, "Live window size in chunks")
	lhlsAdvancedChunks = flag.Int("l", 0, "If > 0 activates LHLS, and it indicates the number of advanced chunks to create")
	manifestTypeInt    = flag.Int("m", int(hls.LiveWindow), "Manifest to generate (0- Vod, 1- Live event, 2- Live sliding window")
	autoPID            = flag.Bool("apids", true, "Enable auto PID detection, if true no need to pass vpid and apid")
	videoPID           = flag.Int("vpid", -1, "Video PID to parse")
	audioPID           = flag.Int("apid", -1, "Audio PID to parse")
	chunkInitType      = flag.Int("i", int(manifestgenerator.ChunkInitStart), "Indicates where to put the init data PAT and PMT packets (0- No ini data, 1- Init segment, 2- At the begining of each chunk")
	destinationType    = flag.Int("d", 1, "Indicates where the destination (0- No output, 1- File + flag indicator, 2- HTTP chunked transfer)")
	httpScheme         = flag.String("protocol", "http", "HTTP Scheme (http, https)")
	httpHost           = flag.String("host", "localhost:9094", "HTTP Host")
	serveTCP           = flag.String("tcp", "localhost:9555", "Enable TCP server at this port instead of stdin input stream")
	serveHttp          = flag.String("http", ":9099", "Enable http server at this port")
)

func main() {
	flag.Parse()

	var log = logger.ConfigureLogger(*verbose)

	log.Info(manifestgenerator.Version)
	log.Info("Started tssegmenter")

	if *autoPID == false && manifestgenerator.ChunkInitTypes(*chunkInitType) != manifestgenerator.ChunkNoIni {
		log.Error("Manual PID mode and Chunk No ini data are not compatible")
		os.Exit(1)
	}

	chunkOutputType := mediachunk.OutputTypes(*destinationType)
	hlsOutputType := hls.OutputTypes(*destinationType)

	// Creating output dir if does not exists
	if chunkOutputType == mediachunk.ChunkOutputModeFile || hlsOutputType == hls.HlsOutputModeFile {
		os.MkdirAll(*baseOutPath, 0744)
	}

	tr := http.DefaultTransport
	client := http.Client{
		Transport: tr,
		Timeout:   0,
	}

	mg := manifestgenerator.New(log,
		chunkOutputType,
		hlsOutputType,
		*baseOutPath,
		*chunkBaseFilename,
		*chunkListFilename,
		*targetSegmentDurS,
		manifestgenerator.ChunkInitTypes(*chunkInitType),
		*autoPID,
		-1,
		-1,
		hls.ManifestTypes(*manifestTypeInt),
		*liveWindowSize,
		*lhlsAdvancedChunks,
		&client,
		*httpScheme,
		*httpHost,
	)

	if *serveHttp != "" && *destinationType == 1 {
		httpserver.HTTPServer(*baseOutPath, *chunkListFilename, *serveHttp, log)
	}

	// Reader
	if *serveTCP == "" {
		inpipe.InPipe(readBufferSize, &mg, log)
	} else {
		inpipe.InTCP(*serveTCP, readBufferSize, &mg, log)
	}

	mg.Close()

	log.Info("Exit because detected EOF in the input pipe")

	os.Exit(0)
}

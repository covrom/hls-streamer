package httpserver

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func HTTPServer(baseOutPath, chunkListFilename, serveHttpAddr string, log *logrus.Logger) {
	fs := http.FileServer(http.Dir(baseOutPath))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html lang="en">
				<head>
					<meta charset=utf-8/>
					<link href="https://unpkg.com/video.js/dist/video-js.min.css" rel="stylesheet">
					<script src="https://unpkg.com/video.js/dist/video.min.js"></script>
				</head>
				<body>
				<video
					id="my-player"
					class="video-js"
					controls
					preload="auto"
					poster=""
					data-setup='{}'>
				<source  src="/video/` + chunkListFilename + `" type="application/x-mpegURL"></source>
				<p class="vjs-no-js">
					To view this video please enable JavaScript, and consider upgrading to a
					web browser that
					<a href="https://videojs.com/html5-video-support/" target="_blank">
					supports HTML5 video
					</a>
				</p>
				</video>
				</body>
			</html>`))
	})
	http.Handle("/video/", http.StripPrefix("/video/", fs))

	go http.ListenAndServe(serveHttpAddr, nil)

	log.Printf("HTTP server listening on %s", serveHttpAddr)
}

package actions

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/canhlinh/hlsdl"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func downloadCmdFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     "directly",
			Usage:    "represent the download of .m3u8 is directly from url or not",
			Required: true,
		},
		&cli.StringFlag{ // string
			Name: "type", // flag 名稱
			// Aliases:  []string{"d"}, // 別名
			Usage:    "specify the type of the target video (.m3u8 or .mp4)",
			Required: true, // flag true就是必填
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify the video name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "worker",
			Usage:    "specify the number of the worker for hlsdl (https://github.com/canhlinh/hlsdl)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "m3u8-path",
			Usage:    "specify the url of .m3u8",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "txt",
			Usage:    "specify the path of m3u8.txt, which contains the content of the target .m3u8",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "download-dir",
			Usage:    "specify the download directory. If download-dir is empty, the default download path will be the same as the executable file",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cookie",
			Usage:    "specify the customized \"cookie\" property of the http header",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "origin",
			Usage:    "specify the customized \"origin\" property of the http header",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "referer",
			Usage:    "specify the customized \"referer\" property of the http header",
			Required: false,
		},
	}
}

func DownloadVideo(c *cli.Context) (err error) {
	videoType := c.String("type")
	if videoType == "" {
		return fmt.Errorf("please specify the video type")
	}

	videoName := c.String("name")
	if videoName == "" {
		return fmt.Errorf("please specify the video name")
	}

	worker := c.Int("worker")
	if worker == 0 {
		worker = 1
	}

	switch strings.ToLower(videoType) {
	case "m3u8":
		return downloadFromM3U8(c, videoName, worker)
	default:
		return fmt.Errorf("unsupported video type: %s", videoType)
	}

}

func downloadFromM3U8(c *cli.Context, videoName string, worker int) (err error) {
	directly := c.Bool("directly")
	if !directly {
		return downloadFromLocalM3U8(c, videoName)
	}

	m3u8Url := c.String("m3u8-path")
	if m3u8Url == "" {
		return fmt.Errorf("please specify the url of .m3u8")
	}

	headers := make(map[string]string)
	cookie := c.String("cookie")
	if cookie != "" {
		headers["Cookie"] = cookie
	}

	origin := c.String("origin")
	if origin != "" {
		headers["Origin"] = origin
	}

	downloadDir := c.String("download-dir")
	outputDir := ""
	if downloadDir == "" {
		// default path is the same as the executable file
		exePath, _ := os.Executable()
		outputDir = filepath.Dir(exePath)
		downloadDir = filepath.Join(outputDir, "download") // filepath.Dirs exePath's directory

		// To support multi-process, the downloaded segments will be put into a subfolder folder under the download folder.
		// check download folder existence for exePath
		if !isExist(downloadDir) {
			os.MkdirAll(downloadDir, os.ModePerm)
		}
		// remove the download folder after this download is done
		defer os.RemoveAll(downloadDir)

		// put every segments under this subfolder
		uuid := uuid.New()
		downloadDir = filepath.Join(downloadDir, uuid.String())
	}

	hlsDL := hlsdl.New(m3u8Url, headers, downloadDir, worker, true, videoName)
	hlsPath, err := hlsDL.Download()
	if err != nil {
		return fmt.Errorf("failed to execute hls download, err=%s", err)
	}

	// check video name, must has suffix .ts
	if !strings.HasSuffix(videoName, ".ts") {
		videoName = fmt.Sprintf("%s.ts", videoName)
	}
	// copy the video (.ts) from hlsPath to outputDir
	videoFile := filepath.Join(outputDir, videoName)
	copy(hlsPath, videoFile)

	log.Printf("download complete, please refer to %s", outputDir)
	return
}

func downloadFromLocalM3U8(c *cli.Context, videoName string) (err error) {
	return
}

// isExist check file/folder was exists
func isExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func copy(src, dst string) (nBytes int64, err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	return io.Copy(destination, source)
}

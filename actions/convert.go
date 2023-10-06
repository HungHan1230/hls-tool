package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func convertCmdFlags() []cli.Flag {

	return []cli.Flag{
		&cli.StringFlag{
			Name:     "video-path",
			Aliases:  []string{"p"},
			Usage:    "specify the path of the video or the directory that includes videos.",
			Required: true,
		},
	}
}

func ConvertVideos(c *cli.Context) (err error) {
	videoPath := c.String("video-path")
	if videoPath == "" {
		return fmt.Errorf("please specify the path of the video or directory")
	}

	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return fmt.Errorf("failed to access videoPath, err=%s", err)
	}

	isDir := fileInfo.IsDir()

	successToDelete := false
	defer func() {
		if successToDelete {
			if isDir {
				files, _ := os.ReadDir(videoPath)
				for _, file := range files {
					if strings.HasSuffix(file.Name(), "ts") {
						os.Remove(filepath.Join(filepath.Dir(videoPath), file.Name()))
					}
				}

			} else {
				os.Remove(videoPath)
			}
		}
	}()

	if isDir {
		fmt.Println("Read files from directory:", videoPath)
		files, direrr := os.ReadDir(videoPath)
		if direrr != nil {
			return fmt.Errorf("failed to read directory, err=%s", direrr)
		}

		for _, file := range files {
			convert_target := filepath.Join(filepath.Dir(videoPath), file.Name())
			fmt.Println("converting", convert_target)
			err = convertVideo(convert_target)
			if err != nil {
				fmt.Println("failed to convert file:", videoPath)
				continue
			}
		}
	} else {
		err = convertVideo(videoPath)
		if err != nil {
			return fmt.Errorf("failed to convert file, err=%s", videoPath)
		}
	}

	// successfully convert .ts to .mp4, delete .ts
	successToDelete = true
	return
}

func convertVideo(videoPath string) (err error) {
	video, err := os.Open(videoPath)
	defer video.Close()

	if err != nil {
		// handle the error and return
		return fmt.Errorf("failed to open video, err=%s", err)
	}

	return RunFFmpeg(videoPath)

	// This returns an *os.FileInfo type
	// fileInfo, err := video.Stat()
	// if err != nil {
	// 	// error handling
	// 	return fmt.Errorf("failed to open video stat, err=%s", err)
	// }

	// IsDir is short for fileInfo.Mode().IsDir()
	// if fileInfo.IsDir() {
	// 	// file is a directory, read file names
	// 	fmt.Printf("%s is directory\n", videoPath)
	// 	return
	// }

}

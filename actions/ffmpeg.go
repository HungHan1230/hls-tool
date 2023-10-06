package actions

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunFFmpeg(videoPath string) (err error) {
	videoExtension := filepath.Ext(videoPath)
	videoMP4Name := strings.Replace(videoPath, videoExtension, ".mp4", -1)

	// fmt.Printf("Run:\nvideoPath: %s\nvideoExtension: %s\nvideoMP4Name: %s\n", videoPath, videoExtension, videoMP4Name)
	
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-acodec", "copy", "-vcodec", "copy", videoMP4Name)

	fmt.Println("cmd: ", cmd.Args)
	// 捕獲標準輸出和錯誤輸出
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute ffmpeg: %s\n", err)
		fmt.Printf("Command output:\n%s\n", output)
		return
	}

	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Printf("failed to execute ffmpeg, err=%s\n", err)
	// }
	return
}

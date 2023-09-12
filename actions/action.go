package actions

import (
	"log"

	"github.com/urfave/cli/v2"
)

type ACTION_TYPE string

const (
	// download: (mp4 or m3u8)
	DOWNLOAD = "download"
	// convert: (m3u8 -> mp4)
	CONVERT = "convert"
)

func GetCmdFlags(at ACTION_TYPE) []cli.Flag {
	switch at {
	case DOWNLOAD:
		return downloadCmdFlags()
	case CONVERT:
		return convertCmdFlags()
	default:
		log.Println("unsupported type")
		return nil
	}

}

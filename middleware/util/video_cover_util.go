package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"strings"
)

//GetSnapshot 获取视频封面
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (error, string) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(videoPath).
		Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err, ""
	}
	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err, ""
	}

	//获取视频名称
	names := strings.Split(videoPath, "/")
	coverName := strings.Split(names[len(names)-1], ".")[0]
	coverPath := snapshotPath + coverName + ".jpeg"

	err = imaging.Save(img, coverPath)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return err, ""
	}

	return nil, coverPath
}

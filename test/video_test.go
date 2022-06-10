package test

import (
	"fmt"
	"go-douyin-demo/middleware/util"
	"log"
	"testing"
)

func TestGetVideoSnapshot(t *testing.T) {
	err, path := util.GetSnapshot("D:/go/go-douyin-demo/public/video_data/独立宣言.mp4",
		"D:/go/go-douyin-demo/public/video_cover/", 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path)
}

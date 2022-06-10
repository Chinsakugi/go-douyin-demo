package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-douyin-demo/middleware/util"
	"go-douyin-demo/store"
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

func TestGetVideoList(t *testing.T) {
	err, list := store.GetVideoList(4)
	if err != nil {
		t.Fatal(err)
	}
	res, _ := json.Marshal(list)
	var out bytes.Buffer
	json.Indent(&out, res, "", "\t")
	fmt.Printf("%v", out.String())
}

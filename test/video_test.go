package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwtHelper "go-douyin-demo/middleware/jwt"
	"go-douyin-demo/middleware/util"
	"go-douyin-demo/store"
	"log"
	"strconv"
	"testing"
	"time"
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

func TestTest(t *testing.T) {
	n, err := strconv.Atoi("")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
}

func TestFormat(t *testing.T) {
	latestTime := 1654851182
	formatTime := time.Unix(int64(latestTime), 0).Format("2006-01-02 15-04-05")
	fmt.Println(formatTime)
}

func TestGetVideoFeed(t *testing.T) {
	err, list := store.GetVideoFeed(1654851182)
	if err != nil {
		log.Fatal(err)
	}
	res, _ := json.Marshal(list)
	var out bytes.Buffer
	json.Indent(&out, res, "", "\t")
	fmt.Printf("%v", out.String())
}

func TestGetNow(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func TestGenMyToken(t *testing.T) {
	token := jwtHelper.GenMyToken("czy")
	fmt.Println(token)
}

func TestParseMytoken(t *testing.T) {
	res := jwtHelper.ParseMyToken("czy_1654919179")
	fmt.Println(res)
}

func TestGetFavoriteVideoList(t *testing.T) {
	list := store.GetFavoriteVideoList(5)
	res, _ := json.Marshal(list)
	var out bytes.Buffer
	json.Indent(&out, res, "", "\t")
	fmt.Printf("%v", out.String())
}

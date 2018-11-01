## install 

```
go get github.com/yale8848/gffmpeg

```

## sample

```go

func getFF() GFFmpeg {

	g,e:=NewGFFmpeg(BIN_PATH)
	if e!=nil {
		panic(e)
	}
	return g

}
func Test_CutVideo(t *testing.T)  {
	ff:=getFF().SetDebug(true)
	bd1:=NewBuilder().SrcPath(TEST_MP4_PATH).KeyInt(25).
		CutVideoStartTime(2).CutVideoEndTime(5).
		CutVideo().DistPath(TEST_DIST_CUT_MP4_PATH)
	ff.Set(bd1).Start(nil)
}
func Test_Thumb(t *testing.T)  {
	ff:=getFF().SetDebug(true)
	bd1:=NewBuilder().SrcPath(TEST_MP4_PATH).ThumbStartTime(2).Thumb().DistPath(TEST_DIST_THUMB_PATH)
	ff.Set(bd1).Start(nil)
}
func TestFfmpeg_GetMediaInfo(t *testing.T) {
	ff:=getFF()
	ff.SetDebug(true)

	bd1:=NewBuilder().SrcPath(TEST_MP4_PATH)
	info:=ff.Set(bd1).GetMediaInfo()
	if info!=nil {
		bd:= NewBuilder().SrcPath(TEST_MP4_PATH).
			BitRate(info.BitRate/2).BufSize(info.BitRate/2).
			KeyInt(25).
			Threads(50).
			DistPath(TEST_DIST_MP4_PATH)
		finish:=ff.Set(bd).Start(nil)
		fmt.Printf("%+v\n", finish.CostDuration)
	}
}


```

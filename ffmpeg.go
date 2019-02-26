// Create by Yale 2018/10/31 15:33
package gffmpeg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MediaInfo struct {
	Duration string
	StartTime string
	BitRate int
	VideoCodeFormat string
	VideoFormat string
	VideoResolution string
	AudioFormat string
	AudioFrequencySampling int
}

type GFFmpeg interface {
	GetMediaInfo()*MediaInfo
	Set(builder Builder)GFFmpeg
	SetDebug(debug bool)GFFmpeg
	Start(result chan *CmdFinish)*CmdFinish
    Run(result chan *CmdFinish,args []string)*CmdFinish
}
type CmdFinish struct {
	StdOut bytes.Buffer
	StdErr bytes.Buffer
	Err error
	CostDuration time.Duration
}
type ffmpeg struct {
	binPath string
	buid Builder
	isDebug bool
}

func NewGFFmpeg(binPath string) (GFFmpeg,error) {
	_,err:=os.Stat(binPath)
	if err!=nil {
		return nil,err
	}
	return  &ffmpeg{binPath:binPath},nil
}

func  (ff *ffmpeg) SetDebug(debug bool) GFFmpeg{
	ff.isDebug = debug
	return ff
}
func  (ff *ffmpeg)Set(bd Builder) GFFmpeg{
	ff.buid = bd
	return ff
}
func  (ff *ffmpeg)formatMediaInfo(med *MediaInfo){
	med.Duration = strings.Split(med.Duration,".")[0]
	med.StartTime = strings.Split(med.StartTime,".")[0]
	med.VideoCodeFormat = strings.Split(med.VideoCodeFormat," ")[0]
	med.VideoFormat = strings.Split(med.VideoFormat," ")[0]
	med.AudioFormat = strings.Split(med.AudioFormat," ")[0]
}
func  (ff *ffmpeg)GetMediaInfo()*MediaInfo{
	media:=MediaInfo{}
	if ff.buid == nil{
		return &media
	}

	res:= ff.Run(nil,ff.buid.Build())
	regRes:=ff.findStdByRegexp("Duration: (.*?), start: (.*?), bitrate: (\\d*) kb\\/s",res)
	if len(regRes)==4 {
		media.Duration = regRes[1]
		media.StartTime = regRes[2]
		br,err:=strconv.Atoi(regRes[3])
		if err==nil {
			media.BitRate=br
		}
	}
	regRes=ff.findStdByRegexp("Video: (.*?), (.*?), (.*?)[,\\s]",res)
	if len(regRes)==4 {
		media.VideoCodeFormat = regRes[1]
		media.VideoFormat = regRes[2]
		media.VideoResolution = regRes[3]
	}
	regRes=ff.findStdByRegexp("Audio: (.+), (\\d+) Hz",res)

	if len(regRes)==3 {
		media.AudioFormat = regRes[1]
		afs,err:=strconv.Atoi(regRes[2])
		if err==nil {
			media.AudioFrequencySampling=afs
		}
	}
	ff.formatMediaInfo(&media)
	return  &media
}
func (ff *ffmpeg)findStdByRegexp(regStr string,res *CmdFinish) []string{

	if res==nil {
		return nil
	}
	var ts []string
	reg := regexp.MustCompile(regStr)
	if len(res.StdOut.String()) >0 {
		ts=reg.FindStringSubmatch(res.StdOut.String())
		if ts!=nil && len(ts)>1 {
			return ts
		}
	}else if len(res.StdErr.String())>0 {
		ts=reg.FindStringSubmatch(res.StdErr.String())
		if ts!=nil && len(ts)>1 {
			return ts
		}
	}
	return nil
}

func (ff *ffmpeg)Start(result chan *CmdFinish)*CmdFinish{
	if ff.buid == nil {
		if result!=nil {
			result <- nil
		}
		return nil
	}
	return ff.Run(result,ff.buid.Build())
}
func (ff *ffmpeg)Run(result chan *CmdFinish,args []string)*CmdFinish  {
	if ff.isDebug {
		fmt.Printf("%s %s\r\n",ff.binPath,strings.Join(args," "))
	}
	startTime:=time.Now()

	cmd:=exec.Command(ff.binPath,args...)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	cmd.Stdout = &out
	err:=cmd.Run()
	endTime:=time.Now()

	res:=CmdFinish{StdOut:out,StdErr:outErr,Err:err,CostDuration:endTime.Sub(startTime)}
	if result!=nil {
		result <-&res
	}
	if ff.isDebug {
		if len(res.StdErr.String())>0 {
			fmt.Println(res.StdErr.String())
		} else if len(res.StdOut.String()) >0 {
			fmt.Println(res.StdOut.String())
		}
	}
	return &res
}
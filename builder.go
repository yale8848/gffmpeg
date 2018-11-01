// Create by Yale 2018/11/1 14:40
package gffmpeg

import (
	"fmt"
	"strconv"
)

type Builder interface {
	SrcPath(srcPath string) Builder
	DistPath(distPath string) Builder
	KeyInt(keyInt int) Builder
	BitRate(rate int ) Builder
	BufSize(size int) Builder
	Threads(threadsNum int)Builder
	CutVideoStartTime(startTime int) Builder
	CutVideoEndTime(endTime int) Builder
	CutVideo()Builder
	ThumbStartTime(startTime int) Builder
	ThumbResolution(resolution string) Builder
	Thumb() Builder
	Build()([]string)
}

type FFBuilder struct {
	cmds []string
}

func NewBuilder() Builder  {
	return &FFBuilder{cmds:make([]string,0)}
}
func (bd *FFBuilder)CutVideo() Builder{
	bd.addCmds("-codec","copy")
	return bd
}
func (bd *FFBuilder)CutVideoStartTime(startTime int) Builder{
	if startTime >0 {
		bd.addCmds("-ss",strconv.Itoa(startTime))
	}
	return bd
}
func (bd *FFBuilder)CutVideoEndTime(endTime int) Builder{
	if endTime >0 {
		bd.addCmds("-t",strconv.Itoa(endTime))
	}
	return bd
}
func (bd *FFBuilder)Thumb() Builder{
	bd.addCmds("-f","image2")
	return bd
}
func (bd *FFBuilder)ThumbResolution(resolution string)Builder{

	if len(resolution)>0 {
		bd.addCmds("-s",resolution)
	}
	return bd
}
func (bd *FFBuilder)ThumbStartTime(startTime int) Builder{

	if startTime>0 {
		bd.addCmds("-t",strconv.Itoa(startTime))
	}
	return bd
}
func (bd *FFBuilder)Build()([]string){
	return bd.cmds
}
func (bd *FFBuilder)addCmds(cmds ...string){
	for _,v:=range cmds{
		bd.cmds = append(bd.cmds,v)
	}
}
func (bd *FFBuilder)BitRate(rate int ) Builder{
	if rate >0 {
		bd.addCmds("-b:v",fmt.Sprintf("%dk",rate))
	}
	return bd
}
func (bd *FFBuilder)BufSize(size int)  Builder{
	if size >0 {
		bd.addCmds("-bufsize",fmt.Sprintf("%dk",size))
	}
	return bd
}
func (bd *FFBuilder)Threads(threadsNum int)Builder{
	if threadsNum>0 {
		bd.addCmds("-threads",strconv.Itoa(threadsNum))
	}
	return bd
}
func (bd *FFBuilder)SrcPath(srcPath string) Builder {
	if len(srcPath)>0 {
		bd.addCmds("-i",srcPath)
	}
	return bd
}
func (bd *FFBuilder)DistPath(distPath string) Builder {
	if len(distPath)>0 {
		bd.addCmds("-y",distPath)
	}
	return bd
}
func (bd *FFBuilder)KeyInt(keyInt int) Builder {
	if keyInt>0 {
		bd.addCmds("-x264opts",fmt.Sprintf("keyint=%d",keyInt))
	}
	return bd
}

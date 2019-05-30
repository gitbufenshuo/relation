package content

import (
	"bufio"
	"fmt"
	"os"

	"strings"
	"sync"
	"time"

	"github.com/gitbufenshuo/relation/tool"
)

var FILENAME = "content.gob"
var lo sync.Mutex

type Two struct {
	UP   uint64
	DOWN []uint64
}
type Content struct {
	LastSync int64
	NowIdx   uint64
	WriteOK  bool
	//////////////
	All map[uint64]*Two
}

////////
func (con *Content) CheckEx(uid uint64) bool {
	if len(con.All) == 0 {
		return false
	}
	//////////////////////
	if _, ok := con.All[uid]; ok {
		return true
	}
	return false
}
func (con *Content) Add(self, up uint64) bool {
	if con.CheckEx(self) {
		return false
	}
	if !con.CheckEx(up) {
		return false
	}
	////////
	con.NowIdx++
	t := new(Two)
	t.UP = up
	con.All[self] = t
	downlist := con.All[up].DOWN
	downlist = append(downlist, self)
	{
		if time.Now().Unix()-g_content.LastSync >= 500 && con.WriteOK {
			// should sync
			Content_Record(self, up, con.NowIdx)
		}
	}
	return true
}

////////
var g_content *Content

func (con *Content) Init() {
	con.All = make(map[uint64]*Two)
	con.All[0] = new(Two)
	con.All[0].UP = 0

}

func Content_Init() {
	g_content = new(Content)
	g_content.Init()
	// 从 file 里 load 回来
	file, err := os.OpenFile(FILENAME, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		g_content.WriteOK = true
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		segs := strings.Split(txt, " ")
		self := tool.StringToUint64(segs[0])
		up := tool.StringToUint64(segs[1])
		idx := tool.StringToUint64(segs[2])
		if idx <= g_content.NowIdx {
			continue
		}
		g_content.Add(self, up)
	}
	g_content.WriteOK = true
	file.Close()
}

// save 到文件
func Content_Record(self, up, nowidx uint64) {
	// fmt.Println("----", tool.NowMs())
beginf:
	f, err := os.OpenFile(FILENAME, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		f1, _ := os.Create(FILENAME)
		f1.Close()
		goto beginf
	}
	f.WriteString(fmt.Sprintf("%v %v %v\n", self, up, nowidx))
	f.Close()
	// fmt.Println(tool.NowMs())
}

func Content_lock() {
	lo.Lock()
}
func Content_unlock() {
	lo.Unlock()
}
func Content_add(self, up uint64) bool {
	return g_content.Add(self, up)
}
func Content_getallup(self uint64) []uint64 {
	if !g_content.CheckEx(self) {
		return nil
	}
	res := []uint64{}
	nows := self
	for {
		if v, ok := g_content.All[nows]; ok {
			if v.UP == 0 {
				break
			}
			res = append(res, v.UP)
			nows = v.UP
		}
	}
	return res
}
func Content_getoneup(self uint64) uint64 {
	if !g_content.CheckEx(self) {
		return 0
	}
	if v, ok := g_content.All[self]; ok {
		return v.UP
	}
	return 0
}

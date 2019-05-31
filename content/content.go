package content

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"strings"
	"sync"
	"time"

	"github.com/gitbufenshuo/relation/tool"
)

var (
	MaxLevel uint8 = 3
	JuLevel  uint8 = 2
)

var FILENAME = "content.gob"
var lo sync.Mutex

type Two struct {
	UP    uint64
	Level uint8
	BigJu uint8 // 大菊线个数
	DOWN  []uint64
}
type Content struct {
	LastSync int64
	NowIdx   uint64
	WriteOK  bool
	MaxDepth uint64 // 最大深度
	MaxWidth uint64 // 最大宽度
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
func (con *Content) CheckMaxLevel(uid uint64) bool {
	if len(con.All) == 0 {
		return false
	}
	/////////
	if v, ok := con.All[uid]; ok {
		return v.Level < MaxLevel
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
	con.All[up].DOWN = downlist
	{
		if time.Now().Unix()-g_content.LastSync >= 500 && con.WriteOK {
			// should sync
			Content_Record(self, up, con.NowIdx)
		}
	}
	return true
}

func (con *Content) LevelUp(self uint64) bool {
	if !con.CheckEx(self) {
		return false
	}
	if !con.CheckMaxLevel(self) {
		return false
	}
	con.NowIdx++
	v := con.All[self]
	v.Level++
	{
		// dajuxian
		if v.Level == JuLevel && v.BigJu == 0 { // 我已成橘, 且我无线
			if v.UP != 0 {
				_up := con.All[v.UP]
				_up.BigJu++
				if _up.BigJu == 0 {
					_up.BigJu--
				}
				for {
					if _up.BigJu == 1 && _up.UP != 0 {
						_up = con.All[_up.UP]
						_up.BigJu++
						if _up.BigJu == 0 {
							_up.BigJu--
						}
					} else {
						break
					}
				}
			}
		}
	}
	{
		if time.Now().Unix()-g_content.LastSync >= 500 && con.WriteOK {
			// should sync
			// up 值如果是math.MAXUINT64，则说明，这一条消息是 升级，不是注册
			Content_Record(self, math.MaxUint64, con.NowIdx)
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
	g_content = nil
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
		if up == math.MaxUint64 {
			g_content.LevelUp(self)
		} else {
			g_content.Add(self, up)
		}
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

func Content_flushall() {
	os.Remove(FILENAME)
	Content_Init()
}

func Content_sta() (uint64, uint64, uint64) {
	return uint64(len(g_content.All)), g_content.MaxDepth, g_content.MaxWidth
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

func Content_levelup(self uint64) bool {
	return g_content.LevelUp(self)
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
func Content_getalldzj(self uint64) []uint64 {
	if !g_content.CheckEx(self) {
		return nil
	}
	v := g_content.All[self]
	return []uint64{uint64(v.Level), uint64(len(v.DOWN)), uint64(v.BigJu)}
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

func Content_getbigju(self uint64) uint64 {
	if !g_content.CheckEx(self) {
		return 0
	}
	if v, ok := g_content.All[self]; ok {
		return uint64(v.BigJu)
	}
	return 0
}
func Content_getlevel(self uint64) uint64 {
	if !g_content.CheckEx(self) {
		return 0
	}
	if v, ok := g_content.All[self]; ok {
		return uint64(v.Level)
	}
	return 0
}
func Content_getzhixia(self uint64) uint64 {
	if !g_content.CheckEx(self) {
		return 0
	}
	if v, ok := g_content.All[self]; ok {
		return uint64(len(v.DOWN))
	}
	return 0
}

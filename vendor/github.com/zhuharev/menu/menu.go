package menu

import (
	"github.com/Unknwon/com"
	"github.com/ungerik/go-dry"
	"regexp"
)

var (
	lineRe = regexp.MustCompile(`[ ]{0,2}([\d]{1,2}) \[([\S]{1,50})\]\(([\S]{1,50})\)`)
)

type Menus []Menu

type Menu struct {
	Title string
	Link  string
	Order int
	Items Menus
}

func (menus Menus) Len() int           { return len(menus) }
func (menus Menus) Less(i, j int) bool { return menus[i].Order < menus[j].Order }
func (menus Menus) Swap(i, j int)      { menus[i], menus[j] = menus[j], menus[i] }

func NewFromFile(fpath string) (Menus, error) {
	var (
		res Menus
	)
	lines, e := dry.FileGetNonEmptyLines(fpath)
	if e != nil {
		return nil, e
	}
	for _, line := range lines {
		var m Menu
		arr := lineRe.FindStringSubmatch(line)
		if len(arr) != 4 {
			continue
		}
		m.Order = com.StrTo(arr[1]).MustInt()
		m.Title = arr[2]
		m.Link = arr[3]
		if line[0] == ' ' && line[1] == ' ' {
			parent := res[len(res)-1]
			m.Link = parent.Link + m.Link
			res[len(res)-1].Items = append(res[len(res)-1].Items, m)
		} else {
			res = append(res, m)
		}
	}
	return res, nil
}

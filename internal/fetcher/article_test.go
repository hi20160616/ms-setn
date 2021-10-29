package fetcher

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hi20160616/exhtml"
	"github.com/pkg/errors"
)

func TestFetchTitle(t *testing.T) {
	tests := []struct {
		url   string
		title string
	}{
		{"https://www.setn.com/News.aspx?NewsID=1014270", "激瘦秘密就是PP塑崩褲！《驕女》春滿大尺寸也穿得下 | 生活"},
		{"https://www.setn.com/News.aspx?NewsID=1018311", "五倍券消費滿額送好禮sodastream氣泡水機超划算 | 生活"},
	}
	for _, tc := range tests {
		a := NewArticle()
		u, err := url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		a.U = u
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		got, err := a.fetchTitle()
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			} else {
				fmt.Println("ignore pass test: ", tc.url)
			}
		} else {
			if tc.title != got {
				t.Errorf("\nwant: %s\n got: %s", tc.title, got)
			}
		}
	}

}

func TestFetchUpdateTime(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://www.setn.com/News.aspx?NewsID=1014270", "2021-10-19 17:00:00 +0800 UTC"},
		{"https://www.setn.com/News.aspx?NewsID=1018311", "2021-10-28 13:00:00 +0800 UTC"},
	}
	var err error
	for _, tc := range tests {
		a := NewArticle()
		a.U, err = url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		tt, err := a.fetchUpdateTime()
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			}
		}
		ttt := tt.AsTime()
		got := shanghai(ttt)
		if got.String() != tc.want {
			t.Errorf("\nwant: %s\n got: %s", tc.want, got.String())
		}
	}
}

func TestFetchContent(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://www.setn.com/News.aspx?NewsID=1014270", "激瘦秘密就是PP塑崩褲！《驕女》春滿大尺寸也穿得下 | 生活"},
		{"https://www.setn.com/News.aspx?NewsID=1018311", "五倍券消費滿額送好禮sodastream氣泡水機超划算 | 生活"},
	}
	var err error

	for _, tc := range tests {
		a := NewArticle()
		a.U, err = url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		c, err := a.fetchContent()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(c)
	}
}

func TestFetchArticle(t *testing.T) {
	tests := []struct {
		url string
		err error
	}{
		{"https://www.setn.com/News.aspx?NewsID=1014270", ErrTimeOverDays},
		{"https://www.setn.com/News.aspx?NewsID=1018311", nil},
	}
	for _, tc := range tests {
		a := NewArticle()
		a, err := a.fetchArticle(tc.url)
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			} else {
				fmt.Println("ignore old news pass test: ", tc.url)
			}
		} else {
			fmt.Println("pass test: ", a.Content)
		}
	}
}

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func TestFindEntries(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.String())
		if r.URL.String() == "/" {
			w.Write([]byte(`
			<table summary="作家データ">
			<tr><td class="header">分類：</td><td>著者</td></tr>
			<tr><td class="header">作家名：</td><td>テスト　太郎</td></tr>
			<tr><td class="header">作家名読み：</td><td>テスト　タロウ</td></tr>
			<tr><td class="header">ローマ字表記：</td><td>Test, Taro</td></tr>
			</table>
			<ol>
			<li><a href="../cards/999999/card001.html">テスト書籍001</a></li>
			<li><a href="../cards/999999/card002.html">テスト書籍002</a></li>
			<li><a href="../cards/999999/card003.html">テスト書籍003</a></li>
			</ol>
			`))
		} else {
			pat := regexp.MustCompile(`.*/cards/([0-9]+)/card([0-9]+).html$`)
			token := pat.FindStringSubmatch(r.URL.String())
			w.Write([]byte(fmt.Sprintf(`
				<table summary="作家データ">
				<tr><td class="header">分類：</td><td>著者</td></tr>
				<tr><td class="header">作家名：</td><td>テスト　太郎</td></tr>
				<tr><td class="header">作家名読み：</td><td>テスト　タロウ</td></tr>
				<tr><td class="header">ローマ字表記：</td><td>Test, Taro</td></tr>
				</table>
				<table border="1" summary="ダウンロードデータ" class="download">
				<tr>
					<td><a href="./files/%[1]s_%[2]s.zip">%[1]s_%[2]s.zip</a></td>
				</tr>
				</table>
				`, token[1], token[2])))
		}
	}))
	defer ts.Close()

	tmp := pageURLFormat
	pageURLFormat = ts.URL + "/cards/%s/card%s.html"
	defer func() {
		pageURLFormat = tmp
	}()

	got, err := findEntries(ts.URL)
	if err != nil {
		t.Error(err)
		return
	}

	want := []Entry{
		{
			AuthorID: "999999",
			Author:   "テスト　太郎",
			TitleID:  "001",
			Title:    "テスト書籍001",
			SiteURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_001.zip",
		},
		{
			AuthorID: "999999",
			Author:   "テスト　太郎",
			TitleID:  "002",
			Title:    "テスト書籍002",
			SiteURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_002.zip",
		},
		{
			AuthorID: "999999",
			Author:   "テスト　太郎",
			TitleID:  "003",
			Title:    "テスト書籍003",
			SiteURL:  ts.URL,
			ZipURL:   ts.URL + "/cards/999999/files/999999_003.zip",
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %+v, but got %+v", want, got)
	}
}

func TextExtractText(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir(".")))
	defer ts.Close()

	got, err := extractText(ts.URL + "/testdata/example.zip")
	if err != nil {
		t.Fatal(err)
		return
	}

	want := "テストデータ\n"
	if want != got {
		t.Errorf("want %+v, but got %+v", want, got)
	}
}

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/YuChaoGithub/CHU-Ing-Wen/converter"
)

const (
	bodyMaxLen = 4269

	requestErrStr     = "（４００）泥ㄉ要求婐做ㄅ到ㄌ！！！＞／／／＜"
	internalErrStr    = "（５００）ㄜ……婐ㄉ內部粗ㄌ１點問題０.０"
	bodyTooLongErrStr = "（４００）ㄚ！泥給ㄉ文章太長ㄌ啦＝　＝"
)

func getAnthologyList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(converter.GetAnthologyList()))
}

func getFromAnthology(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, internalErrStr, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(converter.GetFromAnthology(string(body))))
}

func convertArticle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, internalErrStr, http.StatusInternalServerError)
		return
	}

	if utf8.RuneCountInString(string(body)) > bodyMaxLen {
		http.Error(w, bodyTooLongErrStr, http.StatusBadRequest)
		return
	}

	w.Write([]byte(converter.Convert(string(body))))
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request to", r.URL, "from", r.RemoteAddr)

		next.ServeHTTP(w, r)
	})
}

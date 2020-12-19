package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
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

func addCORSHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "420420")
		next.ServeHTTP(w, r)
	})
}

func custom404(fs http.FileSystem) http.Handler {
	fsrv := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			w.Header().Add("Location", "https://knowyourrights.github.io/")
			w.WriteHeader(http.StatusSeeOther)
			return
		}
		fsrv.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				log.Printf("%s\n%s\n", fmt.Errorf("%s", err), debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

package senstive

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/syyongx/go-wordsfilter"
)

type SenstiveFilter struct {
	filter *wordsfilter.WordsFilter
	root   map[string]*wordsfilter.Node
}

func NewSenstiveFilter(texts []string) *SenstiveFilter {
	wf := wordsfilter.New()
	root := wf.Generate(texts)
	return &SenstiveFilter{
		filter: wf,
		root:   root,
	}
}

func NewSenstiveFilterWithFile(path string) *SenstiveFilter {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	texts := strings.Split(string(fileContent), ",")
	wf := wordsfilter.New()
	root := wf.Generate(texts)
	if err != nil {
		log.Fatal(err)
	}
	return &SenstiveFilter{
		filter: wf,
		root:   root,
	}
}

func (f *SenstiveFilter) Contains(word string) bool {
	words := []rune(word)
	wordLen := len(words)
	for i := 0; i < wordLen; i++ {
		for j := i + 1; j <= wordLen; j++ {
			sub := string(words[i:j])
			if f.filter.Contains(sub, f.root) {
				return true
			}
		}
	}
	return false
}

func (f *SenstiveFilter) Remove(word string) {
	f.filter.Remove(word, f.root)
}

func (f *SenstiveFilter) Replace(word string) string {
	return f.filter.Replace(word, f.root)
}

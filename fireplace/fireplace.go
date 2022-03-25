package fireplace

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

type Handler struct {
	files map[string]string
}

type ReplaceOpts struct {
	LineNumber uint64
	ReplaceTo  string
	changed    bool
}

func New() Handler {
	return Handler{
		files: make(map[string]string),
	}
}

func (h *Handler) Add(fileID, path string) error {
	h.files[fileID] = path
	return nil
}

type SearchOpts struct {
	InFileID string
}

type SearchResult struct {
	FileID      string
	LineNumber  uint64
	LineContent string
}

// Search returns
func (h *Handler) Search(s string, opts ...SearchOpts) ([]SearchResult, error) {
	var fileID string
	if len(opts) == 1 {
		fileID = opts[0].InFileID
	}

	if fileID == "" {
		var results []SearchResult
		for id, file := range h.files {
			searchResults := h.searchFile(file, s, id)
			if searchResults != nil {
				results = append(results, searchResults...)
			}
		}
		if len(results) == 0 {
			return nil, errors.New("content not found")
		}
		return results, nil
	}
	return h.searchFile(h.files[fileID], s, fileID), nil
}

func (h *Handler) Replace(fileID string, opts *ReplaceOpts) error {
	// TODO
	return nil
}

func (h *Handler) ReplaceMultiple(changes map[string][]*ReplaceOpts) error {
	for id, path := range h.files {
		if replaces, ok := changes[id]; ok {
			if err := h.replaceMultipleInFile(path, replaces); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Handler) replaceMultipleInFile(path string, replaces []*ReplaceOpts) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// create temp file
	tmp, err := os.Create(file.Name() + "_updated")
	if err != nil {
		log.Fatal(err)
	}
	defer tmp.Close()

	smallest := h.smallestOfReplaces(replaces)
	if smallest == nil {
		return errors.New("no smallest")
	}
	var lineNum uint64 = 1
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		if smallest != nil && lineNum == smallest.LineNumber {
			if _, err := io.WriteString(tmp, smallest.ReplaceTo+"\n"); err != nil {
				return err
			}
			smallest.changed = true
			smallest = h.smallestOfReplaces(replaces)
		} else {
			line := sc.Text()
			if _, err := io.WriteString(tmp, line+"\n"); err != nil {
				return err
			}
		}

		lineNum++
	}

	return nil
}

func (h Handler) searchFile(path, searchFor, fileID string) []SearchResult {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()
	return h.searchScanner(f, searchFor, fileID)
}

func (h Handler) searchScanner(f *os.File, s, fileID string) []SearchResult {
	scanner := bufio.NewScanner(f)
	var counter uint64 = 1
	var line []byte
	sb := []byte(s)
	var res []SearchResult
	for scanner.Scan() {
		line = scanner.Bytes()
		if bytes.Contains(line, sb) {
			res = append(res, SearchResult{
				FileID:      fileID,
				LineContent: string(line),
				LineNumber:  counter,
			})
		}
		counter++
	}

	return res
}

func (h Handler) smallestOfReplaces(replaces []*ReplaceOpts) *ReplaceOpts {
	var smallest uint64 = math.MaxUint64
	var result *ReplaceOpts = nil

	for _, replace := range replaces {
		if !replace.changed && replace.LineNumber < smallest {
			smallest = replace.LineNumber
			result = replace
		}
	}
	return result
}

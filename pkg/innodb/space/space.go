package space

import (
	"os"
	"io"
	"fmt"
	"goParseInnodb/pkg/innodb/parse"
)

const PAGE_SIZE = 16384

type PageWrapper struct {
	Number int64
	Page   any
	Err    error
}

type Space struct {
	Path  string
	Size  int64
	Pages int64
}

func OpenSpace(path string) (*Space, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	size := stat.Size()

	return &Space{
		Path:  path,
		Size:  size,
		Pages: size / PAGE_SIZE,
	}, nil
}

func (s *Space) OpenPage(pageNumber int64) (any, error) {
	offset := pageNumber * PAGE_SIZE
	if offset < 0 || offset+PAGE_SIZE > s.Size {
		return nil, fmt.Errorf("invalid page number %d", pageNumber)
	}

	f, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, PAGE_SIZE)
	_, err = f.ReadAt(buf, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	pg, err := parse.ParsePage(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page %d: %w", pageNumber, err)
	}

	return pg, nil
}


func (s *Space) IteratePages() <-chan PageWrapper {
	ch := make(chan PageWrapper)

	go func() {
		defer close(ch)

		for pageNumber := int64(0); pageNumber < s.Pages; pageNumber++ {
			p, err := s.OpenPage(pageNumber)
			if err != nil {
				ch <- PageWrapper{Number: pageNumber, Err: err}
				return
			}
			if p != nil {
				ch <- PageWrapper{Number: pageNumber, Page: p}
			}
		}
	}()

	return ch
}

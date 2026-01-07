package innodb

import (
	"os"
	"io"
	"fmt"
)

const PAGE_SIZE = 16384

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

func (s *Space) OpenPage(pageNumber int64) (*Page, error) {
	offset := pageNumber * PAGE_SIZE

	if offset < 0 || offset+PAGE_SIZE > s.Size {
		return nil, nil
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
	page, err := NewPage(buf)
	if err != nil {
		return nil, fmt.Errorf("page parse failed: %w", err)
	}
	return page, nil
}

func (s *Space) IteratePages() <-chan PageWrapper {
	ch := make(chan PageWrapper)

	go func() {
		defer close(ch)

		for pageNumber := int64(0); pageNumber < s.Pages; pageNumber++ {
			page, err := s.OpenPage(pageNumber)
			if err != nil {
				ch <- PageWrapper{Err: err}
				return
			}
			if page != nil {
				ch <- PageWrapper{Number: pageNumber, Page: page}
			}
		}
	}()

	return ch
}

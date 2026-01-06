package innodb

import (
	"os"
	"io"
)

const PageSize = 16384

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
		Pages: size / PageSize,
	}, nil
}

func (s *Space) OpenPage(pageNumber int64) (*Page, error) {
	offset := pageNumber * PageSize

	if offset < 0 || offset+PageSize > s.Size {
		return nil, nil
	}

	f, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, PageSize)

	_, err = f.ReadAt(buf, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return NewPage(buf), nil
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

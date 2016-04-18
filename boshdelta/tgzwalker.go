package boshdelta

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

type tgzWalker func(h *tar.Header, r *tar.Reader) error

func tgzWalk(tgz io.Reader, walkFn tgzWalker) error {
	gzf, err := gzip.NewReader(tgz)
	if err != nil {
		return err
	}
	tr := tar.NewReader(gzf)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err = walkFn(header, tr)
		if err != nil {
			return err
		}
	}
	return nil
}

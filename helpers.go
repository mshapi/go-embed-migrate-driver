package embed_migrate_driver

import (
	"embed"
	"errors"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"io"
	"io/fs"
	"sync"
)

var (
	fsPath   = map[string]*embed.FS{}
	fsPathMu sync.RWMutex
)

func Register(path string, src *embed.FS) {
	fsPathMu.Lock()
	fsPath[path] = src
	fsPathMu.Unlock()
}

func getFS(path string) *embed.FS {
	fsPathMu.RLock()
	defer fsPathMu.RUnlock()
	return fsPath[path]
}

func makeBindata(eFS *embed.FS) (*bindata.Bindata, error) {
	if eFS == nil {
		return nil, errors.New("empty embed data")
	}

	names := make([]string, 0)
	index := map[string]fs.File{}

	err := fs.WalkDir(eFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		f, e := eFS.Open(path)
		if e != nil {
			return e
		}
		i, e := d.Info()
		if e != nil {
			return e
		}
		name := i.Name()
		names = append(names, name)
		index[name] = f
		return nil
	})
	if err != nil {
		return nil, err
	}

	s := bindata.Resource(names, func(name string) ([]byte, error) {
		f, ok := index[name]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return io.ReadAll(f)
	})

	dd, err := bindata.WithInstance(s)
	if err != nil {
		return nil, err
	}

	switch dd.(type) {
	case *bindata.Bindata:
		return dd.(*bindata.Bindata), nil
	default:
		return nil, errors.New("bindata driver error")
	}
}

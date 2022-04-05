package embed_migrate_driver

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
)

var ErrPathNotRegister = errors.New("path not register")

func init() {
	source.Register("embed", &embedDriver{})
}

type embedDriver struct {
	*bindata.Bindata
}

func (d *embedDriver) Open(url string) (source.Driver, error) {
	eFS := getFS(url)
	if eFS == nil {
		return nil, fmt.Errorf("%w: %s", ErrPathNotRegister, url)
	}
	bd, err := makeBindata(eFS)
	if err != nil {
		return nil, err
	}
	return &embedDriver{bd}, nil
}

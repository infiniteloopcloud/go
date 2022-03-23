package cfg

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Opts struct {
	Path      string
	HotReload bool
	Errorf    func(ctx context.Context, err error, format string, args ...interface{})
	Infof     func(ctx context.Context, format string, args ...interface{})
	Warnf     func(ctx context.Context, format string, args ...interface{})
}

type cfg struct {
	file map[string]string
	opts Opts

	lock *sync.RWMutex
}

type Config interface {
	Get(key string) string
}

func New(opts Opts) (cfg, error) {
	c := cfg{
		file: make(map[string]string),
		opts: opts,
		lock: &sync.RWMutex{},
	}
	if c.opts.Path != "" {
		if err := c.openFile(); err != nil {
			return c, err
		}
		if c.opts.HotReload {
			go c.watcher()
		}
	}
	return c, nil
}

func (c cfg) Get(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if v, ok := c.file[key]; ok {
		return v
	}
	return os.Getenv(key)
}

func (c *cfg) openFile() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var file, err = os.OpenFile(c.opts.Path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c.file)
}

func (c *cfg) watcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		//nolint:gosimple
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := c.openFile(); err != nil {
						c.errorf(err, "openFile in cfg package, fsnotify.Write")
					}
					c.infof("fsnotify file write happened on: %s", c.opts.Path)
				}
			}
		}
	}()

	c.infof("file added to fsnotify watchlist: %s", c.opts.Path)
	err = watcher.Add(c.opts.Path)
	if err != nil {
		c.opts.Errorf(context.Background(), err, "watcher::Add")
	}
	<-done
}

func (c cfg) errorf(err error, format string, args ...interface{}) {
	if c.opts.Errorf != nil {
		c.opts.Errorf(context.Background(), err, format, args...)
	}
}

func (c cfg) infof(format string, args ...interface{}) {
	if c.opts.Infof != nil {
		c.opts.Infof(context.Background(), format, args...)
	}
}

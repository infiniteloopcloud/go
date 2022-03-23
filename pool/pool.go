package pool

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/infiniteloopcloud/log"
	_ "github.com/lib/pq"
)

const (
	defaultDriverName = "postgres"
	ctxTimeout        = 5 * time.Second
	maxTimeout        = 10 * time.Second
)

type Connection interface {
	Get() *sql.DB
}

type Pool struct {
	mutex         *sync.Mutex
	connectionStr string
	connection    *sql.DB
	driverName    string
}

func New(connectionStr string) *Pool {
	pool := &Pool{
		mutex:         &sync.Mutex{},
		connectionStr: connectionStr,
		driverName:    defaultDriverName,
	}
	pool.Get()
	return pool
}

func (p *Pool) SetDriverName(driverName string) {
	p.driverName = driverName
}

func (p *Pool) Get() *sql.DB {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.connection == nil {
		conn, err := sql.Open(p.driverName, p.connectionStr)
		if err != nil {
			log.Error(context.Background(), err, "cockroach connection failed")
			return nil
		}
		p.connection = conn
	}
	return p.connection
}

func (p *Pool) Health(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	var errCh = make(chan error, 1)
	go func() {
		select {
		// Kill the process after a max timeout, to prevent deadlock
		case <-time.After(maxTimeout):
			errCh <- errors.New("overslept")
		case <-timeoutCtx.Done():
			errCh <- timeoutCtx.Err()
		}
	}()

	go func() {
		if p.connection == nil {
			errCh <- errors.New("database connection is nil")
			return
		}
		errCh <- p.connection.PingContext(timeoutCtx)
	}()

	err := <-errCh
	// Database connection became unhealthy, set to nil
	if err != nil {
		p.connection = nil
	}

	// Attempt to reconnect
	p.Get()
	return err
}

func (p *Pool) Close() error {
	return p.connection.Close()
}

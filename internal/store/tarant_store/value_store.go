package tarantstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/tarantool/go-tarantool"
)

const schemaName = "values"

type work struct {
	Bn  store.Bind
	Err error
}

type ValueStore struct {
	db *tarantool.Connection
}

func NewValueStore(db *tarantool.Connection) *ValueStore {
	return &ValueStore{
		db: db,
	}
}

// SetKeys inserts keys into database with help of trasaction
func (s *ValueStore) SetKeys(data map[string]interface{}) ([]store.Bind, error) {
	const op = "valuestore.SetKeys"
	ctx, cancel := context.WithCancel(context.Background())

	// Message channels
	tasks := make(chan *store.Bind, 20)
	erroring := make(chan *work, 20)

	for i := 0; i < len(data) && i < 20; i++ {
		go writeWorker(ctx, s.db, erroring, tasks)
	}

	go mapConvertor(ctx, tasks, data)

	uninserted := []store.Bind{}
	for i := 0; i < len(data); i++ {
		w := <-erroring
		if w != nil {
			if errors.Is(w.Err, store.ErrConnCLosed) {
				cancel()
				return []store.Bind{}, w.Err
			}
			uninserted = append(uninserted, w.Bn)
		}
	}

	return uninserted, nil
}

// TODO:
func (s *ValueStore) GetKeys(keys []string) (map[string]interface{}, []string, error) {
	fmt.Println(keys)
	const op = "valuestore.GeKeys"
	ctx, cancel := context.WithCancel(context.Background())

	// Message channels
	readed := make(chan *work, 20)
	tasks := make(chan string, 20)

	for i := 0; i < len(keys) && i < 20; i++ {
		go readWorker(ctx, s.db, tasks, readed)
	}

	go sliceConvertor(ctx, tasks, keys)

	res := make(map[string]interface{}, len(keys))
	var unfound []string
	for i := 0; i < len(keys); i++ {
		w := <-readed
		if w.Err != nil {
			if errors.Is(w.Err, store.ErrConnCLosed) {
				cancel()
				return map[string]interface{}{}, []string{}, fmt.Errorf("op: %s error: %w", op, w.Err)
			}
			// For future handling
			if errors.Is(w.Err, store.ErrRecordNotFound) {
				unfound = append(unfound, w.Bn.Key)
			} else {
				unfound = append(unfound, w.Bn.Key)
			}
		} else {
			res[w.Bn.Key] = w.Bn.Value
		}
	}

	return res, unfound, nil
}

func writeWorker(ctx context.Context, conn *tarantool.Connection, erroring chan<- *work, tasks <-chan *store.Bind) {
	for {
		select {
		// TODO: Check for closing channel
		case b, ok := <-tasks:
			if !ok {
				return
			}
			_, err := conn.Replace(schemaName, []interface{}{b.Key, b.Value})
			if err != nil {
				if !conn.ConnectedNow() {
					erroring <- &work{
						Err: store.ErrConnCLosed,
					}
				} else {
					erroring <- &work{
						Bn: store.Bind{
							Key:   b.Key,
							Value: b.Value,
						},
						Err: err,
					}
				}
			} else {
				erroring <- nil
			}
		case <-ctx.Done():
			return
		}
	}
}

func readWorker(ctx context.Context, conn *tarantool.Connection, tasks <-chan string, res chan<- *work) {
	for {
		select {
		case b, ok := <-tasks:
			if !ok {
				return
			}
			resp, err := conn.Select(schemaName, "primary", 0, 2, tarantool.IterEq, []interface{}{b})
			if err != nil {
				if !conn.ConnectedNow() {
					res <- &work{
						Err: store.ErrConnCLosed,
					}
				} else {
					res <- &work{
						Err: err,
					}
				}
			} else {
				if len(resp.Tuples()) == 0 {
					fmt.Println("not found", b)
					res <- &work{
						Bn: store.Bind{
							Key: b,
						},
						Err: store.ErrRecordNotFound,
					}
				} else {
					res <- &work{
						Bn: store.Bind{
							Key:   resp.Tuples()[0][0].(string),
							Value: resp.Tuples()[0][1],
						},
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func sliceConvertor(ctx context.Context, result chan<- string, data []string) {
	defer close(result)

	for _, v := range data {
		select {
		case result <- v:
		case <-ctx.Done():
			return
		}
	}
}

func mapConvertor(ctx context.Context, result chan<- *store.Bind, data map[string]interface{}) {
	defer close(result)

	for k, v := range data {
		select {
		case result <- &store.Bind{
			Key:   k,
			Value: v,
		}:
		case <-ctx.Done():
			return
		}
	}
}

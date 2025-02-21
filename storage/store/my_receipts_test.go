package store_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/module/metrics"
	"github.com/onflow/flow-go/storage"
	"github.com/onflow/flow-go/storage/operation/dbtest"
	"github.com/onflow/flow-go/storage/store"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestMyExecutionReceiptsStorage(t *testing.T) {
	withStore := func(t *testing.T, f func(store1 *store.MyExecutionReceipts, db storage.DB)) {
		dbtest.RunWithDB(t, func(t *testing.T, db storage.DB) {
			metrics := metrics.NewNoopCollector()
			results := store.NewExecutionResults(metrics, db)
			receipts := store.NewExecutionReceipts(metrics, db, results, 100)
			store1 := store.NewMyExecutionReceipts(metrics, db, receipts)

			f(store1, db)
		})
	}

	t.Run("store1 one get one", func(t *testing.T) {
		withStore(t, func(store1 *store.MyExecutionReceipts, db storage.DB) {
			block := unittest.BlockFixture()
			receipt1 := unittest.ReceiptForBlockFixture(&block)

			err := db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
				return store1.BatchStoreMyReceipt(receipt1, rw)
			})
			require.NoError(t, err)

			actual, err := store1.MyReceipt(block.ID())
			require.NoError(t, err)

			require.Equal(t, receipt1, actual)
		})
	})

	t.Run("store1 same for the same block", func(t *testing.T) {
		withStore(t, func(store1 *store.MyExecutionReceipts, db storage.DB) {
			block := unittest.BlockFixture()

			receipt1 := unittest.ReceiptForBlockFixture(&block)

			err := db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
				return store1.BatchStoreMyReceipt(receipt1, rw)
			})
			require.NoError(t, err)

			err = db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
				return store1.BatchStoreMyReceipt(receipt1, rw)
			})
			require.NoError(t, err)
		})
	})

	t.Run("store1 different receipt for same block should fail", func(t *testing.T) {
		withStore(t, func(store1 *store.MyExecutionReceipts, db storage.DB) {
			block := unittest.BlockFixture()

			executor1 := unittest.IdentifierFixture()
			executor2 := unittest.IdentifierFixture()

			receipt1 := unittest.ReceiptForBlockExecutorFixture(&block, executor1)
			receipt2 := unittest.ReceiptForBlockExecutorFixture(&block, executor2)

			err := db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
				return store1.BatchStoreMyReceipt(receipt1, rw)
			})
			require.NoError(t, err)

			err = db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
				return store1.BatchStoreMyReceipt(receipt2, rw)
			})
			require.Error(t, err)
		})
	})

	t.Run("store1 different receipt concurrent for same block should fail", func(t *testing.T) {
		withStore(t, func(store1 *store.MyExecutionReceipts, db storage.DB) {
			block := unittest.BlockFixture()

			executor1 := unittest.IdentifierFixture()
			executor2 := unittest.IdentifierFixture()

			receipt1 := unittest.ReceiptForBlockExecutorFixture(&block, executor1)
			receipt2 := unittest.ReceiptForBlockExecutorFixture(&block, executor2)

			var wg sync.WaitGroup
			errCh := make(chan error, 2) // Buffered channel to capture errors

			wg.Add(2)
			go func() {
				defer wg.Done()
				err := db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
					return store1.BatchStoreMyReceipt(receipt1, rw)
				})
				errCh <- err
			}()

			go func() {
				defer wg.Done()
				err := db.WithReaderBatchWriter(func(rw storage.ReaderBatchWriter) error {
					return store1.BatchStoreMyReceipt(receipt2, rw)
				})
				errCh <- err
			}()

			wg.Wait()
			close(errCh)

			// Check that at least one of the operations failed
			errorsCount := 0
			for err := range errCh {
				if err != nil {
					errorsCount++
				}
			}

			require.Equal(t, 1, errorsCount, "One of the concurrent store1 operations should fail")
		})
	})
}

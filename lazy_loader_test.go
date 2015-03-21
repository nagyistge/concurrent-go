package concurrent

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	numTestLazyLoaderGoRoutines = 100
)

func TestLazyLoaderBasic(t *testing.T) {
	value := 0
	lazyLoader := NewLazyLoader(func() (interface{}, error) {
		value++
		return value, nil
	})
	runGoRoutinesAndWait(
		numTestLazyLoaderGoRoutines,
		func() {
			lazyLoaderValue, err := lazyLoader.Load()
			require.NoError(t, err)
			require.Equal(t, 1, lazyLoaderValue.(int))
		},
	)
}

func TestLazyLoaderError(t *testing.T) {
	value := 0
	lazyLoader := NewLazyLoader(func() (interface{}, error) {
		value++
		return nil, fmt.Errorf("error%d", value)
	})
	runGoRoutinesAndWait(
		numTestLazyLoaderGoRoutines,
		func() {
			lazyLoaderValue, err := lazyLoader.Load()
			require.Nil(t, lazyLoaderValue)
			require.Error(t, err)
			require.Equal(t, "error1", err.Error())
		},
	)
}

func runGoRoutinesAndWait(count int, f func()) {
	var wg sync.WaitGroup
	for i := 0; i < numTestLazyLoaderGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}

package concurrent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackgroundWorkerBasic(t *testing.T) {
	backgroundWorker := NewBackgroundWorker()
	backgroundWorker.Do(func() (interface{}, error) {
		return 1, nil
	})
	backgroundWorker.Do(func() (interface{}, error) {
		return 2, nil
	})
	backgroundWorker.Do(func() (interface{}, error) {
		return 3, nil
	})
	values, err := backgroundWorker.WaitToFinish()
	require.NoError(t, err)
	checkForValues(t, values, 1, 2, 3)
}

func TestBackgroundWorkerError(t *testing.T) {
	backgroundWorker := NewBackgroundWorker()
	backgroundWorker.Do(func() (interface{}, error) {
		return 1, nil
	})
	backgroundWorker.Do(func() (interface{}, error) {
		return 2, nil
	})
	backgroundWorker.Do(func() (interface{}, error) {
		return 3, errors.New("error")
	})
	values, err := backgroundWorker.WaitToFinish()
	require.Error(t, err)
	require.Equal(t, "[error]", err.Error())
	checkForValues(t, values, 1, 2, 3)

}

func checkForValues(t *testing.T, actualValues []interface{}, expectedValues ...interface{}) {
	require.Equal(t, len(expectedValues), len(actualValues))
	valuesMap := make(map[interface{}]bool)
	for _, expectedValue := range expectedValues {
		valuesMap[expectedValue] = false
	}
	for _, actualValue := range actualValues {
		valuesMap[actualValue] = true
	}
	for _, valuePresent := range valuesMap {
		require.True(t, valuePresent)
	}
}

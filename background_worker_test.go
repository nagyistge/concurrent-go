package concurrent

import (
	"errors"
	"testing"
)

func TestBackgroundWorkerBasic(t *testing.T) {
	backgroundWorker := NewBackgroundWorker()
	if err := backgroundWorker.Do(func() (interface{}, error) {
		return 1, nil
	}); err != nil {
		t.Error(err)
	}
	if err := backgroundWorker.Do(func() (interface{}, error) {
		return 2, nil
	}); err != nil {
		t.Error(err)
	}
	if err := backgroundWorker.Do(func() (interface{}, error) {
		return 3, nil
	}); err != nil {
		t.Error(err)
	}
	values, err := backgroundWorker.Close()
	if err != nil {
		t.Error(err)
	}
	checkForValues(t, values, 1, 2, 3)
}

func TestBackgroundWorkerError(t *testing.T) {
	backgroundWorker := NewBackgroundWorker()
	if err := backgroundWorker.Do(func() (interface{}, error) {
		return 1, nil
	}); err != nil {
		t.Error(err)
	}
	if err := backgroundWorker.Do(func() (interface{}, error) {
		return 2, nil
	}); err != nil {
		t.Error(err)
	}
	if err := backgroundWorker.Do(func() (interface{}, error) {
		return 3, errors.New("error")
	}); err != nil {
		t.Error(err)
	}
	values, err := backgroundWorker.Close()
	if err == nil || err.Error() != "[error]" {
		t.Errorf("expected [error], got %v", err)
	}
	checkForValues(t, values, 1, 2, 3)
}

func checkForValues(t *testing.T, actualValues []interface{}, expectedValues ...interface{}) {
	if len(expectedValues) != len(actualValues) {
		t.Errorf("expected set of %v, got %v", expectedValues, actualValues)
	}
	valuesMap := make(map[interface{}]bool)
	for _, expectedValue := range expectedValues {
		valuesMap[expectedValue] = false
	}
	for _, actualValue := range actualValues {
		valuesMap[actualValue] = true
	}
	for _, valuePresent := range valuesMap {
		if !valuePresent {
			t.Fatalf("expected set of %v, got %v", expectedValues, actualValues)
		}
	}
}

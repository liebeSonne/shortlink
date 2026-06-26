package resetter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	mu    sync.Mutex
	value int
	data  []int
	name  string
}

func (ts *testStruct) Reset() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.value = 0
	ts.data = ts.data[:0]
	ts.name = ""
}

func newTestStruct() *testStruct {
	return &testStruct{
		value: 42,
		data:  make([]int, 0, 10),
		name:  "test",
	}
}

type testValueStruct struct {
	value int
	name  string
}

func (vs *testValueStruct) Reset() {
	vs.value = 0
	vs.name = ""
}

func TestPool_NewPool(t *testing.T) {
	tests := []struct {
		name    string
		newFunc func() *testStruct
		err     error
	}{
		{
			name:    "create pool with valid constructor",
			newFunc: newTestStruct,
			err:     nil,
		},
		{
			name: "create pool with custom constructor",
			newFunc: func() *testStruct {
				return &testStruct{
					value: 100,
					data:  make([]int, 0, 20),
					name:  "custom",
				}
			},
			err: nil,
		},
		{
			name:    "create pool with nil constructor should return error",
			newFunc: nil,
			err:     ErrInvalidPoolConstructorFunc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := NewPool(tt.newFunc)

			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, pool, "Pool should not be nil")

			obj := pool.Get()
			assert.NotNil(t, obj, "Get should return non-nil object")
		})
	}
}

func TestPool_Get(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *Pool[*testStruct]
		expectReset bool
	}{
		{
			name: "get from empty pool creates new object",
			setupFunc: func() *Pool[*testStruct] {
				p, err := NewPool(newTestStruct)
				require.NoError(t, err)
				return p
			},
			expectReset: false,
		},
		{
			name: "get from pool with objects returns existing object",
			setupFunc: func() *Pool[*testStruct] {
				p, err := NewPool(newTestStruct)
				require.NoError(t, err)

				obj := newTestStruct()
				obj.value = 999
				obj.data = append(obj.data, 1, 2, 3)
				obj.name = "modified"
				p.Put(obj)
				return p
			},
			expectReset: true,
		},
		{
			name: "get multiple times from pool",
			setupFunc: func() *Pool[*testStruct] {
				p, err := NewPool(newTestStruct)
				require.NoError(t, err)
				return p
			},
			expectReset: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setupFunc()
			obj := p.Get()

			require.NotNil(t, obj, "Get() should not return nil")

			if tt.expectReset {
				assert.Equal(t, 0, obj.value, "Object should be reset")
				assert.Empty(t, obj.data, "Object data should be reset")
				assert.Empty(t, obj.name, "Object name should be reset")
			}
		})
	}
}

func TestPool_Put(t *testing.T) {
	tests := []struct {
		name      string
		modifyObj func(*testStruct)
	}{
		{
			name: "put resets object with modified fields",
			modifyObj: func(obj *testStruct) {
				obj.value = 100
				obj.data = append(obj.data, 1, 2, 3)
				obj.name = "modified"
			},
		},
		{
			name: "put resets object with large data",
			modifyObj: func(obj *testStruct) {
				obj.value = 999
				for i := 0; i < 1000; i++ {
					obj.data = append(obj.data, i)
				}
				obj.name = "large_data"
			},
		},
		{
			name: "put with empty object",
			modifyObj: func(obj *testStruct) {
				// do nothing
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPool(newTestStruct)
			require.NoError(t, err)

			obj := newTestStruct()

			tt.modifyObj(obj)

			oldValue := obj.value
			oldDataLen := len(obj.data)
			oldName := obj.name

			p.Put(obj)

			got := p.Get()
			require.NotNil(t, got, "Get after Put should return non-nil object")

			assert.Equal(t, 0, got.value, "Put should reset value from %d to 0", oldValue)
			assert.Empty(t, got.data, "Put should reset data from length %d to empty", oldDataLen)
			assert.Empty(t, got.name, "Put should reset name from '%s' to empty", oldName)

			// Проверяем, что емкость слайса сохранилась
			assert.Equal(t, cap(got.data), cap(obj.data), "Slice capacity should be preserved")
		})
	}
}

func TestPool_Concurrent(t *testing.T) {
	tests := []struct {
		name      string
		goroutine int
		iter      int
	}{
		{
			name:      "concurrent get and put with 10 goroutines",
			goroutine: 10,
			iter:      100,
		},
		{
			name:      "concurrent get and put with 50 goroutines",
			goroutine: 50,
			iter:      1000,
		},
		{
			name:      "concurrent get and put with 100 goroutines",
			goroutine: 100,
			iter:      100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPool(newTestStruct)
			require.NoError(t, err)

			var wg sync.WaitGroup

			for i := 0; i < tt.goroutine; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for j := 0; j < tt.iter; j++ {
						obj := p.Get()
						require.NotNil(t, obj, "Get should return non-nil object")

						obj.value = id*100 + j
						obj.data = append(obj.data, j)
						obj.name = "goroutine"

						p.Put(obj)
					}
				}(i)
			}

			wg.Wait()

			obj := p.Get()
			require.NotNil(t, obj, "Pool should still work after concurrent operations")
			assert.Equal(t, 0, obj.value, "Object should be properly reset")
			assert.Empty(t, obj.data, "Object data should be empty")

			p.Put(obj)
		})
	}
}

func TestPool_GetPutIntegration(t *testing.T) {
	tests := []struct {
		name       string
		operations int
	}{
		{
			name:       "integration test with 100 operations",
			operations: 100,
		},
		{
			name:       "integration test with 1000 operations",
			operations: 1000,
		},
		{
			name:       "integration test with 10000 operations",
			operations: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPool(newTestStruct)
			require.NoError(t, err)

			objects := make([]*testStruct, 0, tt.operations)

			for i := 0; i < tt.operations; i++ {
				obj := p.Get()
				require.NotNil(t, obj, "Get should return non-nil object")

				obj.value = i
				obj.data = append(obj.data, i)
				obj.name = "test"
				objects = append(objects, obj)
			}

			assert.Len(t, objects, tt.operations, "Should have correct number of objects")

			for _, obj := range objects {
				p.Put(obj)
			}

			for i := 0; i < 10; i++ {
				obj := p.Get()
				require.NotNil(t, obj, "Get should return non-nil object")

				assert.Equal(t, 0, obj.value, "Object should be reset")
				assert.Empty(t, obj.data, "Object data should be empty")
				assert.Empty(t, obj.name, "Object name should be empty")

				p.Put(obj)
			}
		})
	}
}

func TestPool_WithDifferentTypes(t *testing.T) {
	t.Run("works with pointer type", func(t *testing.T) {
		p, err := NewPool(newTestStruct)
		require.NoError(t, err)

		obj := p.Get()
		assert.NotNil(t, obj)
		p.Put(obj)
	})

	t.Run("works with value type", func(t *testing.T) {
		func() {
			vs := &testValueStruct{value: 42, name: "test"}
			vs.Reset()
		}()

		newValue := func() *testValueStruct {
			return &testValueStruct{value: 42, name: "test"}
		}

		p, err := NewPool(newValue)
		require.NoError(t, err)

		obj := p.Get()
		assert.NotNil(t, obj)
		p.Put(obj)
	})
}

func TestPool_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "multiple get without put",
			testFunc: func(t *testing.T) {
				p, err := NewPool(newTestStruct)
				require.NoError(t, err)

				obj1 := p.Get()
				obj2 := p.Get()

				assert.NotNil(t, obj1)
				assert.NotNil(t, obj2)

				p.Put(obj1)
				p.Put(obj2)
			},
		},
		{
			name: "put and get many times",
			testFunc: func(t *testing.T) {
				p, err := NewPool(newTestStruct)
				require.NoError(t, err)

				for i := 0; i < 100; i++ {
					obj := p.Get()
					obj.value = i
					p.Put(obj)
				}

				obj := p.Get()
				assert.Equal(t, 0, obj.value, "Object should be reset")
				p.Put(obj)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// Benchmarks

func BenchmarkPool_GetPut(b *testing.B) {
	p, err := NewPool(newTestStruct)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := p.Get()
		obj.value = i
		obj.data = append(obj.data, i)
		p.Put(obj)
	}
}

func BenchmarkPool_GetPutParallel(b *testing.B) {
	p, err := NewPool(newTestStruct)
	require.NoError(b, err)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := p.Get()
			obj.value = 42
			obj.data = append(obj.data, 42)
			p.Put(obj)
		}
	})
}

func BenchmarkPool_GetOnly(b *testing.B) {
	p, err := NewPool(newTestStruct)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := p.Get()
		_ = obj
	}
}

func BenchmarkPool_PutOnly(b *testing.B) {
	p, err := NewPool(newTestStruct)
	require.NoError(b, err)

	objects := make([]*testStruct, b.N)

	for i := 0; i < b.N; i++ {
		objects[i] = p.Get()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Put(objects[i])
	}
}

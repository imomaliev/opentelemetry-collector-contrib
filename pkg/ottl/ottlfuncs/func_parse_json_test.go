// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ottlfuncs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

func Test_ParseJSON(t *testing.T) {
	tests := []struct {
		name   string
		target ottl.Getter[any]
		want   func(pcommon.Map)
	}{
		{
			name: "handle string",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":"string value"}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				expectedMap.PutStr("test", "string value")
			},
		},
		{
			name: "handle bool",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":true}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				expectedMap.PutBool("test", true)
			},
		},
		{
			name: "handle int",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":1}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				expectedMap.PutDouble("test", 1)
			},
		},
		{
			name: "handle float",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":1.1}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				expectedMap.PutDouble("test", 1.1)
			},
		},
		{
			name: "handle nil",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":null}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				expectedMap.PutEmpty("test")
			},
		},
		{
			name: "handle array",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":["string","value"]}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				emptySlice := expectedMap.PutEmptySlice("test")
				emptySlice.AppendEmpty().SetStr("string")
				emptySlice.AppendEmpty().SetStr("value")
			},
		},
		{
			name: "handle nested object",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test":{"nested":"true"}}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				newMap := expectedMap.PutEmptyMap("test")
				newMap.PutStr("nested", "true")
			},
		},
		{
			name: "updates existing",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"existing":"pass"}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				expectedMap.PutStr("existing", "pass")
			},
		},
		{
			name: "complex",
			target: ottl.StandardGetSetter[any]{
				Getter: func(ctx context.Context, tCtx any) (interface{}, error) {
					return `{"test1":{"nested":"true"},"test2":"string","test3":1,"test4":1.1,"test5":[[1], [2, 3],[]],"test6":null}`, nil
				},
			},
			want: func(expectedMap pcommon.Map) {
				newMap := expectedMap.PutEmptyMap("test1")
				newMap.PutStr("nested", "true")
				expectedMap.PutStr("test2", "string")
				expectedMap.PutDouble("test3", 1)
				expectedMap.PutDouble("test4", 1.1)
				slice := expectedMap.PutEmptySlice("test5")
				slice0 := slice.AppendEmpty().SetEmptySlice()
				slice0.AppendEmpty().SetDouble(1)
				slice1 := slice.AppendEmpty().SetEmptySlice()
				slice1.AppendEmpty().SetDouble(2)
				slice1.AppendEmpty().SetDouble(3)
				slice.AppendEmpty().SetEmptySlice()
				expectedMap.PutEmpty("test6")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprFunc, err := ParseJSON(tt.target)
			assert.NoError(t, err)

			result, err := exprFunc(context.Background(), nil)
			assert.NoError(t, err)

			resultMap, ok := result.(pcommon.Map)
			if !ok {
				assert.Fail(t, "pcommon.Map not returned")
			}

			expected := pcommon.NewMap()
			tt.want(expected)

			assert.Equal(t, expected.Len(), resultMap.Len())
			expected.Range(func(k string, v pcommon.Value) bool {
				ev, _ := expected.Get(k)
				av, _ := resultMap.Get(k)
				assert.Equal(t, ev, av)
				return true
			})
		})
	}
}

func Test_ParseJSON_Error(t *testing.T) {
	target := &ottl.StandardGetSetter[interface{}]{
		Getter: func(ctx context.Context, tCtx interface{}) (interface{}, error) {
			return 1, nil
		},
	}
	exprFunc, err := ParseJSON[interface{}](target)
	assert.NoError(t, err)
	_, err = exprFunc(context.Background(), nil)
	assert.Error(t, err)
}

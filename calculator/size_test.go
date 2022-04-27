/*
 * Copyright 2015-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package calculator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	memory "github.com/wreulicke/emc/calculator"
)

func MustParseSize(v string) memory.Size {
	r, err := memory.ParseSize(v)
	if err != nil {
		panic(err)
	}
	return r
}

func TestMemorySize(t *testing.T) {
	assert.Equal(t, memory.Size(1023).String(), "0")
	assert.Equal(t, memory.Size(memory.Kibi+1023).String(), "1K")
	assert.Equal(t, memory.Size(memory.Mibi+1023).String(), "1M")
	assert.Equal(t, memory.Size(memory.Gibi+1023).String(), "1G")
	assert.Equal(t, memory.Size((memory.Tibi*1024)+1023).String(), "1024T")

	assert.Equal(t, MustParseSize("0"), memory.Size(0))
	assert.Equal(t, MustParseSize("1"), memory.Size(1))
	assert.Equal(t, MustParseSize("1b"), memory.Size(1))
	assert.Equal(t, MustParseSize("1k"), memory.Size(memory.Kibi))
	assert.Equal(t, MustParseSize("1K"), memory.Size(memory.Kibi))
	assert.Equal(t, MustParseSize("1m"), memory.Size(memory.Mibi))
	assert.Equal(t, MustParseSize("1M"), memory.Size(memory.Mibi))
	assert.Equal(t, MustParseSize("1g"), memory.Size(memory.Gibi))
	assert.Equal(t, MustParseSize("1G"), memory.Size(memory.Gibi))
	assert.Equal(t, MustParseSize("1t"), memory.Size(memory.Tibi))
	assert.Equal(t, MustParseSize("1T"), memory.Size(memory.Tibi))

	assert.Equal(t, MustParseSize("\t\r\n 1"), memory.Size(1))
	assert.Equal(t, MustParseSize("1 \t\r\n"), memory.Size(1))

	var err error

	_, err = memory.ParseSize("")
	assert.NotNil(t, err)

	_, err = memory.ParseSize("-1")
	assert.NotNil(t, err)

	_, err = memory.ParseSize("1A")
	assert.NotNil(t, err)

	_, err = memory.ParseSize("0x1")
	assert.NotNil(t, err)

	_, err = memory.ParseSize("1.0")
	assert.NotNil(t, err)

	_, err = memory.ParseSize("1 0")
	assert.NotNil(t, err)
}

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

func TestReservedCodeCache(t *testing.T) {
	assert.Equal(t, memory.ReservedCodeCache(memory.Kibi).String(), "-XX:ReservedCodeCacheSize=1K")

	assert.True(t, memory.IsReservedCodeCache("-XX:ReservedCodeCacheSize=1K"))

	assert.False(t, memory.IsReservedCodeCache("-Xss1K"))

	r, _ := memory.ParseReservedCodeCache("-XX:ReservedCodeCacheSize=1K")
	assert.Equal(t, r, memory.ReservedCodeCache(memory.Kibi))
}

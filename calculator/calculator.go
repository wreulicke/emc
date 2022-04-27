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

package calculator

import (
	"fmt"
)

type Calculator struct {
	HeadRoom         int
	JvmOptions       *JVMOptions
	LoadedClassCount int64
	ThreadCount      int64
	TotalMemory      Size
}

func (c Calculator) Calculate() ([]fmt.Stringer, error) {
	var options []fmt.Stringer

	j := c.JvmOptions
	if j == nil {
		j = &JVMOptions{}
	}

	headRoom := c.headRoom()

	directMemory := j.MaxDirectMemory
	if directMemory == nil {
		d := DefaultMaxDirectMemory
		directMemory = &d
		options = append(options, *directMemory)
	}

	metaspace := j.MaxMetaspace
	if metaspace == nil {
		m := c.metaspace()
		metaspace = &m
		options = append(options, *metaspace)
	}

	reservedCodeCache := j.ReservedCodeCache
	if reservedCodeCache == nil {
		r := DefaultReservedCodeCache
		reservedCodeCache = &r
		options = append(options, *reservedCodeCache)
	}

	stack := j.Stack
	if stack == nil {
		s := DefaultStack
		stack = &s
		options = append(options, *stack)
	}

	overhead := c.overhead(headRoom, directMemory, metaspace, reservedCodeCache, stack)
	available := c.TotalMemory

	if overhead > available {
		return nil, fmt.Errorf("required memory %s is greater than %s available for allocation: %s, %s, %s, %s x %d threads",
			overhead, available, directMemory, metaspace, reservedCodeCache, stack, c.ThreadCount)
	}

	heap := j.MaxHeap
	if heap == nil {
		h := c.heap(overhead)
		heap = &h
		options = append(options, *heap)
	}

	if overhead+Size(*heap) > available {
		return nil, fmt.Errorf("required memory %s is greater than %s available for allocation: %s, %s, %s, %s, %s x %d threads",
			overhead+Size(*heap), available, directMemory, heap, metaspace, reservedCodeCache, stack, c.ThreadCount)
	}

	return options, nil
}

func (c Calculator) headRoom() Size {
	return Size(float64(c.TotalMemory) * (float64(c.HeadRoom) / 100))
}

func (c Calculator) heap(overhead Size) MaxHeap {
	return MaxHeap(Size(c.TotalMemory) - overhead)
}

func (c Calculator) metaspace() MaxMetaspace {
	return MaxMetaspace((c.LoadedClassCount * 5800) + 14000000)
}

func (c Calculator) overhead(headRoom Size, directMemory *MaxDirectMemory, metaspace *MaxMetaspace, reservedCodeCache *ReservedCodeCache, stack *Stack) Size {
	return headRoom +
		Size(*directMemory) +
		Size(*metaspace) +
		Size(*reservedCodeCache) +
		Size(int64(*stack)*int64(c.ThreadCount))
}

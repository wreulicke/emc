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
	"strings"
)

var DefaultJVMOptions = JVMOptions{}

const FlagJVMOptions = "jvm-options"

type JVMOptions struct {
	MaxDirectMemory   *MaxDirectMemory
	MaxHeap           *MaxHeap
	MaxMetaspace      *MaxMetaspace
	ReservedCodeCache *ReservedCodeCache
	Stack             *Stack
}

func (j *JVMOptions) Set(s string) error {
	for _, c := range strings.Split(s, " ") {
		if IsMaxDirectMemory(c) {
			m, err := ParseMaxDirectMemory(c)
			if err != nil {
				return err
			}

			j.MaxDirectMemory = &m
		} else if IsMaxHeap(c) {
			m, err := ParseMaxHeap(c)
			if err != nil {
				return err
			}

			j.MaxHeap = &m
		} else if IsMaxMetaspace(c) {
			m, err := ParseMaxMetaspace(c)
			if err != nil {
				return err
			}

			j.MaxMetaspace = &m
		} else if IsReservedCodeCache(c) {
			r, err := ParseReservedCodeCache(c)
			if err != nil {
				return err
			}

			j.ReservedCodeCache = &r
		} else if IsStack(c) {
			s, err := ParseStack(c)
			if err != nil {
				return err
			}

			j.Stack = &s
		}
	}

	return nil
}

func (j *JVMOptions) String() string {
	var values []string

	if j.MaxDirectMemory != nil {
		values = append(values, j.MaxDirectMemory.String())
	}

	if j.MaxHeap != nil {
		values = append(values, j.MaxHeap.String())
	}

	if j.MaxMetaspace != nil {
		values = append(values, j.MaxMetaspace.String())
	}

	if j.ReservedCodeCache != nil {
		values = append(values, j.ReservedCodeCache.String())
	}

	if j.Stack != nil {
		values = append(values, j.Stack.String())
	}

	return strings.Join(values, " ")
}

func (j *JVMOptions) Type() string {
	return "string"
}

func (j *JVMOptions) Validate() error {
	return nil
}

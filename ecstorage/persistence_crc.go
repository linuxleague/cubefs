// Copyright 2018 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package ecstorage

import (
	"encoding/binary"
)

type BlockCrc struct {
	BlockNo int
	Crc     uint32
}
type BlockCrcArr []*BlockCrc

const (
	BaseExtentIDOffset = 0
)

func (arr BlockCrcArr) Len() int           { return len(arr) }
func (arr BlockCrcArr) Less(i, j int) bool { return arr[i].BlockNo < arr[j].BlockNo }
func (arr BlockCrcArr) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }

type EcUpdateCrcFunc func(partitionId, extentId, blockNo uint64, size, crc uint32) (err error)
type GetExtentCrcFunc func(extentID uint64) (crc uint32, err error)

func (s *ExtentStore) PersistenceBaseExtentID(extentID uint64) (err error) {
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, extentID)
	_, err = s.metadataFp.WriteAt(value, BaseExtentIDOffset)
	return
}

func (s *ExtentStore) GetPreAllocSpaceExtentIDOnVerfiyFile() (extentID uint64) {
	value := make([]byte, 8)
	_, err := s.metadataFp.WriteAt(value, 8)
	if err != nil {
		return
	}
	extentID = binary.BigEndian.Uint64(value)
	return
}

func (s *ExtentStore) GetPersistenceBaseExtentID() (extentID uint64, err error) {
	data := make([]byte, 8)
	_, err = s.metadataFp.ReadAt(data, 0)
	if err != nil {
		return
	}
	extentID = binary.BigEndian.Uint64(data)
	return
}

func (s *ExtentStore) PersistenceHasDeleteExtent(extentID uint64) (err error) {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, extentID)
	if _, err = s.normalExtentDeleteFp.Write(data); err != nil {
		return
	}
	return
}
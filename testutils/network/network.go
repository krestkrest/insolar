/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package network

import (
	"github.com/insolar/insolar/core"
)

type testNetwork struct {
}

func (n *testNetwork) GetNodeID() core.RecordRef {
	ref, _ := core.NewRefFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa")
	return *ref
}

func (n *testNetwork) GetGlobuleID() core.GlobuleID {
	return 0
}

func (n *testNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Parcel) ([]byte, error) {
	return make([]byte, 0), nil
}
func (n *testNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error {
	return nil
}
func (n *testNetwork) GetAddress() string {
	return ""
}
func (n *testNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {

}

func GetTestNetwork() core.Network {
	return &testNetwork{}
}

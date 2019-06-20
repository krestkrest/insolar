//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handle

import (
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
)

type GetPendingFilament struct {
	dep *proc.Dependencies

	meta payload.Meta
}

func NewGetPendingFilament(dep *proc.Dependencies, meta payload.Meta) *GetPendingFilament {
	return &GetPendingFilament{
		dep:  dep,
		meta: meta,
	}
}

func (s *GetPendingFilament) Present(ctx context.Context, f flow.Flow) error {
	gpf := payload.GetPendingFilament{}
	err := gpf.Unmarshal(s.meta.Payload)
	if err != nil {
		panic(err)
		return errors.Wrap(err, "failed to unmarshal payload")
	}

	getFilament := proc.NewGetPendingFilament(s.meta, gpf.ObjectID, gpf.StartFrom, gpf.ReadUntil)
	s.dep.GetPendingFilament(getFilament)
	return f.Procedure(ctx, getFilament, false)
}

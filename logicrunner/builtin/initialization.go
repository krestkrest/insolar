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

// Code generated by insgocc. DO NOT EDIT.
// source template in logicrunner/preprocessor/templates

package builtin

import (
	"github.com/pkg/errors"

	account "github.com/insolar/insolar/logicrunner/builtin/contract/account"
	costcenter "github.com/insolar/insolar/logicrunner/builtin/contract/costcenter"
	deposit "github.com/insolar/insolar/logicrunner/builtin/contract/deposit"
	helloworld "github.com/insolar/insolar/logicrunner/builtin/contract/helloworld"
	member "github.com/insolar/insolar/logicrunner/builtin/contract/member"
	migrationadmin "github.com/insolar/insolar/logicrunner/builtin/contract/migrationadmin"
	migrationdaemon "github.com/insolar/insolar/logicrunner/builtin/contract/migrationdaemon"
	migrationshard "github.com/insolar/insolar/logicrunner/builtin/contract/migrationshard"
	nodedomain "github.com/insolar/insolar/logicrunner/builtin/contract/nodedomain"
	noderecord "github.com/insolar/insolar/logicrunner/builtin/contract/noderecord"
	pkshard "github.com/insolar/insolar/logicrunner/builtin/contract/pkshard"
	rootdomain "github.com/insolar/insolar/logicrunner/builtin/contract/rootdomain"
	wallet "github.com/insolar/insolar/logicrunner/builtin/contract/wallet"

	XXX_insolar "github.com/insolar/insolar/insolar"
	XXX_artifacts "github.com/insolar/insolar/logicrunner/artifacts"
)

func InitializeContractMethods() map[string]XXX_insolar.ContractWrapper {
	return map[string]XXX_insolar.ContractWrapper{
		"account":         account.Initialize(),
		"costcenter":      costcenter.Initialize(),
		"deposit":         deposit.Initialize(),
		"helloworld":      helloworld.Initialize(),
		"member":          member.Initialize(),
		"migrationadmin":  migrationadmin.Initialize(),
		"migrationdaemon": migrationdaemon.Initialize(),
		"migrationshard":  migrationshard.Initialize(),
		"nodedomain":      nodedomain.Initialize(),
		"noderecord":      noderecord.Initialize(),
		"pkshard":         pkshard.Initialize(),
		"rootdomain":      rootdomain.Initialize(),
		"wallet":          wallet.Initialize(),
	}
}

func shouldLoadRef(strRef string) XXX_insolar.Reference {
	ref, err := XXX_insolar.NewReferenceFromBase58(strRef)
	if err != nil {
		panic(errors.Wrap(err, "Unexpected error, bailing out"))
	}
	return *ref
}

func InitializeCodeRefs() map[XXX_insolar.Reference]string {
	rv := make(map[XXX_insolar.Reference]string, 13)

	rv[shouldLoadRef("0111A7rimrANEAnwBT1kvAhHeHp9NPTFJMLKVng8GLH5.record")] = "account"
	rv[shouldLoadRef("0111A7tUo1FeZ5DSoroiinMCKwzLacaYBAAcwAaNj6bc.record")] = "costcenter"
	rv[shouldLoadRef("0111A79KGpeDUjYhRJP1n1AwYgwU9KEWmc2TNNc3KQjV.record")] = "deposit"
	rv[shouldLoadRef("0111A5w1GcnTsht82duVrnWdVHVNyrxCUVcSPLtgQCPR.record")] = "helloworld"
	rv[shouldLoadRef("0111A72gPKWyrF9c7yzDoccRoPQ62g1uQQDBecWJwAYr.record")] = "member"
	rv[shouldLoadRef("0111A6516TVnMLh8DAzTWbtEJrgZkESeCpdn2viV6D61.record")] = "migrationadmin"
	rv[shouldLoadRef("0111A7PzUnidJKg3DDo82FyyYFukEyKJYmLKoCFfQmoK.record")] = "migrationdaemon"
	rv[shouldLoadRef("0111A66L3aoDPf2wedyRo2gyns8ghV9vdeJdJntVaGEf.record")] = "migrationshard"
	rv[shouldLoadRef("0111A7Q5FK2ebPG9WnSiUc4iqF45w9oYkJkRjEtBohGe.record")] = "nodedomain"
	rv[shouldLoadRef("0111A86xPKUQ1ZxSscgv5brbw93LkwiVhUWgGrYYsMar.record")] = "noderecord"
	rv[shouldLoadRef("0111A5tzn16hnKGCZCyYA8Dv9FALvPYYQu4VA41SVx6s.record")] = "pkshard"
	rv[shouldLoadRef("0111A63R5cAgGHC5DJffqF16vUkCuSVj3GExbMLy56cS.record")] = "rootdomain"
	rv[shouldLoadRef("0111A5e49cJW6GKGegWBhtgrJs7nFh1kSWhBtT2VgK4t.record")] = "wallet"

	return rv
}

func InitializePrototypeRefs() map[XXX_insolar.Reference]string {
	rv := make(map[XXX_insolar.Reference]string, 13)

	rv[shouldLoadRef("0111A62X73fkPeY5vK6NjcXgmL9d37DgRRNtHNLGaEse")] = "account"
	rv[shouldLoadRef("0111A62HrJvAimG7M1r8XdeBVMw4X6ge8hGzVStfnn4e")] = "costcenter"
	rv[shouldLoadRef("0111A7ctasuNUug8BoK4VJNuAFJ73rnH8bH5zqd5HrDj")] = "deposit"
	rv[shouldLoadRef("0111A85JAZugtAkQErbDe3eAaTw56DPLku8QGymJUCt2")] = "helloworld"
	rv[shouldLoadRef("0111A7UqbgvFXj9vkCAaNYSAkWLapu62eU5AUSv3y4JY")] = "member"
	rv[shouldLoadRef("0111A8DhUhw5pzyvzVg1qXomNEHXs7kDtJRQGSD1PUpc")] = "migrationadmin"
	rv[shouldLoadRef("0111A7jZX41e1SpH9oW3F2dgUvVQdjSqXEAGQSxhbqmD")] = "migrationdaemon"
	rv[shouldLoadRef("0111A7FNYLZLYXYWZPbkMhCAPwV9nYrWWE7L57CtdJCj")] = "migrationshard"
	rv[shouldLoadRef("0111A6NKbCjpzFr9MttfcWV8vX8eFjiyGPPfSH1AMtwN")] = "nodedomain"
	rv[shouldLoadRef("0111A5fZeApbGhcsLrbfGy82kKLgapF93GhNPMLSYaPY")] = "noderecord"
	rv[shouldLoadRef("0111A5x8N1VJTm7BKYgzSe6TWHcFi98QZgw3AnkYiKML")] = "pkshard"
	rv[shouldLoadRef("0111A84uiiTD1LXAHNP4GMA6YJFjbnCdkRia2pCqwBV5")] = "rootdomain"
	rv[shouldLoadRef("0111A5gmRD1ZbHjQh7DgH9SrCK4a1qfwEUP5xAir6i8L")] = "wallet"

	return rv
}

func InitializeCodeDescriptors() []XXX_artifacts.CodeDescriptor {
	rv := make([]XXX_artifacts.CodeDescriptor, 0, 13)

	// account
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A7rimrANEAnwBT1kvAhHeHp9NPTFJMLKVng8GLH5.record"),
	))
	// costcenter
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A7tUo1FeZ5DSoroiinMCKwzLacaYBAAcwAaNj6bc.record"),
	))
	// deposit
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A79KGpeDUjYhRJP1n1AwYgwU9KEWmc2TNNc3KQjV.record"),
	))
	// helloworld
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A5w1GcnTsht82duVrnWdVHVNyrxCUVcSPLtgQCPR.record"),
	))
	// member
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A72gPKWyrF9c7yzDoccRoPQ62g1uQQDBecWJwAYr.record"),
	))
	// migrationadmin
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A6516TVnMLh8DAzTWbtEJrgZkESeCpdn2viV6D61.record"),
	))
	// migrationdaemon
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A7PzUnidJKg3DDo82FyyYFukEyKJYmLKoCFfQmoK.record"),
	))
	// migrationshard
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A66L3aoDPf2wedyRo2gyns8ghV9vdeJdJntVaGEf.record"),
	))
	// nodedomain
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A7Q5FK2ebPG9WnSiUc4iqF45w9oYkJkRjEtBohGe.record"),
	))
	// noderecord
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A86xPKUQ1ZxSscgv5brbw93LkwiVhUWgGrYYsMar.record"),
	))
	// pkshard
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A5tzn16hnKGCZCyYA8Dv9FALvPYYQu4VA41SVx6s.record"),
	))
	// rootdomain
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A63R5cAgGHC5DJffqF16vUkCuSVj3GExbMLy56cS.record"),
	))
	// wallet
	rv = append(rv, XXX_artifacts.NewCodeDescriptor(
		/* code:        */ nil,
		/* machineType: */ XXX_insolar.MachineTypeBuiltin,
		/* ref:         */ shouldLoadRef("0111A5e49cJW6GKGegWBhtgrJs7nFh1kSWhBtT2VgK4t.record"),
	))

	return rv
}

func InitializePrototypeDescriptors() []XXX_artifacts.PrototypeDescriptor {
	rv := make([]XXX_artifacts.PrototypeDescriptor, 0, 13)

	{ // account
		pRef := shouldLoadRef("0111A62X73fkPeY5vK6NjcXgmL9d37DgRRNtHNLGaEse")
		cRef := shouldLoadRef("0111A7rimrANEAnwBT1kvAhHeHp9NPTFJMLKVng8GLH5.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // costcenter
		pRef := shouldLoadRef("0111A62HrJvAimG7M1r8XdeBVMw4X6ge8hGzVStfnn4e")
		cRef := shouldLoadRef("0111A7tUo1FeZ5DSoroiinMCKwzLacaYBAAcwAaNj6bc.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // deposit
		pRef := shouldLoadRef("0111A7ctasuNUug8BoK4VJNuAFJ73rnH8bH5zqd5HrDj")
		cRef := shouldLoadRef("0111A79KGpeDUjYhRJP1n1AwYgwU9KEWmc2TNNc3KQjV.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // helloworld
		pRef := shouldLoadRef("0111A85JAZugtAkQErbDe3eAaTw56DPLku8QGymJUCt2")
		cRef := shouldLoadRef("0111A5w1GcnTsht82duVrnWdVHVNyrxCUVcSPLtgQCPR.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // member
		pRef := shouldLoadRef("0111A7UqbgvFXj9vkCAaNYSAkWLapu62eU5AUSv3y4JY")
		cRef := shouldLoadRef("0111A72gPKWyrF9c7yzDoccRoPQ62g1uQQDBecWJwAYr.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // migrationadmin
		pRef := shouldLoadRef("0111A8DhUhw5pzyvzVg1qXomNEHXs7kDtJRQGSD1PUpc")
		cRef := shouldLoadRef("0111A6516TVnMLh8DAzTWbtEJrgZkESeCpdn2viV6D61.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // migrationdaemon
		pRef := shouldLoadRef("0111A7jZX41e1SpH9oW3F2dgUvVQdjSqXEAGQSxhbqmD")
		cRef := shouldLoadRef("0111A7PzUnidJKg3DDo82FyyYFukEyKJYmLKoCFfQmoK.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // migrationshard
		pRef := shouldLoadRef("0111A7FNYLZLYXYWZPbkMhCAPwV9nYrWWE7L57CtdJCj")
		cRef := shouldLoadRef("0111A66L3aoDPf2wedyRo2gyns8ghV9vdeJdJntVaGEf.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // nodedomain
		pRef := shouldLoadRef("0111A6NKbCjpzFr9MttfcWV8vX8eFjiyGPPfSH1AMtwN")
		cRef := shouldLoadRef("0111A7Q5FK2ebPG9WnSiUc4iqF45w9oYkJkRjEtBohGe.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // noderecord
		pRef := shouldLoadRef("0111A5fZeApbGhcsLrbfGy82kKLgapF93GhNPMLSYaPY")
		cRef := shouldLoadRef("0111A86xPKUQ1ZxSscgv5brbw93LkwiVhUWgGrYYsMar.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // pkshard
		pRef := shouldLoadRef("0111A5x8N1VJTm7BKYgzSe6TWHcFi98QZgw3AnkYiKML")
		cRef := shouldLoadRef("0111A5tzn16hnKGCZCyYA8Dv9FALvPYYQu4VA41SVx6s.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // rootdomain
		pRef := shouldLoadRef("0111A84uiiTD1LXAHNP4GMA6YJFjbnCdkRia2pCqwBV5")
		cRef := shouldLoadRef("0111A63R5cAgGHC5DJffqF16vUkCuSVj3GExbMLy56cS.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	{ // wallet
		pRef := shouldLoadRef("0111A5gmRD1ZbHjQh7DgH9SrCK4a1qfwEUP5xAir6i8L")
		cRef := shouldLoadRef("0111A5e49cJW6GKGegWBhtgrJs7nFh1kSWhBtT2VgK4t.record")
		rv = append(rv, XXX_artifacts.NewPrototypeDescriptor(
			/* head:         */ pRef,
			/* state:        */ *pRef.GetLocal(),
			/* code:         */ cRef,
		))
	}

	return rv
}

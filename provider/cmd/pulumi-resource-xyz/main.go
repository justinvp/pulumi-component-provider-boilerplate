// Copyright 2016-2022, Pulumi Corporation.
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

package main

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

func main() {
	p.RunProvider("xyz", Version,
		// We tell the provider what resources it needs to support.
		// In this case, a single custom resource.
		infer.Provider(infer.Options{
			Components: []infer.InferredComponent{
				infer.Component[*RandomLogin, RandomLoginArgs, *RandomLoginState](),
			},
		}))
}

type RandomLogin struct{}

type RandomLoginArgs struct {
	PasswordLength pulumi.IntPtrInput `pulumi:"passwordLength"`
	PetName        bool               `pulumi:"petName"`
}

type RandomLoginState struct {
	pulumi.ResourceState
	Username pulumi.StringOutput `pulumi:"username"`
	Password pulumi.StringOutput `pulumi:"password"`
}

func (r *RandomLogin) Construct(ctx *pulumi.Context, name, typ string, args RandomLoginArgs, opts pulumi.ResourceOption) (*RandomLoginState, error) {
	comp := &RandomLoginState{}
	err := ctx.RegisterComponentResource(typ, name, comp, opts)
	if err != nil {
		return nil, err
	}
	if args.PetName {
		pet, err := random.NewRandomPet(ctx, name+"-pet", &random.RandomPetArgs{}, pulumi.Parent(comp))
		if err != nil {
			return nil, err
		}
		comp.Username = pet.ID().ToStringOutput()
	} else {
		id, err := random.NewRandomId(ctx, name+"-id", &random.RandomIdArgs{
			ByteLength: pulumi.Int(8),
		}, pulumi.Parent(comp))
		if err != nil {
			return nil, err
		}
		comp.Username = id.ID().ToStringOutput()
	}
	var length pulumi.IntInput = pulumi.Int(16)
	if args.PasswordLength != nil {
		length = args.PasswordLength.ToIntPtrOutput().Elem()
	}
	password, err := random.NewRandomPassword(ctx, name+"-password", &random.RandomPasswordArgs{
		Length: length,
	}, pulumi.Parent(comp))
	if err != nil {
		return nil, err
	}
	comp.Password = password.Result

	return comp, nil
}

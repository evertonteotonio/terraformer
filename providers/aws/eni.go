// Copyright 2018 The Terraformer Authors.
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

package aws

import (
	"context"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var EniAllowEmptyValues = []string{"tags."}

type EniGenerator struct {
	AWSService
}

func (EniGenerator) createResources(enis *ec2.DescribeNetworkInterfacesOutput) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, eni := range enis.NetworkInterfaces {
		resources = append(resources, terraformutils.NewSimpleResource(
			aws.StringValue(eni.NetworkInterfaceId),
			aws.StringValue(eni.NetworkInterfaceId),
			"aws_network_interface",
			"aws",
			EniAllowEmptyValues,
		))
	}
	return resources
}

// Generate TerraformResources from AWS API,
// from each ENI creates 1 TerraformResource.
func (g *EniGenerator) InitResources() error {
	config, e := g.generateConfig()
	if e != nil {
		return e
	}
	svc := ec2.New(config)
	p := ec2.NewDescribeNetworkInterfacesPaginator(
		svc.DescribeNetworkInterfacesRequest(&ec2.DescribeNetworkInterfacesInput{}))
	for p.Next(context.Background()) {
		g.Resources = append(g.Resources, g.createResources(p.CurrentPage())...)
	}
	return p.Err()
}

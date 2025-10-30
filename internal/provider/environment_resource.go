package provider

import (
	"github.com/PipeOpsHQ/terraform-provider-pipeops/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewEnvironmentResource() resource.Resource {
	return resources.NewEnvironmentResource()
}

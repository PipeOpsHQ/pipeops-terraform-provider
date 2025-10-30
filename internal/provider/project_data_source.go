package provider

import (
	"github.com/PipeOpsHQ/terraform-provider-pipeops/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewProjectDataSource() datasource.DataSource {
	return datasources.NewProjectDataSource()
}

// Copyright (c) HashiCorp, Inc.
// Copyright (C) Aaron Edwards
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &TimeUtilsProvider{}

type TimeUtilsProvider struct {
	version string
}

func (p *TimeUtilsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "timeutils"
	resp.Version = p.version
}

func (p *TimeUtilsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A provider for advanced time manipulation functions including RFC3339 parsing and strftime formatting.",
	}
}

func (p *TimeUtilsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No configuration needed for this provider
}

func (p *TimeUtilsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *TimeUtilsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *TimeUtilsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function { return NewUnixTimestampFunction() },
		func() function.Function { return NewStrftimeFunction() },
		func() function.Function { return NewDaysDifferenceFunction() },
		func() function.Function { return NewParseRFC3339Function() },
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TimeUtilsProvider{
			version: version,
		}
	}
}

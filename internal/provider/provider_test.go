/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	internal "github.com/metio/terraform-provider-migadu/internal/provider"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigaduProvider_Metadata(t *testing.T) {
	t.Parallel()
	p := &internal.MigaduProvider{}
	request := provider.MetadataRequest{}
	response := &provider.MetadataResponse{}
	p.Metadata(context.TODO(), request, response)

	assert.Equal(t, "migadu", response.TypeName, "TypeName")
}

func providerConfig(endpoint string) string {
	return fmt.Sprintf(`
		provider "migadu" {
			username = "username"
			token    = "token"
			endpoint = "%s"
		}
	`, endpoint)
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"migadu": providerserver.NewProtocol6WithError(internal.New()),
	}
)

type ConfigurationErrorTestCase struct {
	Configuration string
	ErrorRegex    string
}

type APIErrorTestCase struct {
	StatusCode int
	ErrorRegex string
}

type ResourceTestCase[T any] struct {
	Create       ResourceTestStep[T]
	Update       ResourceTestStep[T]
	ImportIgnore []string
}

type ResourceTestStep[T any] struct {
	Send T
	Want T
}

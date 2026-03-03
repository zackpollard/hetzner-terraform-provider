// Copyright (c) Zack
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccBootLinux_CRUD tests activating, updating, and deactivating Linux install boot.
func TestAccBootLinux_CRUD(t *testing.T) {
	serverNumber := testAccGetOrCreateServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Activate Linux install.
			{
				Config: testAccBootLinuxConfig(serverNumber, "Debian 12 base", "en"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hetzner_boot_linux.test", "active", "true"),
					resource.TestCheckResourceAttrSet("hetzner_boot_linux.test", "password"),
				),
			},
			// Import.
			{
				ResourceName:                         "hetzner_boot_linux.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "server_number",
				// password and authorized_key are not returned on read
				ImportStateVerifyIgnore: []string{"password", "authorized_key", "arch"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hetzner_boot_linux.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return rs.Primary.Attributes["server_number"], nil
				},
			},
			// Destroy deactivates.
		},
	})
}

// TestAccBootLinux_DataSource reads Linux boot config via data source.
func TestAccBootLinux_DataSource(t *testing.T) {
	serverNumber := testAccGetOrCreateServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBootLinuxDataSourceConfig(serverNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hetzner_boot_linux.test", "server_number"),
				),
			},
		},
	})
}

// TestAccBootVNC_DataSource reads VNC boot config via data source.
func TestAccBootVNC_DataSource(t *testing.T) {
	serverNumber := testAccGetOrCreateServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBootVNCDataSourceConfig(serverNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hetzner_boot_vnc.test", "server_number"),
				),
			},
		},
	})
}

// TestAccBootRescue_Import tests importing an existing rescue boot configuration.
func TestAccBootRescue_Import(t *testing.T) {
	serverNumber := testAccGetOrCreateServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Activate rescue.
			{
				Config: testAccBootRescueConfig(serverNumber, "linux"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hetzner_boot_rescue.test", "active", "true"),
				),
			},
			// Import.
			{
				ResourceName:                         "hetzner_boot_rescue.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "server_number",
				ImportStateVerifyIgnore:              []string{"password", "authorized_key", "arch", "os", "keyboard"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hetzner_boot_rescue.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return rs.Primary.Attributes["server_number"], nil
				},
			},
		},
	})
}

package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixSubnetsDataSource_basic(t *testing.T) {
	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetsExists("data.nutanix_subnets.test1"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}

func testAccSubnetsDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"

	cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  	}

	vlan_id = %d
	subnet_type = "VLAN"
	
	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]
	
	dhcp_options {
		boot_file_name = "bootfile"
		tftp_server_name = "192.168.0.252"
		domain_name = "nutanix"
	}

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list = ["nutanix.com", "calm.io"]
	
}

data "nutanix_subnet" "test" {
	subnet_id = "${nutanix_subnet.test.id}"
}

data "nutanix_subnets" "test1" {
	metadata {
		length = 1
	}
}`, r, r)
}
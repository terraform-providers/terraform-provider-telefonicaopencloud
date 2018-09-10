provider "telefonicaopencloud" {
  user_name   = "c2c_admin"
  domain_name = "Huawei China"
  password    = "Newuser@123"
  auth_url    = "https://iam.na-mexico-1.telefonicaopencloud.com/v3"
  region      = "na-mexico-1"
  tenant_id   = "d3369c1b9cfa4956825838293f9c0e7d"
}

//stack
resource "telefonicaopencloud_rts_stack_v1" "stack_1" {
  name = "new-stack"
  disable_rollback= true
  timeout_mins=60
  template_body = <<JSON
          {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}

//Stack DS
data "telefonicaopencloud_rts_stack_v1" "stacks" {
        name = "new-stack"
}

output "resource_data3" {
    value = "${data.telefonicaopencloud_rts_stack_v1.stacks.parameters}"
}

output "resource_data4" {
    value = "${data.telefonicaopencloud_rts_stack_v1.stacks.template_body}"
}

//Stack Resource DS
data "telefonicaopencloud_rts_stack_resource_v1" "resource_1" {
  stack_name = "new-stack"
  resource_name = "random"
}


// Config
resource "telefonicaopencloud_rts_software_config_v1" "config"{
        name = "rts-config"
        output_values = [{
          type = "String"
          name = "result"
          error_output = "false"
          description = "value1"
        }]
        input_values=[{
          default = "0"
          type = "String"
          name = "foo"
          description = "value2"
        }]
        group = "script"
   }


data "telefonicaopencloud_rts_software_config_v1" "ds_config" {
  id = "d6085150-acd0-4196-955f-adeea8c1b7aa"
}

output "resource_data11" {
    value = "${data.telefonicaopencloud_rts_software_config_v1.ds_config.input_values}"
}

output "resource_data12" {
    value = "${data.telefonicaopencloud_rts_software_config_v1.ds_config.output_values}"
}

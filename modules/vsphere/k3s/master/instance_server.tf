data "vsphere_datacenter" "datacenter" {
  name = var.vsphere_datacenter
}

data "vsphere_datastore" "datastore" {
  name          = var.vsphere_datastore
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_compute_cluster" "cluster" {
  name          = var.vsphere_cluster
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_network" "network" {
  name          = var.vsphere_network
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

data "vsphere_virtual_machine" "template" {
  name          = var.vm_template
  datacenter_id = data.vsphere_datacenter.datacenter.id
}

resource "vsphere_virtual_machine" "master" {
  name             = "${var.resource_name}-${local.resource_tag}-server1"
  resource_pool_id = data.vsphere_compute_cluster.cluster.resource_pool_id
  folder           = var.vsphere_cluster_folder
  datastore_id     = data.vsphere_datastore.datastore.id
  num_cpus         = var.vm_num_cpus
  memory           = var.vm_memory
  guest_id         = data.vsphere_virtual_machine.template.guest_id
  scsi_type        = data.vsphere_virtual_machine.template.scsi_type
  network_interface {
    network_id   = data.vsphere_network.network.id
    adapter_type = data.vsphere_virtual_machine.template.network_interface_types[0]
  }
  disk {
    label            = data.vsphere_virtual_machine.template.disks.0.label
    size             = data.vsphere_virtual_machine.template.disks.0.size
    thin_provisioned = data.vsphere_virtual_machine.template.disks.0.thin_provisioned
  }
  clone {
    template_uuid = data.vsphere_virtual_machine.template.id
    # customize {
    #   linux_options {
    #     host_name = "${var.resource_name}"
    #     domain = ""
    #   }
  }
}


# resource "vsphere_virtual_machine" "master" {
#   name             = "${var.resource_name}-${local.resource_tag}-server1"
#   resource_pool_id = data.vsphere_compute_cluster.cluster.resource_pool_id
#   folder           = var.vsphere_cluster_folder
#   datastore_id     = data.vsphere_datastore.datastore.id
#   num_cpus         = var.vm_num_cpus
#   memory           = var.vm_memory
#   guest_id         = var.vm_guest_os_id
#   enable_disk_uuid = true
#   network_interface {
#     network_id = data.vsphere_network.network.id
#   }
#   disk {
#     label = var.vm_disk_label
#     size  = var.vm_disk_size
#   }


#   cdrom {
#     datastore_id = data.vsphere_datastore.datastore.id
#     path         = ""
#   }
# }

locals {
  resource_tag     =  "distros-qa"
}
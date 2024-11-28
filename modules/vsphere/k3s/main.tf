module master {
   source ="./master"

   vsphere_user           = var.vsphere_user
   vsphere_password       = var.vsphere_password
   vsphere_server         = var.vsphere_server
   vsphere_datastore      = var.vsphere_datastore
   vsphere_datacenter     = var.vsphere_datacenter
   vsphere_cluster        = var.vsphere_cluster
   vsphere_cluster_folder = var.vsphere_cluster_folder
   vsphere_network        = var.vsphere_network

   ## VM related
   resource_name          = var.resource_name
   vm_template            = var.vm_template
   vm_num_cpus            = var.vm_num_cpus
   vm_memory              = var.vm_memory
   vm_guest_os_id         = var.vm_guest_os_id
   vm_disk_label          = var.vm_disk_label
   vm_disk_size           = var.vm_disk_size

}

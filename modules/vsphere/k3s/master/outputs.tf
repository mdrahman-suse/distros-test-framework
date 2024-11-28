output "master_ips" {
  value = vsphere_virtual_machine.master.default_ip_address
  description = "The IP of the vSphere node"
}
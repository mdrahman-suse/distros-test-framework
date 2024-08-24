output "master_ips" {
  value       = module.instance.master_ips
  description = "The private IP of the AWS node"
}

output "worker_ips" {
  value       = module.instance.worker_ips
  description = "The private IP of the AWS node"
}

output "windows_worker_ips" {
  value       = module.instance.windows_worker_ips
  description = "The public IP of the AWS node"
}

output "windows_worker_password" {
  value = module.instance.windows_worker_password
  description = "Windows worker password"
}

output "bastion_ip" {
  value       = module.instance.bastion_ip
  description = "The public IP of the AWS node"
}

output "bastion_dns" {
  value       = module.instance.bastion_dns
  description = "The public DNS name of the AWS node"
}

output "check_airgap" {
  value = module.instance.check_airgap.rendered
}

output "check_ipv6only" {
  value = module.instance.check_ipv6only.rendered
}
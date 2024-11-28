output "master_ips" {
  value       = module.master.master_ips
  description = "The public IP of the AWS node"
}
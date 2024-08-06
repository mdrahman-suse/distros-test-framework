data "template_file" "is_airgap" {
    template = (var.enable_public_ip == false && var.enable_ipv6 == false) ? true : false
}

data "template_file" "is_ipv6only" {
    template = (var.enable_public_ip == false && var.enable_ipv6 == true) ? true : false
}

resource "aws_instance" "master" {

  depends_on = [ null_resource.prepare_bastion ]

  ami                         = var.aws_ami
  instance_type               = var.ec2_instance_class  
  associate_public_ip_address = false
  ipv6_address_count          = var.enable_ipv6 == true ? 1 : 0
  count                       = var.no_of_server_nodes
  
  root_block_device {
    volume_size          = var.volume_size
    volume_type          = "standard"
  }
  subnet_id              = var.subnets
  availability_zone      = var.availability_zone
  vpc_security_group_ids = [var.sg_id]
  key_name               = var.key_name
  tags = {
    Name                 = "${var.resource_name}-server${count.index + 1}"
  }

}

resource "aws_instance" "worker" {

  depends_on = [ null_resource.prepare_bastion ]

  ami                         = var.aws_ami
  instance_type               = var.ec2_instance_class  
  associate_public_ip_address = false
  ipv6_address_count          = var.enable_ipv6 == true ? 1 : 0
  count                       = var.no_of_worker_nodes
  
  root_block_device {
    volume_size          = var.volume_size
    volume_type          = "standard"
  }
  subnet_id              = var.subnets
  availability_zone      = var.availability_zone
  vpc_security_group_ids = [var.sg_id]
  key_name               = var.key_name
  tags = {
    Name                 = "${var.resource_name}-worker${count.index + 1}"
  }
}

# resource "aws_instance" "windows_worker" {
#   ami                         = var.windows_aws_ami
#   instance_type               = var.windows_ec2_instance_class  
#   associate_public_ip_address = false
#   ipv6_address_count          = var.enable_ipv6 ? 1 : 0
#   count                       = ((data.template_file.is_airgap == true || data.template_file.is_ipv6only == true)) ? 1 : 0
  
#   root_block_device {
#     volume_size          = "50"
#     volume_type          = "standard"
#   }
#   subnet_id              = var.subnets
#   availability_zone      = var.availability_zone
#   vpc_security_group_ids = [var.sg_id]
#   key_name               = var.key_name
#   tags = {
#     Name                 = "${var.resource_name}-agent"
#   }
# }

resource "aws_instance" "bastion" {
  ami                         = var.aws_ami
  instance_type               = var.ec2_instance_class  
  associate_public_ip_address = true
  ipv6_address_count          = var.enable_ipv6 == true ? 1 : 0
  count                       = var.no_of_bastion_nodes
  
  connection {
    type          = "ssh"
    user          = var.aws_user
    host          = self.public_ip
    private_key   = file(var.access_key)
  }
  root_block_device {
    volume_size          = var.volume_size
    volume_type          = "standard"
  }
  subnet_id              = var.bastion_subnets
  availability_zone      = var.availability_zone
  vpc_security_group_ids = [var.sg_id]
  key_name               = var.key_name
  tags = {
    Name                 = "${var.resource_name}-bastion"
  }

  provisioner "file" {
    source = "../../config/.ssh/aws_key.pem"
    destination = "/tmp/${var.key_name}.pem"
  }

  provisioner "file" {
    source = "setup/get_artifacts.sh"
    destination = "/tmp/get_artifacts.sh"
  }

  provisioner "file" {
    source = "setup/install_product.sh"
    destination = "/tmp/install_product.sh"
  }

  provisioner "file" {
    source = "setup/bastion_prepare.sh"
    destination = "/tmp/bastion_prepare.sh"
  }

  provisioner "file" {
    source = "setup/private_registry.sh"
    destination = "/tmp/private_registry.sh"
  }

  provisioner "file" {
    source = "setup/basic-registry"
    destination = "/tmp"
  }
  
  #sudo bastion_prepare.sh
  #sudo /tmp/download_product.sh ${var.product} ${var.product_version} ${var.arch}

}

resource "null_resource" "prepare_bastion" {

  depends_on = [ aws_instance.bastion[0] ]
  connection {
    type          = "ssh"
    user          = var.aws_user
    host          = aws_instance.bastion[0].public_ip
    private_key   = file(var.access_key)
  }

  provisioner "remote-exec" {
    inline = [<<-EOT
      echo ${aws_instance.bastion[0].public_ip} > /tmp/${var.resource_name}_bastion_ip
      sudo cp /tmp/${var.key_name}.pem /tmp/bastion_prepare.sh /tmp/get_artifacts.sh /tmp/install_product.sh /tmp/private_registry.sh ~/
      sudo cp -r /tmp/basic-registry ~/
      sudo chmod +x bastion_prepare.sh
      sudo ./bastion_prepare.sh
    EOT
    ]
  }
}


# resource "null_resource" "uploader" {
#   depends_on = [ 
#     aws_instance.bastion,
#     aws_instance.server,
#     aws_instance.agent 
#   ]

#   count = (tobool(data.template_file.is_airgap.rendered) == true || tobool(data.template_file.is_ipv6only.rendered) == true) ? 1 : 0

#   connection {
#     type          = "ssh"
#     user          = var.aws_user
#     host          = "${aws_instance.bastion[0].public_ip}"
#     private_key   = file(var.access_key)
#   }

#   provisioner "remote-exec" {
#     inline = [<<-EOT
#       chmod +x /tmp/download_product.sh
#       sudo /tmp/download_product.sh ${var.product} ${var.product_version} ${var.arch} "${var.channel}"
#       chmod 400 /tmp/${var.key_name}.pem
#       scp ${local.shell_options} -i /tmp/${var.key_name}.pem /tmp/install_product.sh ${var.aws_user}@${var.enable_ipv6 ? format("[%#v]",aws_instance.server[0].ipv6_addresses[0]) : aws_instance.server[0].private_ip}:/tmp/install_product.sh
#       scp ${local.shell_options} -i /tmp/${var.key_name}.pem /tmp/${var.resource_name}_bastion_ip ${var.aws_user}@${var.enable_ipv6 ? format("[%#v]",aws_instance.server[0].ipv6_addresses[0]) : aws_instance.server[0].private_ip}:/tmp/${var.resource_name}_bastion_ip
#       scp -r ${local.shell_options} -i /tmp/${var.key_name}.pem /tmp/${var.product}-assets ${var.aws_user}@${var.enable_ipv6 ? format("[%#v]",aws_instance.server[0].ipv6_addresses[0]) : aws_instance.server[0].private_ip}:/tmp/${var.product}-assets
#       scp ${local.shell_options} -i /tmp/${var.key_name}.pem /tmp/install_product.sh ${var.aws_user}@${var.enable_ipv6 ? format("[%#v]",aws_instance.agent[0].ipv6_addresses[0]) : aws_instance.agent[0].private_ip}:/tmp/install_product.sh
#       scp ${local.shell_options} -i /tmp/${var.key_name}.pem /tmp/${var.resource_name}_bastion_ip ${var.aws_user}@${var.enable_ipv6 ? format("[%#v]",aws_instance.agent[0].ipv6_addresses[0]) : aws_instance.agent[0].private_ip}:/tmp/${var.resource_name}_bastion_ip
#       scp -r ${local.shell_options} -i /tmp/${var.key_name}.pem /tmp/${var.product}-assets ${var.aws_user}@${var.enable_ipv6 ? format("[%#v]",aws_instance.agent[0].ipv6_addresses[0]) : aws_instance.agent[0].private_ip}:/tmp/${var.product}-assets
#      EOT
#      ]
#   }
# }

# resource "null_resource" "installer" {
#   depends_on = [ 
#     aws_instance.bastion,
#     aws_instance.server,
#     aws_instance.agent,
#     null_resource.uploader
#   ]

#   count = (tobool(data.template_file.is_airgap.rendered) == true || tobool(data.template_file.is_ipv6only.rendered) == true) ? 1 : 0

#   connection {
#     type          = "ssh"
#     user          = var.aws_user
#     host          = "${aws_instance.bastion[0].public_ip}"
#     private_key   = file(var.access_key)
#   }

#   provisioner "remote-exec" {
#     inline = [<<-EOT
#       chmod 400 /tmp/${var.key_name}.pem
#       ssh -i /tmp/${var.key_name}.pem ${local.shell_options} ${var.aws_user}@${var.enable_ipv6 ? aws_instance.server[0].ipv6_addresses[0] : aws_instance.server[0].private_ip} 'sudo cp -p /tmp/${var.product}-assets/${var.product} /usr/local/bin/${var.product} && chmod +x /usr/local/bin/${var.product}'
#       ssh -i /tmp/${var.key_name}.pem ${local.shell_options} ${var.aws_user}@${var.enable_ipv6 ? aws_instance.agent[0].ipv6_addresses[0] : aws_instance.agent[0].private_ip} 'sudo cp -p /tmp/${var.product}-assets/${var.product} /usr/local/bin/${var.product} && chmod +x /usr/local/bin/${var.product}'
#      EOT
#      ]
#   }
  
# }

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  # Ensure that enable_dns_support and enable_dns_hostnames are true for Aurora DNS resolution
  enable_dns_support   = true
  enable_dns_hostnames = true
}


data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_subnet" "private" {
  count = 2 # Limit to 2 availability zones

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  availability_zone = element(data.aws_availability_zones.available.names, count.index)

  map_public_ip_on_launch = false
}





provider "random" {

}

resource "aws_secretsmanager_secret" "fga_rds_db_credentials" {
  name = "/STORER/FGA/RDS"
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "!#$&"
}


resource "aws_secretsmanager_secret_version" "fga_rds_db_credentials_version" {
  secret_id     = aws_secretsmanager_secret.fga_rds_db_credentials.id
  secret_string = jsonencode({
    username = "rdsroot"
    password = random_password.password.result
  })
}


resource "aws_db_subnet_group" "fga_subnet_group" {
  name       = "fga_subnet_group"
  subnet_ids = aws_subnet.private.*.id

  tags = {
    Name = "FGA DB Subnet Group"
  }
}

/*
data "aws_secretsmanager_secret_version" "rds_secret" {
  secret_id = "/STORER/FGA/RDS"
}
*/

resource "aws_security_group" "aurora_sg" {
  name   = "aurora-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port = 0
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block] # Allow access from within the VPC only
  }
}

resource "aws_rds_cluster" "aurora_postgres_cluster" {
  cluster_identifier          = "openfga-db"
  engine                      = "aurora-postgresql"
  //engine_mode                 = "serverless"
  engine_version              = "15.4"  # Specify your desired version
  db_subnet_group_name        = aws_db_subnet_group.fga_subnet_group.name
  vpc_security_group_ids      = [aws_security_group.aurora_sg.id]
  skip_final_snapshot         = true
  deletion_protection         = false
  database_name               = "openfga"
  master_username             = jsondecode(aws_secretsmanager_secret_version.fga_rds_db_credentials_version.secret_string)["username"]
  master_password             = jsondecode(aws_secretsmanager_secret_version.fga_rds_db_credentials_version.secret_string)["password"]

  # Enable Performance Insights
  #enable_performance_insights = true
  #performance_insights_retention_period = 7  # days, choose between 7 and 731

  # Serverless specific settings
  serverlessv2_scaling_configuration {
    max_capacity = 1.0
    min_capacity = 0.5
  }
 
}

resource "aws_rds_cluster_instance" "cluster_instances" {
  identifier = "openfga-db-serverless"
  cluster_identifier = aws_rds_cluster.aurora_postgres_cluster.id
  instance_class = "db.serverless"
  engine = aws_rds_cluster.aurora_postgres_cluster.engine
  engine_version = aws_rds_cluster.aurora_postgres_cluster.engine_version
}

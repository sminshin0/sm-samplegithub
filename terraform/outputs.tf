# 중요한 정보들을 출력
output "instance_id" {
  description = "EC2 인스턴스 ID"
  value       = aws_instance.web.id
}

output "public_ip" {
  description = "EC2 인스턴스 공개 IP"
  value       = aws_instance.web.public_ip
}

output "public_dns" {
  description = "EC2 인스턴스 공개 DNS"
  value       = aws_instance.web.public_dns
}

output "website_url" {
  description = "웹사이트 URL"
  value       = "http://${aws_instance.web.public_ip}"
}
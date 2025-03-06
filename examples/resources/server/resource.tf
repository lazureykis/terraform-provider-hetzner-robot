resource "hetzner_robot_server" "example" {
  id   = "123456" # Your server ID
  name = "terraform-managed-server"

  # Optional: Enable rescue mode
  # rescue = true
  # rescue_os = "linux64"
}

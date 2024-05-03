terraform {
  required_providers {
    email = {
      version = "0.2.3"
      source  = "watonomous.ca/tf/email"
    }
  }
}

provider "email" {}

variable "recipients" {
  default = ["recipient1@example.com", "recipient2@example.com", "recipient3@example.com", "recipient4@example.com", "recipient5@example.com", "recipient6@example.com"]
}

resource "email_email" "example" {
  # to = "infra-outreach@watonomous.ca"
  for_each = toset(var.recipients)

  to = each.value
  from = "sentry-outgoing@watonomous.ca"
  reply_to = "infrastructure@watonomous.ca"
  subject = "Hello from Terraform"
  body = "This is a test email sent from Terraform using a custom email provider."
  # smtp_server = "smtp.gmail.com"
  # smtp_port = "587"
  smtp_server = "localhost"
  smtp_port = "2525"
  smtp_username = "mailbot@watonomous.ca"
  smtp_password = "<replace_me>"
}

resource "email_email" "example_with_styling" {
  to = "infra-outreach@watonomous.ca"
  from = "sentry-outgoing@watonomous.ca"
  reply_to = "infrastructure@watonomous.ca"
  subject = "Hello from Terraform"
  preamble = <<EOT
MIME-Version: 1.0
Content-Type: text/html; charset="utf-8"
EOT
  body = <<EOT
<!DOCTYPE html>
<html>
<head>
<title>Welcome to WATonomous!</title>
</head>
<body>
<pre style="font-family: 'Courier New', monospace;">
    _∩_
  __|_|_
 /|__|__\____
|            |
`.(o)-----(o).'
</pre>

Test email sent from Terraform using a custom email provider.

</body>
</html>
EOT
  # smtp_server = "smtp.gmail.com"
  # smtp_port = "587"
  smtp_server = "localhost"
  smtp_port = "2525"
  smtp_username = "mailbot@watonomous.ca"
  smtp_password = "<replace_me>"
}
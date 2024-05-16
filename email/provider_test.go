package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEmail_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"email": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccEmailConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailExists("email_email.example"),
					resource.TestCheckResourceAttr(
						"email_email.example", "to", "recipient@example.com"),
					resource.TestCheckResourceAttr(
						"email_email.example", "from", "sender@example.com"),
					resource.TestCheckResourceAttr(
						"email_email.example", "reply_to", "reply_to@example.com"),
					resource.TestCheckResourceAttr(
						"email_email.example", "subject", "Test Subject"),
					resource.TestCheckResourceAttr(
						"email_email.example", "body", "Test Body"),
					resource.TestCheckResourceAttr(
						"email_email.example", "smtp_server", "localhost"),
					resource.TestCheckResourceAttr(
						"email_email.example", "smtp_port", "2525"),
					resource.TestCheckResourceAttr(
						"email_email.example", "smtp_username", "username"),
					resource.TestCheckResourceAttr(
						"email_email.example", "smtp_password", "password"),
				),
			},
		},
	})
}

func testAccCheckEmailExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		return nil
	}
}
func mockSendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return errors.New("421 Service not available")
}

func TestRetryWorkflow(t *testing.T) {
	maxRetries := 5

	userName := "username"
	password := "password"

	server := "localhost"
	port := "2525"

	to := "hello@example.com"
	from := "wato@example.com"
	msg := "hello"

	// TO-DO: Create Struct With Errors, Expected Values, Iterate Through Them

	// grab error & test error code
	err := sendMail(mockSendMail, maxRetries, server, port, userName, password, from, to, msg)
	errCode := extractStatusCode(err.Error())
	// Error Guard Statements
	if err == nil {
		t.Errorf("Expected: %s, Got: %s", "nil", err)
	}
	if errCode != "421" {
		t.Errorf("Expected: %s, Got: %s", "421", errCode)
	}

}

// Requires a local SMTP server running on port 2525
// `docker run --rm -it -p 3000:80 -p 2525:25 rnwood/smtp4dev:v3`
const testAccEmailConfig = `
resource "email_email" "example" {
	to = "recipient@example.com"
	from = "sender@example.com"
	reply_to = "reply_to@example.com"
	subject = "Test Subject"
	body = "Test Body"
	smtp_server = "localhost"
	smtp_port = "2525"
	smtp_username = "username"
	smtp_password = "password"
}
`

package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// Global Testing Variables
var sendMailInvocations int = 0

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
func mockSendMailReturn421(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	sendMailInvocations += 1
	return errors.New("421 Service not available")
}

// function to test non-421 error
func mockSendMailReturn500(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	sendMailInvocations += 1
	return errors.New("500 Internal Server Error")
}

// email to return no error
func mockSendMailSuccess(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	sendMailInvocations += 1
	return nil
}

func TestRetryWorkflow(t *testing.T) {
	// define max retries
	maxRetries := 5
	// write the tests
	tests := []struct {
		name                        string
		sendMailFunc                SendMailFunc
		expectedErr                 error
		expectedSendMailInvocations int
	}{
		{
			name:                        "Retry with 421 error",
			sendMailFunc:                mockSendMailReturn421,
			expectedErr:                 errors.New("421 Service not available"),
			expectedSendMailInvocations: maxRetries,
		},
		{
			name:                        "Success on first try",
			sendMailFunc:                mockSendMailSuccess,
			expectedErr:                 nil,
			expectedSendMailInvocations: 1,
		},
		{
			name:                        "Other error",
			sendMailFunc:                mockSendMailReturn500,
			expectedErr:                 errors.New("500 Internal Server Error"),
			expectedSendMailInvocations: 1,
		},
	}
	// execute the tests
	for _, test := range tests {
		// run subtests
		t.Run(test.name, func(t *testing.T) {
			// set test retries zero before each test
			sendMailInvocations = 0
			err := sendMail(test.sendMailFunc, maxRetries, "localhost", "2525", "username", "password", "from@example.com", "to@example.com", "message")
			if test.expectedErr != nil {
				// assert that the errors are equal
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			// assert testing number invocations regardless
			assert.Equal(t, sendMailInvocations, test.expectedSendMailInvocations)
		})
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

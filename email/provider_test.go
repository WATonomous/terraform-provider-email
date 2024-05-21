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

// function to test non-421 error
func mockReturnNon421(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return errors.New("500 Internal Server Error")
}

// email to return no error
func mockSendMailSuccess(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return nil
}

func TestRetryWorkflow(t *testing.T) {
	// define max retries
	maxRetries := 5
	// write the tests
	tests := []struct {
		name         string
		sendMailFunc SendMailFunc
		expectedErr  error
	}{
		{
			name:         "Retry with 421 error",
			sendMailFunc: mockSendMail,
			expectedErr:  errors.New("421 Service not available"),
		},
		{
			name:         "Success on first try",
			sendMailFunc: mockSendMailSuccess,
			expectedErr:  nil,
		},
		{
			name:         "Other error",
			sendMailFunc: mockReturnNon421,
			expectedErr:  errors.New("500 Internal Server Error"),
		},
	}
	// execute the tests
	for _, test := range tests {
		// run subtests
		t.Run(test.name, func(t *testing.T) {
			err := sendMail(test.sendMailFunc, maxRetries, "localhost", "2525", "username", "password", "from@example.com", "to@example.com", "message")
			if test.expectedErr != nil {
				// assert that the errors are equal
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
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

package email

import (
	"log"
	"math/rand"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// define function type
type SendMailFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

func resourceEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailCreate,
		Read:   resourceEmailRead,
		Update: resourceEmailUpdate,
		Delete: resourceEmailDelete,

		Schema: map[string]*schema.Schema{
			"to": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"to_display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"from": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"from_display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"reply_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"subject": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"preamble": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"body": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"smtp_server": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"smtp_port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"smtp_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"smtp_password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"dry_run": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

// regex function to extract the status code
func extractStatusCode(errMsg string) string {
	// Regex to find the first three-digit number, which is the SMTP status code
	re := regexp.MustCompile(`\b\d{3}\b`)
	matches := re.FindString(errMsg)
	if matches != "" {
		return matches // Returns the first match (three-digit number) if found
	}
	return "No status code found"
}

func resourceEmailCreate(d *schema.ResourceData, m interface{}) error {
	rawTo := d.Get("to").([]interface{})
	to := make([]string, len(rawTo))
	for i, v := range rawTo {
		to[i] = v.(string)
	}
	toDisplayName := d.Get("to_display_name").(string)
	from := d.Get("from").(string)
	fromDisplayName := d.Get("from_display_name").(string)
	replyTo := d.Get("reply_to").(string)
	subject := d.Get("subject").(string)
	preamble := d.Get("preamble").(string)
	body := d.Get("body").(string)
	smtpServer := d.Get("smtp_server").(string)
	smtpPort := d.Get("smtp_port").(string)
	smtpUsername := d.Get("smtp_username").(string)
	smtpPassword := d.Get("smtp_password").(string)
	dryRun := d.Get("dry_run").(bool)

	if toDisplayName == "" {
		toDisplayName = strings.Join(to, ", ")
	}
	if fromDisplayName == "" {
		fromDisplayName = from
	}

	msg := "From: " + fromDisplayName + "\n" +
		"To: " + toDisplayName + "\n" +
		"Reply-To: " + replyTo + "\n" +
		"Subject: " + subject + "\n" +
		preamble + "\n\n" +
		body

	if dryRun {
		log.Printf("In dry-run mode. Not sending email from %s to %s with subject %s", from, to, subject)
	} else {
		// TODO: make this tf configurable
		maxRetries := 5
		// send mail using exponential back-off
		err := sendMail(smtp.SendMail, maxRetries, smtpServer, smtpPort, smtpUsername, smtpPassword, from, to, msg)
		// log error if not cleared after retries
		if err != nil {
			log.Printf("smtp error: %s", err)
			return err
		}
	}

	timestamp := time.Now().Unix()
	d.SetId(strings.Join(to, ",") + " | " + subject + " | " + strconv.FormatInt(timestamp, 10))

	return resourceEmailRead(d, m)
}

func sendMail(sendEmailImpl SendMailFunc, maxRetries int, smtpServer string, smtpPort string, smtpUsername string, smtpPassword string, from string, to []string, msg string) error {
	// Set up a random number for exponential backoff
	minRandInt := 10
	maxRandInt := 150
	// generate random number in that range
	randomNumber := rand.Intn(maxRandInt-minRandInt) + minRandInt
	var err error
	for retries := 0; retries < maxRetries; retries++ {
		// send smtp email
		err = sendEmailImpl(smtpServer+":"+smtpPort,
			smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer),
			from, to, []byte(msg))

		if err == nil {
			break
		}
		// extract error code
		errorCode := extractStatusCode(err.Error())
		log.Printf("Extracted Error Code: %s", errorCode)
		// guard statement for error 421
		if errorCode != "421" {
			break
		}
		// implement exponential back off
		time.Sleep(time.Duration(randomNumber << 1))
	}
	return err
}

func resourceEmailRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceEmailDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func resourceEmailUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceEmailRead(d, m)
}

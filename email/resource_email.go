package email

import (
	"log"
	"net/smtp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailCreate,
		Read:   resourceEmailRead,
		Update: resourceEmailUpdate,
		Delete: resourceEmailDelete,

		Schema: map[string]*schema.Schema{
			"to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func resourceEmailCreate(d *schema.ResourceData, m interface{}) error {
	to := d.Get("to").(string)
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

	if toDisplayName == "" {
		toDisplayName = to
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

	err := smtp.SendMail(smtpServer+":"+smtpPort,
		smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}

	// Create unique ID using current timestamp
	timestamp := time.Now().Unix()
	d.SetId(to + " | " + subject + " | " + strconv.FormatInt(timestamp, 10))
	return resourceEmailRead(d, m)
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

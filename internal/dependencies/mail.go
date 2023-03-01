package dependencies

import (
	"log"

	"github.com/dafaath/iot-server/v2/configs"
	"gopkg.in/gomail.v2"
)

func NewMailDialer(config *configs.Config) (*gomail.Dialer, error) {
	dialer := gomail.NewDialer(
		config.Mail.SMTPHost,
		config.Mail.SMTPPort,
		config.Mail.AuthenticationMail,
		config.Mail.AuthenticationPassword,
	)
	log.Println("Dialing mail server...")
	// sendCloser, err := dialer.Dial()
	// if err != nil {
	// 	return dialer, err
	// }
	// err = sendCloser.Close()
	// if err != nil {
	// 	return dialer, err
	// }
	log.Println("Finish dialing mail server...")

	return dialer, nil
}

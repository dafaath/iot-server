package dependencies

import (
	"github.com/dafaath/iot-server/configs"
	"gopkg.in/gomail.v2"
)

func NewMailDialer(config *configs.Config) (*gomail.Dialer, error) {
	dialer := gomail.NewDialer(
		config.Mail.SMTPHost,
		config.Mail.SMTPPort,
		config.Mail.AuthenticationMail,
		config.Mail.AuthenticationPassword,
	)

	// log.Println("Testing mail dialer...")
	// sendCloser, err := dialer.Dial()
	// if err != nil {
	// 	return dialer, err
	// }
	// err = sendCloser.Close()
	// if err != nil {
	// 	return dialer, err
	// }
	// log.Println("Success testing mail dialer")

	return dialer, nil
}

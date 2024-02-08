package mail

import (
	"petpal-backend/src/configs"
)

func SendEmailWithGmail(toEmailAddress string) error {
	// Add a function to validate an Email address
	sender := NewGmailSender("PetpalAdmin", configs.GetEmailSenderAddress(), configs.GetEmailSenderPassword())

	subject := "Test Petpal Email Sender"
	content := `
	<h1>สวัสดีครับ</h1>
	<p>ยินดีด้วย คุณยืนยันตัวตนกับทาง Petpal สำเร็จ!!!</p>
	`
	to := []string{toEmailAddress}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

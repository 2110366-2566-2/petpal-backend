package mail

import (
	"petpal-backend/src/configs"
)

// Send an email to "toEmailAddress" with "emailSubject" and "emailContent"
func SendEmailWithGmail(toEmailAddress string, emailSubject string, emailContent string) error {
	// Add a function to validate an Email address
	sender := NewGmailSender("PetpalAdmin", configs.GetInstance().GetEmailSenderAddress(), configs.GetInstance().GetEmailSenderPassword())
	/* Example emailContent
	`
	<h4>สวัสดีครับ</h4>
	<p>ยินดีด้วย คุณยืนยันตัวตนกับทาง Petpal สำเร็จ!!!</p>
	`
	*/
	to := []string{toEmailAddress}
	err := sender.SendEmail(emailSubject, emailContent, to, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

package utills

import (
	"fmt"
	"petpal-backend/src/models"
	mail "petpal-backend/src/utills/email"
)

// User create booking
func NotifyCreateBooking(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking confirmed"
	emailContent := fmt.Sprintf(`<h4>สวัสดีครับ</h4>
<p>ผู้ใช้บริการ %s<br>
ได้สร้างคำขอการจองบริการ %s<br>
ในเวลา %s ถึงเวลา %s`,
		user.FullName,
		booking.ServiceName,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)

	emailContent += `<br><br>โปรดตรวจสอบความถูกต้องของการจอง แล้วยืนยันการจองในระบบ<br><br>ขอขอบคุณที่ใช้บริการกับเรา<br><br>Petpal Team</p>`
	err := mail.SendEmailWithGmail(svcp.SVCPEmail, emailSubject, emailContent)
	if err != nil {
		return err
	}
	return nil
}

// User reschedule booking
func NotifyRescheduleBooking(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking scheduled changed"
	emailContent := fmt.Sprintf(`<h4>สวัสดีครับ</h4>
<p>ผู้ใช้บริการ %s<br>
ได้เปลี่ยนเวลาการจองบริการ %s<br>
เป็นเวลา %s ถึงเวลา %s`,
		user.FullName,
		booking.ServiceName,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)

	emailContent += `<br><br>โปรดตรวจสอบความถูกต้องของการจอง แล้วยืนยันการจองในระบบ<br><br>ขอขอบคุณที่ใช้บริการกับเรา<br><br>Petpal Team</p>`
	err := mail.SendEmailWithGmail(svcp.SVCPEmail, emailSubject, emailContent)
	if err != nil {
		return err
	}
	return nil
}

// SVCP confirm booking
func NotifyConfirmBookingToUser(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking confirmed"
	emailContent := fmt.Sprintf(`<h4>สวัสดีครับ</h4>
<p>ผู้ให้บริการ %s<br>
ได้ยืนยันการจองบริการ %s<br>
ในเวลา %s ถึงเวลา %s<br>
เรียบร้อยแล้ว`,
		svcp.SVCPUsername,
		booking.ServiceName,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)

	emailContent += `<br><br>ขอขอบคุณที่ใช้บริการกับเรา<br><br>Petpal Team</p>`
	err := mail.SendEmailWithGmail(user.Email, emailSubject, emailContent)
	if err != nil {
		return err
	}
	return nil
}

// Cancel booking
func NotifyCancelBookingToUser(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking cancelled"
	emailContent := fmt.Sprintf(`<h4>สวัสดีครับ</h4>
<p>การจองบริการ %s<br>
กับผู้ให้บริการ %s<br>
เวลา %s ถึงเวลา %s<br>
ได้ถูกยกเลิกแล้ว`,
		booking.ServiceName,
		svcp.SVCPUsername,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)
	if booking.Cancel.CancelReason != "" {
		emailContent += fmt.Sprintf(`<br>
ด้วยสาเหตุ %s`, booking.Cancel.CancelReason)
	}
	emailContent += `<br><br>ขออภัยในความไม่สะดวก<br><br>Petpal Team</p>`

	err := mail.SendEmailWithGmail(user.Email, emailSubject, emailContent)
	if err != nil {
		return err
	}
	return nil
}

func NotifyCancelBookingToSVCP(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking cancelled"
	emailContent := fmt.Sprintf(`<h4>สวัสดีครับ</h4>
<p>การจองบริการ %s<br>
กับผู้ใช้บริการ %s<br>
เวลา %s ถึงเวลา %s<br>
ได้ถูกยกเลิกแล้ว`,
		booking.ServiceName,
		user.FullName,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)
	if booking.Cancel.CancelReason != "" {
		emailContent += fmt.Sprintf(`<br>
ด้วยสาเหตุ %s`, booking.Cancel.CancelReason)
	}
	emailContent += `<br><br>ขออภัยในความไม่สะดวก<br><br>Petpal Team</p>`

	err := mail.SendEmailWithGmail(svcp.SVCPEmail, emailSubject, emailContent)
	if err != nil {
		return err
	}

	return nil
}

// Complete booking
func NotifyCompleteBookingToUser(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking marked as completed"
	emailContent := fmt.Sprintf(`
<h4>สวัสดีครับ</h4>
<p>ผู้ให้บริการ %s<br>
ได้ยืนยันการจองบริการ %s<br>
เวลา %s ถึงเวลา %s<br>
ว่าสำเร็จแล้ว`,
		booking.ServiceName,
		svcp.SVCPUsername,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)
	emailContent += `<br><br>ขอขอบคุณที่ใช้บริการกับเรา<br><br>Petpal Team</p>`

	err := mail.SendEmailWithGmail(user.Email, emailSubject, emailContent)
	if err != nil {
		return err
	}
	return nil
}

func NotifyCompleteBookingToSVCP(booking *models.Booking, user *models.User, svcp *models.SVCP) error {
	emailSubject := "Petpal service booking marked as completed"
	emailContent := fmt.Sprintf(`<h4>สวัสดีครับ</h4>
<p>ผู้ใช้บริการ %s<br>
ได้ยืนยันการจองบริการ %s<br>
เวลา %s ถึงเวลา %s<br>
ว่าสำเร็จแล้ว`,
		user.FullName,
		booking.ServiceName,
		booking.StartTime.Format("15:04:05 02-01-2006 MST"),
		booking.EndTime.Format("15:04:05 02-01-2006 MST"),
	)
	emailContent += `<br><br>ขอขอบคุณที่ใช้บริการกับเรา<br><br>Petpal Team</p>`
	err := mail.SendEmailWithGmail(svcp.SVCPEmail, emailSubject, emailContent)
	if err != nil {
		return err
	}
	return nil
}

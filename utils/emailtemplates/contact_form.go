package emailtemplates

import (
	"fmt"
	"html"
	"strings"
)

func GenerateContactFormEmailBody(name, email, subject, message string) string {
	safeMessage := html.EscapeString(message)
	formattedMessage := strings.ReplaceAll(safeMessage, "\n", "<br>")

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>New Contact Form Submission</title>
        <style>
            body { 
				font-family: Arial, sans-serif; 
				line-height: 1.6; 
				margin: 0; 
				padding: 0; 
			}
            .container { 
				max-width: 500px; 
				margin: 20px auto; 
				background-color: #cecece; 
				padding: 20px; 
				border-radius: 8px; 
				border: 2px solid rgb(0, 0, 0);
			}
            .header { 
				padding: 10px 0; 
				text-align: center; 
			}
            .header h1 { 
				margin: 0; 
				font-size: 24px; 
			}
            .content { 
				padding: 20px 0; 
				font-size: 16px;
			}
            .detail-row { 
				margin-bottom: 10px;
			 }
            .detail-label { 
				font-weight: bold; 
				color: #555; 
			}
            .detail-value {
				margin-left: 5px; 
			}
            .message-box { 
				background-color: #f9f9f9; 
				border: 1px solid #333333; 
				padding: 20px; 
				margin-top: 20px; 
				border-radius: 5px; 
			}
            .footer { 
				text-align: center; 
				margin-top: 30px; 
				font-size: 12px; 
			}
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>New Contact Form Submission</h1>
            </div>
            <div class="content">
                <div class="detail-row">
                    <span class="detail-label">Name:</span>
                    <span class="detail-value">%s</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Email:</span>
                    <span class="detail-value">%s</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Subject:</span>
                    <span class="detail-value">%s</span>
                </div>
                <p class="detail-label">Message:</p>
                <div class="message-box">
                    <p>%s</p>
                </div>
            </div>
            <div class="footer">
                <p>This email was sent from the LMW Fitness contact form.</p>
            </div>
        </div>
    </body>
    </html>
    `, name, email, subject, formattedMessage)
}

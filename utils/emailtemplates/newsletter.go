package emailtemplates

import "fmt"

func GenerateNewsletterConfirmationEmailBody(recipientEmail, confirmLink string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>Confirm Your LMW Fitness Newsletter Subscription</title>
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
                border: 2px solid rgb(0, 0, 0);
                padding: 20px; 
                border-radius: 8px; 
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
            .button { 
                display: inline-block; 
                padding: 10px 20px; 
                margin: 20px 0; 
                background-color: #ffcf00; 
                color: #000000; 
                text-decoration: none !important; 
                border-radius: 8px; 
                font-weight: bold;
            }
            .footer { 
                text-align: center; 
                margin-top: 30px; 
                font-size: 12px; 
            }
            .footer a { 
                color: #ffcf00; 
                text-decoration: none; 
                font-size: 12px;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>
                    <span style="color: #21fc0d;">L</span>
                    <span style="color: #ffcf00;">M</span>
                    <span style="color: #ff11ff;">W</span>
                    <span> Fitness</span>
                </h1>
            </div>
            <div class="content">
                <p>Hello %s,</p>
                <p>Thank you for subscribing to the 
                <strong>
                    <span style="color: #21fc0d;">L</span>
                    <span style="color: #ffcf00;">M</span>
                    <span style="color: #ff11ff;">W</span>
                    <span> Fitness</span>
                </strong> newsletter!</p>
                <p>To complete your subscription, please click the button below:</p>
                <p style="text-align: center;">
                    <a href="%s" class="button">Confirm Subscription</a>
                </p>
                <p>If you did not subscribe to this newsletter, please ignore this email.</p>
                <p style="margin-bottom: 2px";>Many thanks,</p><br>
                <p><span style="color: #21fc0d;">L</span><span style="color: #ffcf00;">M</span><span style="color: #ff11ff;">W</span> Fitness</p>
            </div>
            <div class="footer">
                <p>&copy; 2025 LMW Fitness. All rights reserved.</p>
                <p><a href="https://www.lmwfitness.co.uk">Visit our website</a></p>
            </div>
        </div>
    </body>
    </html>
    `, recipientEmail, confirmLink)
}

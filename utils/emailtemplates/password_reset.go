package emailtemplates

import "fmt"

func GeneratePasswordResetEmailBody(recipientEmail, resetLink string) string {
	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>LMW Fitness - Password Reset</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                margin: 0;
                padding: 0;
            }
            .container {
                max-width: 500px;
                margin: 20px auto;
                padding: 20px;
				border: 2px solid rgb(0, 0, 0);
				border-radius: 8px;
				background-color: #cecece;
            }
            .header {
                text-align: center;
                padding: 10px 0;
            }
            .header h1 {
                font-size: 24px;
                margin: 0;
            }
            .content {
                padding: 20px 0;
                font-size: 16px;
            }
            .button {
                display: inline-block;
                padding: 12px 24px;
                margin: 20px 0;
                background-color: #ffcf00;
                text-decoration: none !important;
				color: #000000;
                border-radius: 8px;
                font-weight: bold;
            }
            .footer {
                text-align: center;
                font-size: 12px;
                margin-top: 30px;
            }
            .footer a {
                color: #ffcf00;
				font-size: 12px;
                text-decoration: none;
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
                <p>You have requested a password reset for your <strong>
                    <span style="color: #21fc0d;">L</span>
                    <span style="color: #ffcf00;">M</span>
                    <span style="color: #ff11ff;">W</span>
                    <span> Fitness</span>
                </strong> account.</p>
                <p style="text-align: center;">
                    <a href="%s" class="button" style="text-decoration: none !important">Please reset Your Password</a>
                </p>
                <p>This link will expire in 12 hours. If you did not request a password reset, please ignore this email.</p>
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
    `, recipientEmail, resetLink)
}

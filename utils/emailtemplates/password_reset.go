package emailtemplates

import "fmt"

func GeneratePasswordResetEmailBody(recipientEmail, resetLink string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
  <head>
    <meta charset="utf-8"> 
    <meta name="viewport" content="width=device-width"> 
    <meta http-equiv="x-ua-compatible" content="ie=edge"> 
    
    <title>LMW Fitness - Password Reset Request</title>

    <style>
      :root {
        --color-limeGreen: #21fc0d;
        --color-brightYellow: #ffcf00;
        --color-hotPink: #ff11ff;
        --color-customGray: #2a3241;
        --color-logoGray: #cecece;
        --color-customWhite: #f3f4f6;

        --font-higherJump: higherJump, sans-serif;
        --font-westerners: westerners, cursive;
        --font-titillium: titillium, sans-serif;
      }

      .preheader { display:none !important; visibility:hidden; opacity:0; color:transparent; height:0; width:0; overflow:hidden; mso-hide:all; }
      @media only screen and (max-width:600px){
        .container{ width:100%% !important; }
        .px-24{ padding-left:16px !important; padding-right:16px !important; }
        .py-24{ padding-top:16px !important; padding-bottom:16px !important; }
        .h1{ font-size:24px !important; line-height:32px !important; }
        .lead{ font-size:16px !important; line-height:24px !important; }
        .btn{ display:block !important; width:100%% !important; }
      }
      a { color:var(--color-customGray); }
      a.btn-link { color:var(--color-customGray) !important; text-decoration:none !important; }
      .reset-btn {
        display: inline-block;
        padding: 12px 24px;
        background-color: var(--color-brightYellow);
        color: var(--color-customGray) !important;
        text-decoration: none !important;
        border-radius: 8px;
        font-weight: bold;
        font-family: var(--font-titillium);
      }
    </style>
  </head>
  <body style="margin:0; padding:0; background-color:var(--color-customWhite);">

    <!-- Preheading -->
    <div class="preheader">Password Reset Request — We've received a request to reset your password. Click the link below to create a new password.</div>

    <center role="article" aria-roledescription="email" lang="en" style="width:100%%; background-color:var(--color-customWhite);">
      <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="100%%" style="background-color:var(--color-customWhite);">
        <tr>
          <td align="center">
            <table role="presentation" cellpadding="0" cellspacing="0" border="0" width="600" class="container" style="width:600px; max-width:600px;">
              <tr><td style="height:24px; line-height:24px;">&nbsp;</td></tr>

              <!-- Card -->
              <tr>
                <td class="px-24" style="padding:0 24px;">
                  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#ffffff; border-radius:12px; box-shadow:0 4px 14px rgba(0,0,0,0.06);">
                    <tr>
                      <td class="py-24 px-24" style="padding:28px;">

                        <!-- Main Heading -->
                        <h1 class="h1" style="margin:16px; padding-bottom:24px; font-family:var(--font-higherJump); font-size:26px; line-height:34px; color:var(--color-customGray);">
                          Password Reset
                        </h1>

                        <!-- Greeting -->
                        <p class="lead" style="margin:16px; font-family:var(--font-titillium); font-size:17px; line-height:26px; color:#444444;">
                         Hello %s,
                        </p>

                        <!-- Main message -->
                        <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
                         We received a request to reset your password for your LMW Fitness account.
                        </p>

                        <!-- Reset button -->
                        <div style="text-align: center; margin: 24px 0;">
                          <a href="%s" class="reset-btn" style="display: inline-block; padding: 12px 24px; background-color: var(--color-brightYellow); color: var(--color-customGray) !important; text-decoration: none !important; border-radius: 8px; font-weight: bold; font-family: var(--font-titillium);">
                            Reset Your Password
                          </a>
                        </div>

                        <!-- Line separator -->
                        <hr style="border:none; border-top:1px solid #efefef; margin:18px 0;">

                        <!-- Security notice -->
                        <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
                          <strong>Important:</strong> This link will expire in 12 hours for your security.
                        </p>
                        
                        <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
                          If you didn't request this password reset, please ignore this email. Your password will remain unchanged.
                        </p>

                        <!-- Support section -->
                        <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
                          Need help? You can 
                          <a href="https://wa.me/447516606668"
                            target="_blank"
                            rel="noopener noreferrer"
                            style="font-family:var(--font-titillium); color:var(--color-hotPink); font-weight:600; text-decoration:none;">
                            WhatsApp me
                          </a>
                          with any questions.
                        </p>

                        <!-- Ending -->
                        <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:var(--color-customGray);">
                            All the best,
                        </p>
                        <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:var(--color-customGray);">
                            Laura
                        </p>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>

              <!-- Footer -->
              <tr>
                <td class="px-24" align="center" style="padding:18px 24px 32px;">
                  <p style="margin:12px 0 6px; font-family:var(--font-titillium); font-size:12px; line-height:18px; color:var(--color-logoGray);">
                    © 2025 LMW Fitness • Live More With Fitness
                  </p>
                  <p style="margin:0; font-family:var(--font-titillium); font-size:12px; line-height:18px; color:var(--color-logoGray);">
                    You're receiving this because you requested a password reset.
                  </p>
                  <table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:16px 0 0;">
                    <tr>
                      <td style="padding:0 8px;">
                          <a href="https://www.facebook.com/profile.php?id=61573194721199" target="_blank">
                              <img src="https://api.lmwfitness.co.uk/images/LMW_fitness_facebook.png" alt="Facebook" width="30" height="30" style="display:block; border:0; width:30px; height:30px; margin:auto;" />
                          </a>
                      </td>
                      <td style="padding:0 8px;">
                        <a href="https://www.instagram.com/lmw__fitness/" target="_blank">
                          <img src="https://api.lmwfitness.co.uk/images/LMW_fitness_instagram.png" alt="Instagram" width="30" height="30" style="display:block; border:0; width:30px; height:30px; margin:auto;" />
                        </a>
                      </td>
                      <td style="padding:0 8px;">
                        <a href="https://www.tiktok.com/en/" target="_blank">
                          <img src="https://api.lmwfitness.co.uk/images/LMW_fitness_tiktok.png" alt="Tik Tok" width="30" height="30" style="display:block; border:0; width:30px; height:30px; margin:auto;" />
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              
            </table>
          </td>
        </tr>
      </table>
    </center>
  </body>
</html>
    `, recipientEmail, resetLink)
}

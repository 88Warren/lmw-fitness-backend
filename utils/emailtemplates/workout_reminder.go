package emailtemplates

import "fmt"

func GenerateWorkoutReminderEmailBody(recipientEmail string, daysSince int, currentStreak int, frontendURL string) string {
	var daysSinceText string
	if daysSince == 1 {
		daysSinceText = "yesterday"
	} else {
		daysSinceText = fmt.Sprintf("%d days ago", daysSince)
	}

	var streakSection string
	if currentStreak > 1 {
		streakSection = fmt.Sprintf(`
		<p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
			You had a <strong>%d day streak</strong> going — don't let it slip away! 🔥
		</p>`, currentStreak)
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <title>LMW Fitness - Time to Work Out!</title>
    <style>
      :root {
        --color-limeGreen: #21fc0d;
        --color-brightYellow: #ffcf00;
        --color-hotPink: #ff11ff;
        --color-customGray: #2a3241;
        --color-logoGray: #cecece;
        --color-customWhite: #f3f4f6;
        --font-titillium: titillium, sans-serif;
        --font-higherJump: higherJump, sans-serif;
      }
      .preheader { display:none !important; visibility:hidden; opacity:0; color:transparent; height:0; width:0; overflow:hidden; }
      @media only screen and (max-width:600px){
        .container{ width:100%% !important; }
      }
    </style>
  </head>
  <body style="margin:0; padding:0; background-color:#f3f4f6;">
    <div class="preheader">Your workout is waiting — it's been %s since your last session. Come back and keep the momentum going!</div>
    <center style="width:100%%; background-color:#f3f4f6;">
      <table cellpadding="0" cellspacing="0" border="0" width="100%%" style="background-color:#f3f4f6;">
        <tr><td align="center">
          <table cellpadding="0" cellspacing="0" border="0" width="600" class="container" style="width:600px; max-width:600px;">
            <tr><td style="height:24px;">&nbsp;</td></tr>
            <tr>
              <td style="padding:0 24px;">
                <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#ffffff; border-radius:12px; box-shadow:0 4px 14px rgba(0,0,0,0.06);">
                  <tr>
                    <td style="padding:28px;">
                      <h1 style="margin:16px; padding-bottom:8px; font-family:var(--font-higherJump); font-size:26px; color:var(--color-customGray);">
                        Your workout is waiting 💪
                      </h1>
                      <p style="margin:16px; font-family:var(--font-titillium); font-size:17px; line-height:26px; color:#444444;">
                        Hey %s,
                      </p>
                      <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
                        It's been <strong>%s</strong> since your last workout. Your programme is still here, ready when you are.
                      </p>
                      %s
                      <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:#444444;">
                        Even a short session counts. Consistency is what gets results — let's keep going!
                      </p>
                      <div style="text-align:center; margin:28px 0;">
                        <a href="%s/profile" style="display:inline-block; padding:14px 32px; background-color:#ffcf00; color:#2a3241; text-decoration:none; border-radius:8px; font-weight:bold; font-family:var(--font-titillium); font-size:16px;">
                          Back to My Workout →
                        </a>
                      </div>
                      <hr style="border:none; border-top:1px solid #efefef; margin:18px 0;">
                      <p style="margin:16px; font-family:var(--font-titillium); font-size:13px; line-height:20px; color:#888888;">
                        Don't want these reminders? 
                        <a href="%s/profile" style="color:#2a3241;">Update your preferences in your profile.</a>
                      </p>
                      <p style="margin:16px; font-family:var(--font-titillium); font-size:16px; line-height:24px; color:var(--color-customGray);">
                        All the best,<br>Laura
                      </p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>
            <tr>
              <td align="center" style="padding:18px 24px 32px;">
                <p style="margin:0; font-family:var(--font-titillium); font-size:12px; color:var(--color-logoGray);">
                  © 2025 LMW Fitness • Live More With Fitness
                </p>
              </td>
            </tr>
          </table>
        </td></tr>
      </table>
    </center>
  </body>
</html>
`, daysSinceText, recipientEmail, daysSinceText, streakSection, frontendURL, frontendURL)
}

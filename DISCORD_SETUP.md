# 🔔 Discord Webhook Setup Guide

This guide explains how to configure Discord webhooks for receiving bug reports and contact form submissions from MyLittlePrice.

## Overview

MyLittlePrice supports two types of Discord notifications:
- **Bug Reports** - Technical issues reported by users with screenshots and context
- **Contact Forms** - General inquiries and messages from the contact page

Each notification type can be sent to a **different Discord channel** using separate webhooks.

## Prerequisites

- Discord server with admin permissions
- Access to backend `.env` configuration

## Setup Steps

### 1. Create Discord Webhooks

#### For Bug Reports

1. Open your Discord server
2. Navigate to the channel where you want to receive bug reports (e.g., `#bug-reports`)
3. Click the **⚙️ gear icon** next to the channel name → **Integrations**
4. Click **Webhooks** → **New Webhook**
5. Configure the webhook:
   - **Name**: `MyLittlePrice Bug Reports`
   - **Channel**: Select your bug reports channel
   - **Avatar**: (Optional) Upload a bug icon
6. Click **Copy Webhook URL**
7. Save it for later

#### For Contact Forms

1. Navigate to a different channel (e.g., `#contact-messages`)
2. Follow the same steps as above
3. Configure the webhook:
   - **Name**: `MyLittlePrice Contact Form`
   - **Channel**: Select your contact messages channel
   - **Avatar**: (Optional) Upload a mail icon
4. Click **Copy Webhook URL**
5. Save it for later

### 2. Configure Backend Environment

Add the webhook URLs to your backend `.env` file:

```bash
# Bug Reports Webhook (for #bug-reports channel)
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/123456789012345678/AbCdEfGhIjKlMnOpQrStUvWxYz1234567890

# Contact Form Webhook (for #contact-messages channel)
CONTACT_WEBHOOK_URL=https://discord.com/api/webhooks/987654321098765432/ZyXwVuTsRqPoNmLkJiHgFeDcBa0987654321
```

**Important:**
- Leave empty to disable notifications for that type
- Both webhooks are optional
- You can use the same webhook for both if you want all notifications in one channel

### 3. Restart Backend

After updating the `.env` file, restart the backend:

```bash
cd backend
go run ./cmd/api/main.go
```

## Notification Examples

### Bug Report Notification

When a user submits a bug report, Discord receives:

```
🐛 New Bug Report

Description of the bug with user's detailed explanation...

👤 User
user@example.com
user-id-123

📍 Location
https://mylittleprice.com/chat

🔧 Session ID
session-abc-123

📝 Steps to Reproduce
1. Go to chat page
2. Click on product
3. Error appears

💻 Technical Details
Browser: Mozilla/5.0...
Screen: 1920x1080
Viewport: 1920x943

📎 Attachments: screenshot.png

Report ID: bug_report:20240101_120000:user-id
```

**Features:**
- Rich embed with color coding (red for bugs)
- User information and contact email
- Direct link to the problematic page
- Session ID for debugging
- Technical details (browser, screen resolution)
- Up to 3 screenshot attachments
- Unique report ID for tracking

### Contact Form Notification

When someone fills out the contact form, Discord receives:

```
📬 New Contact Form Submission

John Doe sent a message

👤 Name
John Doe

📧 Email
john@example.com

📋 Subject
Partnership Inquiry

💬 Message
Hi, I'm interested in partnering with MyLittlePrice for...

Contact ID: contact_form:20240101_120000:john_example_com
```

**Features:**
- Rich embed with color coding (blue for contact)
- Full name and email address
- Subject line
- Complete message content
- Unique contact ID for tracking
- Timestamp

## Testing

### Test Bug Report

1. Go to your frontend: `http://localhost:3000`
2. Click the **🐛 Bug Report** button (usually in the corner)
3. Fill out the form with test data
4. (Optional) Attach a screenshot
5. Submit
6. Check your Discord `#bug-reports` channel

### Test Contact Form

1. Go to: `http://localhost:3000/contact`
2. Fill out all fields:
   - Name
   - Email
   - Subject
   - Message
3. Click **Send Message**
4. Check your Discord `#contact-messages` channel

## Rate Limiting

To prevent spam, both endpoints have rate limits:

**Bug Reports:**
- Max 5 submissions per minute per IP address
- Message: "Too many bug reports, please try again later"

**Contact Forms:**
- Max 3 submissions per minute per IP address
- Message: "Too many contact form submissions, please try again later"

## Data Storage

Both submissions are stored in Redis for 90 days:

**Bug Reports:**
- Key format: `bug_report:YYYYMMDD_HHMMSS:user_id`
- Retention: 30 days
- Contains: Full bug report data + attachments (base64)

**Contact Forms:**
- Key format: `contact_form:YYYYMMDD_HHMMSS:email_sanitized`
- Retention: 90 days
- Contains: Name, email, subject, message, timestamp

## Troubleshooting

### Notifications Not Appearing

**Check webhook URLs:**
```bash
# Verify URLs are set correctly
echo $DISCORD_WEBHOOK_URL
echo $CONTACT_WEBHOOK_URL
```

**Check backend logs:**
```bash
# Look for Discord-related logs
✅ Discord notification sent successfully
❌ Failed to send Discord notification: ...
```

**Test webhook manually:**
```bash
curl -X POST "YOUR_WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{"content": "Test message"}'
```

### Webhook Returns Error

**Common issues:**

1. **Invalid webhook URL** - Verify the URL is correct and not expired
2. **Rate limit exceeded** - Discord has rate limits (30 requests/60 seconds per webhook)
3. **Webhook deleted** - Recreate the webhook in Discord
4. **Attachment too large** - Max 8MB per file, 25MB total per message

### Backend Errors

**"undefined: ContactWebhookURL"**
- Make sure you added `ContactWebhookURL string` to `config.go`
- Restart backend after config changes

**"Failed to marshal Discord payload"**
- Check message content for invalid characters
- Review backend logs for JSON errors

## Security Best Practices

1. ✅ **Never commit webhook URLs to git** - Use `.env` files only
2. ✅ **Regenerate webhooks if exposed** - Discord allows easy regeneration
3. ✅ **Use separate channels** - Keep bug reports and contact messages organized
4. ✅ **Monitor for abuse** - Check Discord regularly for spam
5. ✅ **Rate limiting is enabled** - Backend prevents spam submissions

## Optional: Webhook Customization

### Change Notification Colors

Edit the embed color in handlers:

**Bug Reports** (`backend/internal/handlers/bug_report.go:256`):
```go
"color": 15158332, // Red (default)
```

**Contact Forms** (`backend/internal/handlers/contact.go:129`):
```go
"color": 3447003, // Blue (default)
```

Color values use decimal format. [Color picker](https://www.spycolor.com/)

### Add More Fields

You can customize what information is sent by editing the `fields` array in:
- `backend/internal/handlers/bug_report.go` (line 213-250)
- `backend/internal/handlers/contact.go` (line 106-119)

## Advanced: Discord Bot Integration

For more advanced features (like responding to reports via Discord), consider:
- Creating a Discord bot instead of webhooks
- Implementing two-way communication
- Adding reaction-based triage

See [Discord Developer Portal](https://discord.com/developers/docs) for bot documentation.

## Support

If you encounter issues:
1. Check backend logs for detailed error messages
2. Verify webhook URLs are valid
3. Test with manual curl request
4. Check Discord server permissions
5. Ensure backend can reach Discord API (no firewall blocking)

---

**Last Updated**: 2024
**Maintained By**: MyLittlePrice Team

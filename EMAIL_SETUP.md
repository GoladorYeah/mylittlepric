# 📧 Email Configuration for Password Reset

This guide explains how to configure email sending for password reset functionality using Zoho Mail SMTP.

## Overview

The password reset feature requires SMTP configuration to send reset links to users. We're using **Zoho Mail** as our email provider with the following configuration:

- **SMTP Server**: smtp.zoho.eu
- **Port**: 587 (TLS) or 465 (SSL)
- **Email Address**: info@mylittleprice.com

## Configuration Steps

### 1. Update `.env` File

Add the following environment variables to your backend `.env` file:

```bash
# Email Configuration (SMTP)
SMTP_HOST=smtp.zoho.eu
SMTP_PORT=587
SMTP_USERNAME=info@mylittleprice.com
SMTP_PASSWORD=your-zoho-email-password-here
SMTP_FROM_EMAIL=info@mylittleprice.com
SMTP_FROM_NAME=MyLittlePrice

# Frontend URL for reset links
# For development use: http://localhost:3000
# For production use: https://mylittleprice.com
FRONTEND_URL=https://mylittleprice.com
```

### 2. Get Zoho Mail App Password

For security, it's recommended to use an **App-Specific Password** instead of your main email password:

1. Log in to [Zoho Mail](https://mail.zoho.eu)
2. Go to **Settings** → **Security** → **App Passwords**
3. Click **Generate New Password**
4. Give it a name (e.g., "MyLittlePrice Backend")
5. Copy the generated password
6. Use this password in `SMTP_PASSWORD` environment variable

### 3. DNS Configuration (Already Done)

Your domain DNS is already configured with the following records:

**MX Records:**
```
Priority 10: mx.zoho.eu
Priority 20: mx2.zoho.eu
Priority 50: mx3.zoho.eu
```

**SPF Record:**
```
v=spf1 include:zohomail.eu ~all
```

### 4. Test Email Sending

Restart your backend server and test the password reset flow:

```bash
cd backend
go run ./cmd/api/main.go
```

Then:
1. Go to http://localhost:3000/login
2. Click "Forgot password?"
3. Enter an email address (must be a registered user with `provider='email'`)
4. Check the email inbox for the reset link

## Email Template

The password reset email includes:

- Professional HTML template with MyLittlePrice branding
- Clickable "Reset Password" button
- Reset link as plain text (fallback)
- 1-hour expiration warning
- Security tips

## Troubleshooting

### Email Not Sending

**Check logs for errors:**
```bash
# Look for email-related logs in backend console
⚠️ Failed to send password reset email: ...
```

**Common issues:**

1. **Invalid credentials**: Verify `SMTP_USERNAME` and `SMTP_PASSWORD`
2. **Port blocked**: Try port 465 (SSL) instead of 587 (TLS)
3. **Firewall**: Ensure outbound SMTP connections are allowed

### Testing Without Email

If email is not configured, the backend will:
- Generate the token successfully
- Return it in the API response (fallback mode)
- Log a warning: `⚠️ Failed to send password reset email`

You can still test password reset by:
1. Using the token from the API response
2. Manually navigating to: `http://localhost:3000/reset-password?token=YOUR_TOKEN`

### Check SMTP Configuration

```bash
# Test SMTP connection with openssl
openssl s_client -connect smtp.zoho.eu:587 -starttls smtp
```

## Production Configuration

For production deployment:

### Update Environment Variables

```bash
# Production SMTP (same server)
SMTP_HOST=smtp.zoho.eu
SMTP_PORT=587
SMTP_USERNAME=info@mylittleprice.com
SMTP_PASSWORD=your-production-app-password

# Production Frontend URL
FRONTEND_URL=https://mylittleprice.com
```

### Security Best Practices

1. ✅ Use App-Specific Password (not main password)
2. ✅ Store credentials in environment variables (never commit to git)
3. ✅ Use TLS/SSL for SMTP connection
4. ✅ Monitor failed email attempts
5. ✅ Rate limit password reset requests (already implemented)

### Remove Testing Fallback

In production, remove the fallback that returns the token in the API response:

**File**: `backend/internal/handlers/auth.go`

```go
// Remove this fallback in production:
if err := h.container.EmailService.SendPasswordResetEmail(req.Email, resetToken); err != nil {
    // Don't return the token - just log the error
    fmt.Printf("⚠️ Failed to send password reset email: %v\n", err)

    // Return generic success message
    return c.JSON(fiber.Map{
        "message": "If an account exists with this email, a password reset link has been sent",
    })
}
```

## Email Delivery Tips

### Avoid Spam Folder

1. **SPF Record**: ✅ Already configured (`v=spf1 include:zohomail.eu ~all`)
2. **DKIM**: Configure in Zoho Mail settings
3. **DMARC**: Add DMARC record for domain
4. **Reverse DNS**: Ensure server IP has proper PTR record

### Monitor Email Deliverability

- Check Zoho Mail logs for bounce/rejection rates
- Monitor user complaints about missing emails
- Set up alerts for SMTP errors

## API Endpoints

### Request Password Reset
```http
POST /api/auth/request-password-reset
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Response (Email Sent):**
```json
{
  "message": "If an account exists with this email, a password reset link has been sent"
}
```

**Response (Email Failed - Development Only):**
```json
{
  "message": "Password reset token generated (email failed, using fallback)",
  "token": "eyJhbGc..."
}
```

### Reset Password
```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "eyJhbGc...",
  "new_password": "newSecurePassword123"
}
```

**Response:**
```json
{
  "message": "Password reset successfully"
}
```

## Frontend Routes

- **Login Page**: `/login` - Has "Forgot password?" link
- **Reset Password Page**: `/reset-password?token=...` - Enter new password
- **Settings Page**: `/settings` - Change password (for authenticated users)

## Rate Limiting

Password reset requests are rate-limited to prevent abuse:

- **Endpoint**: `/api/auth/request-password-reset`
- **Limit**: Uses auth rate limiter (configurable in `.env`)
- **Default**: 5 requests per minute per IP

## Token Security

- **Storage**: Redis with 1-hour expiration
- **Hashing**: SHA-256 before storage
- **One-time use**: Marked as used after successful reset
- **Auto-cleanup**: Expired tokens removed automatically

## Support

If you encounter issues:

1. Check backend logs for detailed error messages
2. Verify SMTP credentials are correct
3. Test SMTP connection manually with openssl
4. Contact Zoho support if email delivery fails

---

**Last Updated**: 2024
**Maintained By**: MyLittlePrice Team

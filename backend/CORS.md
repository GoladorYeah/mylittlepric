# CORS Configuration Guide

## Overview

Cross-Origin Resource Sharing (CORS) is a security feature that controls which domains can access your API. This guide helps you configure CORS correctly for MyLittlePrice.

## 🚨 Quick Fix for Production Errors

If you're seeing this error in production:
```
Access to fetch at 'https://api.mylittleprice.com/api/product-details' from origin 'https://mylittleprice.com'
has been blocked by CORS policy: Response to preflight request doesn't pass access control check:
No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

**Root Cause:** Your production backend server doesn't have the correct `CORS_ORIGINS` environment variable configured.

**Solution:** Set the `CORS_ORIGINS` environment variable on your production server:

```bash
CORS_ORIGINS=https://mylittleprice.com,https://www.mylittleprice.com
```

### How to Fix (Deployment-Specific)

#### Docker Deployment
1. Add to your production `.env` file:
   ```bash
   CORS_ORIGINS=https://mylittleprice.com,https://www.mylittleprice.com
   ```

2. Restart your backend container:
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml restart backend
   ```

#### Environment Variable Deployment (Vercel, Heroku, AWS, etc.)
1. Go to your hosting platform's environment variables section
2. Add or update: `CORS_ORIGINS=https://mylittleprice.com,https://www.mylittleprice.com`
3. Redeploy or restart your backend service

#### Verify the Fix
After updating, test with curl:
```bash
curl -X OPTIONS https://api.mylittleprice.com/api/product-details \
  -H "Origin: https://mylittleprice.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v
```

Expected response should include:
```
Access-Control-Allow-Origin: https://mylittleprice.com
```

## Configuration

### Environment Variable

The `CORS_ORIGINS` environment variable accepts a comma-separated list of allowed origins:

```bash
# Development
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# Production (example)
CORS_ORIGINS=https://mylittleprice.com,https://www.mylittleprice.com

# Multiple domains
CORS_ORIGINS=https://app.example.com,https://admin.example.com,https://example.com
```

### Docker Deployment

For Docker deployments, ensure your `.env` file contains the correct `CORS_ORIGINS`:

```bash
# In your .env file
CORS_ORIGINS=https://yourdomain.com
```

Then deploy using:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Important Rules

1. **NO trailing slashes**
   - ✅ Correct: `https://mylittleprice.com`
   - ❌ Wrong: `https://mylittleprice.com/`

2. **Include ALL frontend domains**
   - If users access your site via multiple URLs, include all of them
   - Common pattern: include both `www` and non-`www` versions

3. **Use HTTPS in production**
   - ✅ Production: `https://mylittleprice.com`
   - ✅ Development: `http://localhost:3000`

4. **DO NOT include the API domain**
   - Only include domains where your **frontend** is hosted
   - The backend API domain should NOT be in CORS_ORIGINS

## Google OAuth Configuration

For Google OAuth to work, you also need to configure:

### 1. Google Cloud Console Settings

Go to [Google Cloud Console](https://console.cloud.google.com/apis/credentials):

1. Select your OAuth 2.0 Client ID
2. Add **Authorized JavaScript origins**:
   ```
   https://mylittleprice.com
   https://www.mylittleprice.com
   ```

3. Add **Authorized redirect URIs**:
   ```
   https://mylittleprice.com/auth/callback
   https://www.mylittleprice.com/auth/callback
   ```

### 2. Backend Environment Variables

```bash
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=https://mylittleprice.com/auth/callback
```

### 3. Frontend Environment Variables

```bash
NEXT_PUBLIC_API_URL=https://api.mylittleprice.com
```

## Troubleshooting

### Error: "No 'Access-Control-Allow-Origin' header"

**Cause:** The origin making the request is not in `CORS_ORIGINS`

**Solution:**
1. Check what domain the error message shows (e.g., `https://mylittleprice.com`)
2. Add that exact domain to `CORS_ORIGINS`
3. Restart your backend service

### Error: "The 'Access-Control-Allow-Origin' header contains multiple values"

**Cause:** Your reverse proxy (nginx, Caddy, etc.) is adding CORS headers AND the backend is adding them

**Solution:** Remove CORS headers from your reverse proxy config and let the backend handle it

### Google Sign-In Button Doesn't Work

**Possible causes:**
1. CORS not configured correctly (see above)
2. Google Client ID not set in environment variables
3. Domain not authorized in Google Cloud Console
4. Frontend using wrong API URL

**Debug steps:**
1. Open browser DevTools (F12) → Console tab
2. Look for errors related to CORS or Google OAuth
3. Check Network tab for failed requests
4. Verify environment variables are set correctly

### Testing CORS Locally

Test if CORS is working:

```bash
curl -X OPTIONS https://api.mylittleprice.com/api/auth/google \
  -H "Origin: https://mylittleprice.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v
```

Expected response should include:
```
Access-Control-Allow-Origin: https://mylittleprice.com
Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
Access-Control-Allow-Headers: Origin,Content-Type,Accept,Authorization
Access-Control-Allow-Credentials: true
```

## Production Deployment Checklist

Before deploying to production, verify:

- [ ] `CORS_ORIGINS` includes your production domain(s)
- [ ] `GOOGLE_CLIENT_ID` is set
- [ ] `GOOGLE_CLIENT_SECRET` is set
- [ ] `GOOGLE_REDIRECT_URL` points to your production domain
- [ ] Google Cloud Console has your domain in authorized origins
- [ ] Google Cloud Console has your callback URL in authorized redirects
- [ ] Frontend `NEXT_PUBLIC_API_URL` points to your API domain
- [ ] All domains use HTTPS (not HTTP)
- [ ] JWT secrets are set to secure random values

## Code Reference

CORS configuration is located in:

- **Configuration loading**: `backend/internal/config/config.go:219`
- **Middleware setup**: `backend/cmd/api/main.go:47-52`
- **Google auth endpoint**: `backend/internal/handlers/auth.go:68-94`

## Support

If you continue to have CORS issues after following this guide:

1. Check the backend logs for CORS-related messages
2. Verify your environment variables are loaded correctly
3. Ensure your reverse proxy isn't interfering with CORS headers
4. Test with the curl command above to isolate the issue

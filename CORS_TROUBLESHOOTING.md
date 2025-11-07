# CORS Troubleshooting Guide - Production

## 🚨 Problem: Product Details API Blocked by CORS

**Error Message:**
```
Access to fetch at 'https://api.mylittleprice.com/api/product-details' from origin 'https://mylittleprice.com'
has been blocked by CORS policy: Response to preflight request doesn't pass access control check:
No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

## ✅ Solution

Your production backend server is **missing the CORS configuration** for your frontend domain.

### Step 1: Set the Environment Variable

Add this to your production environment:

```bash
CORS_ORIGINS=https://mylittleprice.com,https://www.mylittleprice.com
```

**Important Notes:**
- ✅ Use HTTPS (not HTTP) for production
- ✅ Include both `www` and non-`www` versions if users can access both
- ✅ NO trailing slashes (`https://mylittleprice.com` not `https://mylittleprice.com/`)
- ✅ Comma-separated for multiple domains
- ❌ DO NOT include the API domain (`api.mylittleprice.com`)

### Step 2: Restart Your Backend Service

**Docker:**
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml restart backend
```

**Other Platforms:**
- Vercel: Redeploy after adding environment variable
- Heroku: `heroku restart`
- AWS: Restart your service/container
- Railway: Automatic restart after env var change

### Step 3: Verify the Fix

Test with curl:
```bash
curl -X OPTIONS https://api.mylittleprice.com/api/product-details \
  -H "Origin: https://mylittleprice.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -i
```

**Expected Output:**
```
HTTP/1.1 204 No Content
Access-Control-Allow-Origin: https://mylittleprice.com
Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS
Access-Control-Allow-Headers: Origin,Content-Type,Accept,Authorization
Access-Control-Allow-Credentials: true
```

### Step 4: Check Backend Logs

After the fix, check your backend logs. You should see:
```
🌍 Allowed origins: [https://mylittleprice.com https://www.mylittleprice.com]
```

If CORS is still failing, the logs will show:
```
❌ CORS REJECTED - Origin: 'https://mylittleprice.com' not in allowed list: [http://localhost:3000]
💡 Fix: Add 'https://mylittleprice.com' to CORS_ORIGINS environment variable
```

## 🔍 Common Issues

### Issue: Still Getting CORS Error After Setting Environment Variable

**Possible Causes:**
1. Backend service wasn't restarted
2. Environment variable syntax error (extra spaces, trailing slashes)
3. Wrong environment file loaded
4. Reverse proxy (nginx/Caddy) interfering with CORS headers

**Debug Steps:**
1. Check backend logs to verify the configured origins
2. Restart the backend service completely (not just reload)
3. Verify no reverse proxy is adding conflicting CORS headers
4. Test directly against the backend (bypass proxy if any)

### Issue: CORS Works for Some Endpoints but Not `/api/product-details`

**Cause:** This is unlikely with the current codebase - CORS is configured globally.

**Check:**
1. Verify the endpoint path is correct: `/api/product-details` (POST)
2. Check if there's a specific middleware on that route
3. Review backend logs for that specific request

### Issue: Multiple CORS Headers in Response

**Cause:** Both your reverse proxy AND the backend are adding CORS headers.

**Solution:** Remove CORS configuration from your reverse proxy (nginx, Caddy, etc.) and let the backend handle it.

## 📚 Additional Resources

- Full CORS guide: [backend/CORS.md](backend/CORS.md)
- Environment variables: [backend/.env.example](backend/.env.example)
- Production deployment: [DOCKER.md](DOCKER.md)

## 🆘 Need Help?

1. Check backend logs for CORS rejection messages
2. Verify your `.env` file or environment variables
3. Test with the curl command above
4. Review [backend/CORS.md](backend/CORS.md) for detailed configuration

## Quick Reference

| Environment | CORS_ORIGINS Example |
|-------------|---------------------|
| Development | `http://localhost:3000,http://localhost:3001` |
| Staging | `https://staging.mylittleprice.com` |
| Production | `https://mylittleprice.com,https://www.mylittleprice.com` |
| Multiple Prod Domains | `https://app.example.com,https://admin.example.com,https://example.com` |

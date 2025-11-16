# Bug Report System

Custom bug reporting solution for MyLittlePrice.

## Features

- **Easy to use**: Simple button + modal dialog
- **File attachments**: Upload screenshots via drag-and-drop, file picker, or Ctrl+V paste
- **Automatic context collection**: Captures user info, session, browser details, last messages
- **Rate limited**: 5 reports per minute to prevent spam
- **Stored in Redis**: 30-day retention for quick access
- **Logged to slog**: Can be picked up by external monitoring (Loki, etc.)
- **Discord notifications**: Optional webhook integration with screenshot attachments

## Components

### Frontend

- **BugReportButton** - Button component with 2 variants:
  - `full`: Full button with icon + text (for expanded sidebar)
  - `icon`: Icon-only with tooltip (for collapsed sidebar)

- **BugReportDialog** - Modal form for bug report submission

### Backend

- **POST /api/bug-report** - Submit bug report
- **GET /api/bug-reports** - List all reports (admin only, commented out)

## Auto-Collected Context

The system automatically collects:

- User ID & email (if logged in)
- Session ID
- Current URL
- User agent (browser info)
- Screen resolution & viewport size
- Last 5 chat messages (truncated to 100 chars each)
- Timestamp
- Screenshots (optional, up to 3 files, max 5MB each)

## Usage

### In Sidebar (Already Added)

```tsx
import { BugReportButton } from "@/features/bug-report";

// Expanded sidebar
<BugReportButton variant="full" />

// Collapsed sidebar (icon only)
<BugReportButton variant="icon" />
```

### In Chat Header

```tsx
import { BugReportButton } from "@/features/bug-report";

// Header with responsive text
<BugReportButton variant="header" />
```

### In Other Components

```tsx
import { BugReportButton } from "@/features/bug-report";

<BugReportButton variant="full" className="custom-class" />
```

## Attaching Screenshots

Users can attach screenshots in three ways:

1. **Drag and drop**: Drag image files into the upload area
2. **File picker**: Click "Click to upload" to select files
3. **Paste from clipboard**: Press `Ctrl+V` (or `Cmd+V` on Mac) to paste screenshots

**Limitations:**
- Max 3 files per report
- Max 5MB per file
- Only image files allowed (PNG, JPG, GIF, etc.)

## Data Storage

Bug reports are stored in Redis with key format:
```
bug_report:YYYYMMDD_HHMMSS:user_id
```

TTL: 30 days

## Logs

All bug reports are logged with slog:
```go
slog.Info("Bug report submitted",
    "user_id", req.Context.UserID,
    "description", req.Description,
    // ... more fields
)
```

## Discord Notifications

Bug reports are automatically sent to Discord if `DISCORD_WEBHOOK_URL` is configured.

### Setup Discord Webhook:

1. Go to your Discord server
2. Server Settings → Integrations → Webhooks
3. Click "New Webhook"
4. Choose channel (e.g., #bug-reports)
5. Copy webhook URL
6. Add to `.env`:
```env
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN
```

### Discord Embed Format

```text
🐛 New Bug Report
[Description]

👤 User: user@email.com
        user_id_123

📍 Location: https://app.com/chat

🔧 Session ID: session_123

📝 Steps to Reproduce:
[Steps if provided]

💻 Technical Details:
Browser: Mozilla/5.0...
Screen: 1920x1080
Viewport: 1366x768

Report ID: bug_report:20250116_143052:user_id

[Screenshot attachments displayed inline]
```

**Note:** Screenshots are automatically attached and displayed in Discord messages (up to 3 images).

## Future Enhancements

- [x] Discord webhook notifications ✅
- [x] Screenshot upload support ✅
- [ ] Email notifications
- [ ] Admin dashboard to view reports
- [ ] Persistent database storage (current: Redis only)
- [ ] Priority/severity selection
- [ ] Bug status tracking (open/in-progress/resolved)
- [ ] Slack webhook integration
- [ ] Video recording support
- [ ] Console logs capture

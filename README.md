# Mail Cleaner

Automatic email cleanup utility using IMAP protocol with configurable deletion rules.

## 📋 Features

- Connect to any IMAP server
- Delete emails by rules:
  - By sender address
  - By sender domain
- Optimized for large mailboxes (tens of thousands of emails)
- Multiple email accounts support via configuration

---

## 🚀 Quick Start

### Installation

```bash
git clone <repository-url>
cd mail-cleaner
go mod download
```

### Configuration

1. Create `.env.<service-name>` file (e.g., `.env.ukrnet`):

```env
IMAP_SERVER=imap.ukr.net
IMAP_PORT=993
EMAIL=your-email@ukr.net
PASSWORD=your-app-password
```

⚠️ **Important:** Use app-specific passwords, not your main email password!

2. Create rules file (e.g., `rules.json`):

```json
[
  {
    "type": "address_rule",
    "address": "spam@example.com"
  },
  {
    "type": "domain_rule",
    "domain": "newsletter.com"
  }
]
```

### Usage

```bash
go run . <service-name> <rules-file>
```

**Example:**

```bash
# Clean ukr.net mailbox with rules from rules.json
go run . ukrnet rules.json
```

### Build

```bash
go build -o mail-cleaner
./mail-cleaner ukrnet rules.json
```

---

## 📝 Deletion Rules

### Rule Types

#### 1. Address Rule - delete by exact email address
```json
{
  "type": "address_rule",
  "address": "noreply@spam.com"
}
```
Deletes all emails from `noreply@spam.com`.

#### 2. Domain Rule - delete by domain
```json
{
  "type": "domain_rule",
  "domain": "marketing.com"
}
```
Deletes all emails from domains containing `marketing.com` (e.g., `news@marketing.com`, `promo@marketing.com`).

### Full Rules File Example

```json
[
  {
    "type": "address_rule",
    "address": "no-reply@mail.instagram.com"
  },
  {
    "type": "address_rule",
    "address": "notification@facebookmail.com"
  },
  {
    "type": "domain_rule",
    "domain": "newsletter.com"
  },
  {
    "type": "domain_rule",
    "domain": "marketing"
  }
]
```

---

## ⚙️ Email Service Configuration

### Ukr.net
Create `.env.ukrnet`:
```env
IMAP_SERVER=imap.ukr.net
IMAP_PORT=993
EMAIL=your-email@ukr.net
PASSWORD=your-password
```

⚠️ **Note:** This tool was tested with Ukr.net. For other IMAP servers, create appropriate `.env.<service>` file with your server settings.

---

## 🔍 Usage Examples

### Delete all newsletters
```json
[
  {"type": "domain_rule", "domain": "newsletter"},
  {"type": "domain_rule", "domain": "mailing"},
  {"type": "domain_rule", "domain": "noreply"}
]
```

### Delete social networks
```json
[
  {"type": "domain_rule", "domain": "facebook"},
  {"type": "domain_rule", "domain": "instagram"},
  {"type": "domain_rule", "domain": "twitter"}
]
```

### Delete specific addresses
```json
[
  {"type": "address_rule", "address": "spam@example.com"},
  {"type": "address_rule", "address": "ads@company.com"}
]
```

---

## ⚠️ Warnings

- ⚠️ **Emails are deleted permanently!** Test rules on a test account first
- ⚠️ Make sure rules won't delete important emails
- ⚠️ Use app-specific passwords, not main passwords
- ⚠️ For large mailboxes (>10000 emails) processing may take time

---

## 🐛 Troubleshooting

### Connection Error
```
Error connecting to IMAP server: dial tcp: lookup failed
```
**Solution:** Check IMAP_SERVER and IMAP_PORT in `.env.<service-name>`

### Authentication Error
```
Error connecting to IMAP server: LOGIN failed
### Authentication Error
```
Error connecting to IMAP server: LOGIN failed
```
**Solution:** 
- Check EMAIL and PASSWORD in `.env.<service-name>`
- Enable IMAP in email settings
No valid rules found in the rules file.
```
**Solution:** Check JSON syntax in rules file

---

## 🛠️ Adding Custom Rules

### Step 1: Create rule file

Create `internal/rules/rule/your_rule.go`:

```go
package rule

import (
    "fmt"
    "strings"
    "mail-cleaner/internal/rules"
    "github.com/emersion/go-imap"
)

type SubjectRule struct {
    Keyword string
}

func init() {
    RegisterRuleFactory("subject_rule", func(data map[string]any) (rules.Rule, error) {
        keyword, ok := data["keyword"].(string)
        if !ok {
            return nil, fmt.Errorf("invalid 'keyword' field")
        }
        return NewSubjectRule(keyword)
    })
}

func NewSubjectRule(keyword string) (*SubjectRule, error) {
    if keyword == "" {
        return nil, fmt.Errorf("keyword cannot be empty")
    }
    return &SubjectRule{Keyword: keyword}, nil
}

func (r *SubjectRule) ShouldDelete(msg *imap.Message) bool {
    if msg.Envelope == nil {
        return false
    }
    return strings.Contains(strings.ToLower(msg.Envelope.Subject), 
                          strings.ToLower(r.Keyword))
}
```

### Step 2: That's it!

The factory will automatically register your rule via `init()`.

### Step 3: Use in JSON

```json
{
  "type": "subject_rule",
  "keyword": "unsubscribe"
}
```

**Available fields in `msg.Envelope`:**
- `Subject` - email subject
- `From` - sender info (slice of `*imap.Address`)
- `To` - recipients
- `Date` - email date
- `ReplyTo`, `Cc`, `Bcc` - other headers

**Example - delete by subject keyword:**
```go
func (r *SubjectRule) ShouldDelete(msg *imap.Message) bool {
    if msg.Envelope == nil || msg.Envelope.Subject == "" {
        return false
    }
    return strings.Contains(strings.ToLower(msg.Envelope.Subject), 
                          strings.ToLower(r.Keyword))
}
```

---

## 📄 License

MIT

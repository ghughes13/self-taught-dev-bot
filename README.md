# Self-Taught Dev Bot

A Discord bot written in Go that allows users to assign and remove roles themselves using simple commands.

## Features

- **Channel-specific commands**: Bot only responds in the `#set-role` channel (ID: 614590342724845580)
- **Self-service role assignment**: Users can assign themselves developer roles using `.iam <role>`
- **Role removal**: Users can remove roles using `.iamnot <role>`
- **Short name mapping**: Short names like "Frontend" map to full role names like "Frontend Developer"
- **Case-insensitive**: Role names are matched case-insensitively for convenience
- **Pre-configured roles**: Supports Frontend, Backend, Fullstack, and Mobile developers

## Prerequisites

- Go 1.21 or higher
- A Discord bot token (see setup instructions below)

## Setup

### 1. Create a Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to the "Bot" section in the left sidebar
4. Click "Add Bot" and confirm
5. Under "Token", click "Reset Token" or "Copy" to get your bot token
6. **IMPORTANT**: Enable the following Privileged Gateway Intents:
   - **MESSAGE CONTENT INTENT** (Required for the bot to read message content)
7. Under "OAuth2" > "URL Generator":
   - Select scopes: `bot`
   - Select bot permissions: `Manage Roles`
   - Copy the generated URL and open it in your browser to invite the bot to your server

### 2. Configure Role Hierarchy

**Important**: In your Discord server:
1. Go to Server Settings > Roles
2. Make sure the bot's role is placed **above** the roles you want users to be able to assign themselves
3. The bot needs to have "Manage Roles" permission
4. The roles users can assign should have "Display role members separately" if desired

### 3. Set Environment Variable

Set the `DISCORD_BOT_TOKEN` environment variable with your bot token:

**Windows (PowerShell):**
```powershell
$env:DISCORD_BOT_TOKEN="your-bot-token-here"
```

**Windows (Command Prompt):**
```cmd
set DISCORD_BOT_TOKEN=your-bot-token-here
```

**Linux/Mac:**
```bash
export DISCORD_BOT_TOKEN="your-bot-token-here"
```

Alternatively, you can create a `.env` file (make sure it's in `.gitignore`):
```
DISCORD_BOT_TOKEN=your-bot-token-here
```

And load it before running (you may need a package like `godotenv`).

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Bot

```bash
go run main.go
```

Or build and run:

```bash
go build -o bot.exe  # Windows
# or
go build -o bot      # Linux/Mac

./bot.exe  # Windows
# or
./bot      # Linux/Mac
```

## Usage

Once the bot is running, users can use these commands in the `#set-role` channel:

- `.iam <role>` - Assign yourself a developer role (e.g., `.iam Frontend`)
- `.iamnot <role>` - Remove a developer role from yourself (e.g., `.iamnot Frontend`)
- `.help` - Show help information

### Available Roles

The bot supports the following roles (short names map to full role names):

- `Frontend` → "Frontend Developer"
- `Backend` → "Backend Developer"
- `Fullstack` → "Fullstack Developer"
- `Mobile` → "Mobile Developer"

### Example Commands

```
.iam Frontend      # Assigns "Frontend Developer" role
.iam Backend       # Assigns "Backend Developer" role
.iam Fullstack     # Assigns "Fullstack Developer" role
.iam Mobile        # Assigns "Mobile Developer" role
.iamnot Frontend   # Removes "Frontend Developer" role
.help              # Shows help information
```

**Note**: Commands are case-insensitive, so `.iam frontend` works the same as `.iam Frontend`.

## Troubleshooting

### Bot doesn't assign roles

1. **Check bot permissions**: The bot needs "Manage Roles" permission
2. **Check role hierarchy**: The bot's role must be higher than the roles it's trying to assign
3. **Check bot role**: The bot needs a role in the server (you can give it one during the invite process)

### Bot doesn't respond to commands

1. **Check channel**: The bot only responds in the `#set-role` channel (ID: 614590342724845580)
2. **Check Message Content Intent**: Make sure "MESSAGE CONTENT INTENT" is enabled in the Discord Developer Portal
3. **Check bot is online**: The bot should show as "Online" in Discord
4. **Check command prefix**: Commands must start with `.` (not `!`)

### "Missing Permissions" error

- Make sure the bot's role is placed above the target role in Server Settings > Roles
- Ensure the bot has "Manage Roles" permission

## Security Notes

- **Never commit your bot token** to version control
- The `.gitignore` file is set up to exclude `.env` files
- Keep your bot token secret and rotate it if it's ever exposed

## License

MIT License - see LICENSE file for details

## Contributing

Feel free to submit issues and pull requests!

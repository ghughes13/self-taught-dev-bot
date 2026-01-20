package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	botToken         string
	botID            string
	setRoleChannelID string
	// Role mapping from short names to full role names
	roleMapping = map[string]string{
		"frontend":  "Frontend Developer",
		"backend":   "Backend Developer",
		"fullstack": "Fullstack Developer",
		"mobile":    "Mobile Developer",
		"student":   "Student",
	}
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Get bot token from environment variable
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	setRoleChannelID = os.Getenv("SET_ROLE_CHANNEL_ID")
	if botToken == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is not set")
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
	}

	// Get bot user info
	user, err := dg.User("@me")
	if err != nil {
		log.Fatal("Error obtaining account details: ", err)
	}
	botID = user.ID
	fmt.Printf("Bot is running as: %s#%s\n", user.Username, user.Discriminator)

	// Register message handler
	dg.AddHandler(messageCreate)

	// Register ready handler
	dg.AddHandler(ready)

	// Open websocket connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	<-make(chan struct{})
}

// ready is called when the bot is ready
func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}

// messageCreate handles incoming messages
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == botID {
		return
	}

	// Only respond in the set-role channel
	if m.ChannelID != setRoleChannelID {
		return
	}

	// Check if message is a command (starts with .)
	if !strings.HasPrefix(m.Content, ".") {
		return
	}

	// Parse command
	parts := strings.Fields(m.Content)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0][1:]) // Remove the . prefix

	switch command {
	case "iam":
		handleIamCommand(s, m, parts[1:])
	case "iamnot":
		handleIamNotCommand(s, m, parts[1:])
	case "help":
		handleHelpCommand(s, m)
	default:
		// Unknown command
	}
}

// handleIamCommand handles the .iam command
func handleIamCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `.iam <role>` - Available roles: Frontend, Backend, Fullstack, Mobile, Student")
		return
	}

	// Get the short role name (case-insensitive)
	shortName := strings.ToLower(strings.Join(args, " "))

	// Look up the full role name
	fullRoleName, exists := roleMapping[shortName]
	if !exists {
		// Build available roles message
		var availableRoles []string
		for key := range roleMapping {
			availableRoles = append(availableRoles, strings.Title(key))
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Role `%s` not found. Available roles: %s", strings.Join(args, " "), strings.Join(availableRoles, ", ")))
		return
	}

	// Assign the role
	assignRole(s, m, fullRoleName)
}

// handleIamNotCommand handles the .iamnot command for removing roles
func handleIamNotCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `.iamnot <role>` - Available roles: Frontend, Backend, Fullstack, Mobile")
		return
	}

	// Get the short role name (case-insensitive)
	shortName := strings.ToLower(strings.Join(args, " "))

	// Look up the full role name
	fullRoleName, exists := roleMapping[shortName]
	if !exists {
		// Build available roles message
		var availableRoles []string
		for key := range roleMapping {
			availableRoles = append(availableRoles, strings.Title(key))
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Role `%s` not found. Available roles: %s", strings.Join(args, " "), strings.Join(availableRoles, ", ")))
		return
	}

	// Remove the role
	removeRole(s, m, fullRoleName)
}

// assignRole assigns a role to the user
func assignRole(s *discordgo.Session, m *discordgo.MessageCreate, roleName string) {
	// Get the guild (server)
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error: Could not access server information.")
		return
	}

	// Find the role by name (case-insensitive)
	var targetRole *discordgo.Role
	for _, role := range guild.Roles {
		if strings.EqualFold(role.Name, roleName) {
			targetRole = role
			break
		}
	}

	if targetRole == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Role `%s` not found.", roleName))
		return
	}

	// Check if user already has the role
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Printf("Error getting member: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error: Could not access your member information.")
		return
	}

	for _, roleID := range member.Roles {
		if roleID == targetRole.ID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ You already have the role `%s`!", roleName))
			return
		}
	}

	// Assign the role
	err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, targetRole.ID)
	if err != nil {
		log.Printf("Error assigning role: %v", err)

		// Check if it's a permission error
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "Missing Permissions") {
			s.ChannelMessageSend(m.ChannelID, "❌ Error: Bot doesn't have permission to assign this role. Please make sure the bot's role is higher than the target role in the server settings.")
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Error assigning role: %v", err))
		}
		return
	}

	// Send confirmation message
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Role `%s` has been assigned to you!", roleName))
}

// removeRole removes a role from the user
func removeRole(s *discordgo.Session, m *discordgo.MessageCreate, roleName string) {
	// Get the guild (server)
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Printf("Error getting guild: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error: Could not access server information.")
		return
	}

	// Find the role by name (case-insensitive)
	var targetRole *discordgo.Role
	for _, role := range guild.Roles {
		if strings.EqualFold(role.Name, roleName) {
			targetRole = role
			break
		}
	}

	if targetRole == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Role `%s` not found.", roleName))
		return
	}

	// Check if user has the role
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Printf("Error getting member: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error: Could not access your member information.")
		return
	}

	hasRole := false
	for _, roleID := range member.Roles {
		if roleID == targetRole.ID {
			hasRole = true
			break
		}
	}

	if !hasRole {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ You don't have the role `%s`.", roleName))
		return
	}

	// Remove the role
	err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, targetRole.ID)
	if err != nil {
		log.Printf("Error removing role: %v", err)

		// Check if it's a permission error
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "Missing Permissions") {
			s.ChannelMessageSend(m.ChannelID, "❌ Error: Bot doesn't have permission to remove this role. Please make sure the bot's role is higher than the target role in the server settings.")
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Error removing role: %v", err))
		}
		return
	}

	// Send confirmation message
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅ Role `%s` has been removed from you!", roleName))
}

// handleHelpCommand shows help information
func handleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Build available roles list
	var availableRoles []string
	for key, value := range roleMapping {
		availableRoles = append(availableRoles, fmt.Sprintf("`%s` → %s", strings.Title(key), value))
	}

	helpText := "**Bot Commands:**\n\n" +
		"`.iam <role>` - Assign yourself a developer role\n" +
		"`.iamnot <role>` - Remove a developer role from yourself\n" +
		"`.help` - Show this help message\n\n" +
		"**Available Roles:**\n" +
		strings.Join(availableRoles, "\n") + "\n\n" +
		"**Examples:**\n" +
		"`.iam Frontend` - Assigns \"Frontend Developer\" role\n" +
		"`.iam Backend` - Assigns \"Backend Developer\" role\n" +
		"`.iamnot Frontend` - Removes \"Frontend Developer\" role"

	s.ChannelMessageSend(m.ChannelID, helpText)
}

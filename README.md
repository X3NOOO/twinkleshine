# ![logo](.github/assets/logo.png)

Twinkleshine is a Discord bot with AI-powered knowledge management capabilities that powers the BudCare Galaxy project.
It allows users to ask questions and receive answers based on a knowledge base that can be continuously expanded.

## Features

- **AI-Powered Responses**: Uses Google AI (or potentially other LLM providers) to generate responses to user questions
- **Knowledge Management**: 
  - Add text, files, URLs, or entire Discord channels to the bot's knowledge base
  - Automatically learns from messages sent by users with specific roles
  - Retrieval-Augmented Generation (RAG) for accurate, source-cited responses
- **Discord Integration**:
  - Slash commands for easy interaction
  - Message context menu commands
  - Role-based access control for sensitive operations
  - Cooldown system to prevent abuse
- **Vector Database Storage**:
  - Uses Qdrant for efficient similarity search
  - Stores document chunks with metadata for source attribution
- **Document Processing**:
  - Automatic content type detection
  - Parsing of various file formats using LlamaParse
  - Text chunking with configurable length and overlap

## Requirements

- Discord Bot Token
- LLM API access (currently supports Google AI)
- Vector Database (currently supports Qdrant)

## Configuration

### Environment Variables

Create a `.env` file based on the provided `.env.example`:

```
DISCORD_TOKEN=your_discord_bot_token

LLM_PROVIDER=google
LLM_MODEL=your_model_name
LLM_API_KEY=your_api_key

VDB_PROVIDER=qdrant
VDB_API_KEY=your_qdrant_api_key
VDB_HOST=your_qdrant_host
VDB_COLLECTION_NAME=your_collection_name
```

### Configuration File

Create a `config.yaml` file based on the provided `config.yaml.example`:

```yaml
system_prompt: |
  Your system prompt here that defines the bot's personality and behavior

discord:
  security:
    staff_role_id: your_staff_role_id
    cooldown_seconds: 60
  learn_messages_role_id: your_learn_messages_role_id

llm:
  max_tokens: 1024
  temperature: 0.6
  min_message_length: 12

rag:
  parse_timeout_seconds: 60
  chunking:
    length: 2048
    overlap: 256
  matches:
    root_count: 50
    count: 15
  rag_prompt: |
    Your RAG prompt here that defines how the bot should use retrieved information
```

> [!IMPORTANT]
> The `rag_prompt` makes use of a special formatting syntax to insert retrieved information into the prompt. You NEED to include `{RAG_KNOWLEDGE}` in the prompt to properly embed the RAG results. See the example file for more details.

## Deployment

### Local Development

1. Clone the repository
2. Create `.env` and `config.yaml` files as described above
3. Run the bot:
   ```
   go run .
   ```

### Docker Deployment

1. Clone the repository
2. Create `.env` and `config.yaml` files as described above
3. Build and run with Docker Compose:
   ```
   docker-compose up -d
   ```

The Docker setup includes:
- Multi-stage build for smaller image size
- Volume mounts for configuration files
- Automatic restart unless manually stopped

### Command Line Options

The bot supports the following command line options:

- `-env`: Path to the environment file (default: `.env`)
- `-config`: Path to the configuration file (default: `config.yaml`)
- `-verbose`: Enable verbose logging (default: false)

## Usage

### Discord Commands

- `/about`: Display information about the bot
- `/ask [question]`: Ask a question to the bot
- Right-click on a message > Apps > "Respond to the question": Have the bot respond to a specific message
- `/remember file [file]`: Add a file to the knowledge base
- `/remember text [text]`: Add text to the knowledge base
- `/remember urls [urls]`: Add websites to the knowledge base
- `/remember channel [channel]`: Add all attachments from a channel to the knowledge base
- Right-click on a message > Apps > "Add to the persistent knowledge": Add a specific message to the knowledge base

### Automatic Learning

The bot can automatically learn from messages sent by users with a specific role. Configure the `learn_messages_role_id` in the config file to enable this feature.

## License

See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

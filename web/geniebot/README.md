# GenieBot UI

A React + TypeScript frontend for AI service interactions.

## Features

- **Natural Language Chat**: Conversational interface with AI
- **Intent Detection**: Automatic recognition of user intent
- **Service Recommendations**: Smart suggestions for LLM/Agent/Workflow services
- **Task Management**: Real-time task progress tracking
- **Result Viewer**: View and download task outputs

## Tech Stack

- React 18
- TypeScript
- Vite
- Tailwind CSS
- Lucide React (icons)
- React Markdown (rendering)
- Vitest (testing)

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Build for production
npm run build
```

## Project Structure

```
src/
├── components/       # React components
│   ├── chat/        # Chat-related components
│   ├── layout/      # Layout components (Header, Sidebar)
│   ├── services/    # Service recommendation components
│   ├── tasks/       # Task management components
│   ├── results/     # Result viewer components
│   └── ui/          # UI components (Button, Card, etc.)
├── hooks/           # Custom React hooks
├── services/        # API services
├── types/           # TypeScript types
├── utils/           # Utility functions
└── test/            # Test setup
```

## API Integration

The frontend connects to the ShareToken blockchain API:

- `POST /api/chat/message` - Send chat message
- `POST /api/chat/intent` - Detect intent
- `POST /api/services/recommend` - Get service recommendations
- `POST /api/services/invoke` - Invoke a service
- `GET /api/tasks` - Get task list
- `WebSocket /websocket` - Real-time updates

## Components

### Chat Components
- `ChatContainer` - Main chat interface
- `ChatMessage` - Individual message display
- `ChatInput` - Message input with auto-resize

### Service Components
- `ServiceCard` - Service recommendation card
- `ServiceList` - Grid of service recommendations

### Task Components
- `TaskList` - Task management panel
- Progress tracking with real-time updates

### UI Components
- `Button` - Button with variants
- `Card` - Card container
- `Badge` - Status badges
- `Progress` - Progress bar
- `Dialog` - Modal dialog
- `ScrollArea` - Scrollable container

## Testing

Unit tests are written with Vitest and React Testing Library:

```bash
# Run tests
npm test

# Run tests with UI
npm run test:ui

# Type check
npm run typecheck

# Lint
npm run lint
```

## Environment Variables

```env
VITE_API_URL=http://localhost:1317
VITE_WS_URL=ws://localhost:26657
```

## License

MIT

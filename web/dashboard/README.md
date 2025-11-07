# NebulaBox Dashboard

Modern web dashboard for NebulaBox container management platform.

## Tech Stack

- **React 18** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI components
- **Lucide React** - Icons
- **React Router** - Routing

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
cd web/dashboard
npm install
```

### Development

```bash
npm run dev
```

The dashboard will be available at `http://localhost:3000`

### Build

```bash
npm run build
```

## Features

- 📊 **Dashboard** - System overview and quick actions
- 📦 **Container Management** - List, create, stop, view logs
- 🖼️ **Image Management** - Pull, push, and manage images
- 📈 **System Monitoring** - Real-time CPU, memory, and disk usage
- ⚙️ **Settings** - Configuration panel

## API Integration

The dashboard connects to the NebulaBox Go backend API. By default, it expects the API at `http://localhost:8080/api`.

You can configure the API URL using the `VITE_API_URL` environment variable:

```bash
VITE_API_URL=http://localhost:8080/api npm run dev
```

## Project Structure

```
src/
├── components/
│   ├── layout/        # Sidebar, dashboard layout
│   └── ui/            # shadcn/ui components
├── pages/             # Route pages
│   ├── Dashboard.tsx
│   ├── Containers.tsx
│   ├── CreateContainer.tsx
│   ├── ContainerLogs.tsx
│   ├── Images.tsx
│   ├── Monitor.tsx
│   └── Settings.tsx
├── lib/
│   ├── api.ts         # API client
│   └── utils.ts       # Utility functions
├── App.tsx            # Main app with routing
├── main.tsx           # Entry point
└── index.css          # Global styles
```

## Development Notes

Currently, the dashboard uses mock data when the API is unavailable. This allows for frontend development without requiring the full backend stack.

To integrate with the real API, you'll need to:
1. Create a REST API server in the Go backend
2. Update the API endpoints in `src/lib/api.ts`
3. Remove mock data fallbacks

## Next Steps

- [ ] Connect to real NebulaBox API
- [ ] Add authentication
- [ ] Implement container logs streaming
- [ ] Add container terminal access
- [ ] Enhanced error handling
- [ ] Dark mode toggle
- [ ] Export logs functionality


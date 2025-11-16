# Hobby Tracker

A cross-platform web application for tracking and organizing hobbies with AI-powered categorization and multi-source import capabilities.

## Features

### Core Functionality

- **AI-Powered Categorization**: Automatically categorizes items (movies, restaurants, travel destinations, music, activities, etc.)
- **Multi-Source Import**: Import from Instagram, YouTube, X (Twitter), TikTok, Telegram, Wikipedia, and web links
- **Smart Tagging**: AI learns user preferences and suggests relevant tags
- **List Management**: Create and organize custom lists with different categories
- **Circles**: Package lists together and share with specific groups (Partner, Friends, Family, etc.)
- **Progress Tracking**: Check off completed items with metadata tracking
- **Calendar Integration**: Plan activities with due dates (ready for calendar sync)
- **Access Control**: Share lists and circles with different permission levels

### User Interface

- Clean, distraction-free design inspired by Google Tasks
- Material Design principles
- List-based layout (no visual clutter)
- Time-based grouping (This Week, Last Week, Earlier)
- Detail panel for rich item information
- Images hidden in list view, shown in detail view

## Tech Stack

- **Framework**: React 18 with TypeScript
- **Routing**: React Router v6
- **State Management**: Zustand
- **Styling**: Tailwind CSS
- **Build Tool**: Vite
- **Icons**: Lucide React
- **Date Handling**: date-fns

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn
- Modern web browser

### Installation

1. Clone the repository:

```bash
git clone https://github.com/debMan/collaborative-hobby-tracker.git
cd collaborative-hobby-tracker
```

2. Install dependencies:

```bash
npm install
```

3. Create environment file:

```bash
cp .env.example .env
```

4. Start development server:

```bash
npm run dev
```

5. Open your browser and navigate to `http://localhost:3000`

### Default Test Account

For development purposes, you can use any email/password combination. The app uses mock authentication.

Example:

- Email: `test@example.com`
- Password: `password123`

## Project Structure

```
hobby-tracker-app/
├── src/
│   ├── components/          # React components
│   │   ├── auth/           # Authentication components
│   │   ├── items/          # Item-related components
│   │   ├── layout/         # Layout components (Header, Sidebar)
│   │   └── modals/         # Modal dialogs
│   ├── pages/              # Page components
│   │   ├── Dashboard.tsx   # Main dashboard
│   │   ├── Login.tsx       # Login page
│   │   └── Register.tsx    # Registration page
│   ├── services/           # API services
│   │   └── api.ts          # API service layer with mock implementations
│   ├── store/              # State management
│   │   └── index.ts        # Zustand store
│   ├── types/              # TypeScript type definitions
│   │   └── index.ts        # App-wide types
│   ├── utils/              # Utility functions
│   │   └── mockData.ts     # Mock data for development
│   ├── App.tsx             # Root component with routing
│   ├── main.tsx            # Application entry point
│   └── index.css           # Global styles
├── public/                 # Static assets
├── index.html              # HTML template
├── package.json            # Dependencies and scripts
├── tsconfig.json           # TypeScript configuration
├── tailwind.config.js      # Tailwind CSS configuration
├── vite.config.ts          # Vite configuration
└── README.md               # This file
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Backend Integration

### Current State

The application currently uses **mock data** stored in `localStorage`. All API calls in `src/services/api.ts` are simulated with delays to mimic network requests.

### Backend Connection (When Ready)

The application is designed to easily connect to a real backend. All backend connections are **commented out** in `src/services/api.ts` with clear markers:

```typescript
/* BACKEND CONNECTION - Uncomment when backend is ready
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...
  });
*/
```

#### Steps to Connect Backend

1. **Set API URL**:
   - Update `VITE_API_URL` in `.env` file
   - Example: `VITE_API_URL=https://api.hobbytracker.com/v1`

2. **Uncomment Backend Calls**:
   - In `src/services/api.ts`, uncomment all backend connection blocks
   - Remove or comment out mock implementations

3. **Configure OAuth**:
   - Set up OAuth client IDs in `.env`:
     - `VITE_GOOGLE_CLIENT_ID`
     - `VITE_APPLE_CLIENT_ID`
     - `VITE_TWITTER_CLIENT_ID`

4. **Update OAuth Flow**:
   - Implement OAuth redirect handling
   - Update OAuth login functions in `src/services/api.ts`

### Backend API Endpoints Expected

#### Authentication

- `POST /api/auth/login` - Email/password login
- `POST /api/auth/register` - User registration
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user
- `GET /api/auth/oauth/{provider}` - OAuth login

#### Items

- `GET /api/items` - Get all items
- `GET /api/lists/{listId}/items` - Get items by list
- `POST /api/items` - Create item
- `PATCH /api/items/{id}` - Update item
- `DELETE /api/items/{id}` - Delete item
- `POST /api/items/{id}/toggle` - Toggle completion

#### Import

- `POST /api/import` - Import item with AI categorization

#### Lists

- `GET /api/lists` - Get all lists
- `POST /api/lists` - Create list
- `PATCH /api/lists/{id}` - Update list
- `DELETE /api/lists/{id}` - Delete list

#### Circles

- `GET /api/circles` - Get all circles
- `POST /api/circles` - Create circle
- `PATCH /api/circles/{id}` - Update circle
- `DELETE /api/circles/{id}` - Delete circle

#### Tags

- `GET /api/tags` - Get all tags
- `POST /api/tags/suggest` - Get AI tag suggestions

## Features to Implement (Backend Required)

The following features are designed but require backend implementation:

1. **Real AI Categorization**: Replace mock keyword detection with actual AI models
2. **Image Fetching**: Fetch images from source URLs (Instagram, YouTube, etc.)
3. **External Data Integration**:
   - IMDB for movies
   - Google Maps for locations
   - Music APIs for albums/artists
4. **Calendar Integration**: Sync with Google Calendar, Apple Calendar
5. **Real-time Collaboration**: Share and edit lists with others
6. **Push Notifications**: Due date reminders
7. **Search and Filters**: Full-text search across items

## Mock Data

During development, the app uses mock data defined in `src/utils/mockData.ts`. This includes:

- Sample users
- Pre-populated items (movies, restaurants, travel destinations, etc.)
- Sample lists and circles
- Sample tags

Mock data is stored in `localStorage` and persists between sessions.

## Design Philosophy

- **Simplicity First**: Clean, distraction-free interface
- **Content-First**: Focus on hobby items, not decorative UI
- **List-Based**: No cards, minimal visual noise
- **Hidden Complexity**: Advanced features accessible but not overwhelming
- **Responsive**: Works on desktop and mobile (mobile-optimized version coming)

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Future Roadmap

- [ ] Native mobile apps (iOS, Android)
- [ ] Offline support with PWA
- [ ] Advanced filtering and search
- [ ] Export functionality (PDF, CSV)
- [ ] Browser extensions for quick imports
- [ ] API for third-party integrations
- [ ] Analytics and insights dashboard
- [ ] Recommendations based on history

## Support

For issues and questions, please open an issue on GitHub or contact the development team.

---

Built with ❤️ for hobby enthusiasts

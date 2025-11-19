# Hobby Tracker - Quick Setup Guide

## 🚀 Getting Started in 3 Steps

### 1. Install Dependencies

```bash
cd collaborative-hobby-tracker
npm install
```

### 2. Start Development Server

```bash
npm run dev
```

### 3. Open in Browser

Navigate to: <http://localhost:3000>

## 🎯 First Login

The app uses mock authentication. You can log in with ANY email/password:

**Example credentials:**

- Email: `demo@example.com`
- Password: `password`

Or click "Continue with Google/Apple/X" (also mocked for now)

## 📱 What You'll See

After logging in, you'll see:

- **Sidebar**: Pre-populated with sample lists (Movies, Restaurants, Travel, Music)
- **Main Content**: Sample hobby items organized by time
- **Add Item**: Click the input box to import or add new items

## 🧪 Try These Features

1. **Add a New Item**:
   - Click the "Add a new item or paste a link..." box
   - Enter text like "Visit Tokyo" or paste a URL
   - Click "Analyze" to see AI categorization
   - Choose category and add tags
   - Click "Add Item"

2. **Complete an Item**:
   - Click the checkbox next to any item
   - It will be marked as complete

3. **View Details**:
   - Click on any item to see full details in the right panel
   - View images, description, tags, and metadata

4. **Switch Lists**:
   - Click on different lists in the sidebar
   - See items filtered by category

## 🔧 Current Status

**✅ What Works (Mock Data)**

- Authentication (all login methods work with mock data)
- Adding/editing/deleting items
- List management
- AI-powered categorization (simple keyword-based mock)
- Tag management
- Item completion tracking
- Detail panel views

**🔄 What Needs Backend**

- Real AI categorization
- Image fetching from URLs
- External API integration (IMDB, Google Maps, etc.)
- Calendar synchronization
- Real-time collaboration
- OAuth with real providers

## 📂 Project Structure

```
hobby-tracker-app/
├── src/
│   ├── components/     # UI components
│   ├── pages/          # Main pages (Dashboard, Login, Register)
│   ├── services/       # API layer (mock implementations)
│   ├── store/          # State management (Zustand)
│   ├── types/          # TypeScript types
│   └── utils/          # Mock data and utilities
├── package.json
├── README.md          # Full documentation
└── .env.example       # Environment variables template
```

## 🔌 Connecting to Backend

When your backend is ready:

1. **Update .env file**:

```bash
cp .env.example .env
# Edit .env and set:
VITE_API_URL=http://localhost:8080/api
```

2. **Uncomment Backend Calls**:

- Open `src/services/api.ts`
- Find all `/* BACKEND CONNECTION - Uncomment when backend is ready */`
- Uncomment those blocks
- Comment out or remove the mock implementations

3. **Configure OAuth**:

- Add your OAuth client IDs to `.env`
- Update OAuth redirect URLs

## 🎨 Customization

### Colors

Edit `tailwind.config.js` to change the primary color scheme:

```javascript
colors: {
  primary: {
    // Change these values
    600: '#1a73e8',  // Main blue
    700: '#1967d2',  // Darker blue
    // ...
  }
}
```

### Mock Data

Edit `src/utils/mockData.ts` to change sample items, lists, and circles.

## 📖 Full Documentation

See `README.md` for complete documentation including:

- Detailed feature list
- Backend API endpoints specification
- Architecture decisions
- Contributing guidelines

## 🐛 Troubleshooting

**Port 3000 already in use?**

```bash
# Vite will automatically try the next available port
# Or specify a different port in vite.config.ts
```

**Dependencies not installing?**

```bash
# Try clearing npm cache
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

**TypeScript errors?**

```bash
# Make sure you're using TypeScript 5.2+
npm install -D typescript@latest
```

## 🎉 You're All Set

Start building your hobby tracking experience!

For questions or issues, check the README.md or open an issue.

---

Happy tracking! 🚀

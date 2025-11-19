import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { migrateLocalStorageData } from './utils/migrations'

// Run data migrations before app starts
migrateLocalStorageData();

// Apply dark mode based on localStorage preference or system preference
const getInitialDarkMode = () => {
  const stored = localStorage.getItem('darkMode');
  if (stored !== null) {
    return stored === 'true';
  }
  // Default to system preference if no stored preference
  return window.matchMedia('(prefers-color-scheme: dark)').matches;
};

if (getInitialDarkMode()) {
  document.documentElement.classList.add('dark');
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

import { Menu } from 'lucide-react';
import { useStore } from '../../store';

export default function Header() {
  const { user, logout, toggleMobileMenu } = useStore();

  const handleLogout = async () => {
    if (confirm('Are you sure you want to log out?')) {
      await logout();
    }
  };

  const initials = user?.name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase() || 'U';

  return (
    <header className="bg-white border-b border-gray-200 px-4 sm:px-6 py-3 sm:py-4 flex items-center justify-between">
      <div className="flex items-center gap-3 sm:gap-4">
        {/* Mobile menu button */}
        <button
          onClick={toggleMobileMenu}
          className="lg:hidden p-2 hover:bg-gray-100 rounded-lg transition"
          aria-label="Toggle menu"
        >
          <Menu className="w-5 h-5 text-gray-600" />
        </button>

        <h1 className="text-lg sm:text-xl text-gray-600 font-normal">Hobby Tracker</h1>
      </div>

      <div className="relative group">
        <button
          className="w-8 h-8 sm:w-9 sm:h-9 rounded-full bg-primary-600 text-white flex items-center justify-center text-xs sm:text-sm font-medium cursor-pointer hover:bg-primary-700 transition"
        >
          {initials}
        </button>

        {/* Dropdown menu */}
        <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
          <div className="py-2 px-4 border-b border-gray-100">
            <div className="font-medium text-sm">{user?.name}</div>
            <div className="text-xs text-gray-500">{user?.email}</div>
          </div>
          <div className="py-1">
            <button
              onClick={handleLogout}
              className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition"
            >
              Sign out
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}

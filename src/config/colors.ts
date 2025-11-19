// Centralized color configuration for the application
// All UI colors should be defined here for easy customization

export const colors = {
  // Badge colors
  badges: {
    category: {
      bg: 'bg-gray-100',
      text: 'text-gray-700',
      border: 'border-gray-200',
    },
    circle: {
      bg: 'bg-emerald-100',
      text: 'text-emerald-700',
      border: 'border-emerald-200',
    },
    source: {
      bg: 'bg-rose-100',
      text: 'text-rose-700',
      border: 'border-rose-200',
    },
    tag: {
      bg: 'bg-sky-50',
      text: 'text-sky-700',
      border: 'border-sky-200',
    },
  },

  // Sidebar filter colors
  sidebar: {
    circle: {
      selected: {
        bg: 'bg-emerald-100',
        text: 'text-emerald-700',
        border: 'border-emerald-300',
      },
      unselected: {
        bg: 'bg-transparent',
        text: 'text-gray-700',
        border: 'border-transparent',
      },
      hover: {
        bg: 'hover:bg-gray-100',
      },
    },
    source: {
      selected: {
        bg: 'bg-rose-100',
        text: 'text-rose-700',
        border: 'border-rose-300',
      },
      unselected: {
        bg: 'bg-transparent',
        text: 'text-gray-700',
        border: 'border-transparent',
      },
      hover: {
        bg: 'hover:bg-gray-100',
      },
    },
  },

  // Primary colors (existing)
  primary: {
    50: '#eff6ff',
    100: '#dbeafe',
    200: '#bfdbfe',
    300: '#93c5fd',
    400: '#60a5fa',
    500: '#3b82f6',
    600: '#2563eb',
    700: '#1d4ed8',
    800: '#1e40af',
    900: '#1e3a8a',
  },
};

// Helper function to get all classes for a badge type
export const getBadgeClasses = (type: 'category' | 'circle' | 'source' | 'tag') => {
  const badge = colors.badges[type];
  const darkVariants = {
    category: 'dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600',
    circle: 'dark:bg-emerald-900 dark:text-emerald-300 dark:border-emerald-800',
    source: 'dark:bg-rose-900 dark:text-rose-300 dark:border-rose-800',
    tag: 'dark:bg-sky-900 dark:text-sky-300 dark:border-sky-800',
  };
  return `${badge.bg} ${badge.text} ${badge.border} ${darkVariants[type]}`;
};

// Helper function to get sidebar filter classes
export const getSidebarFilterClasses = (
  type: 'circle' | 'source',
  isSelected: boolean
) => {
  const filter = colors.sidebar[type];
  const darkVariants = {
    circle: {
      selected: 'dark:bg-emerald-900 dark:text-emerald-300 dark:border-emerald-800',
      unselected: 'dark:text-gray-300 dark:hover:bg-gray-800',
    },
    source: {
      selected: 'dark:bg-rose-900 dark:text-rose-300 dark:border-rose-800',
      unselected: 'dark:text-gray-300 dark:hover:bg-gray-800',
    },
  };

  if (isSelected) {
    return `${filter.selected.bg} ${filter.selected.text} ${filter.selected.border} ${darkVariants[type].selected}`;
  }
  return `${filter.unselected.bg} ${filter.unselected.text} ${filter.unselected.border} ${filter.hover.bg} ${darkVariants[type].unselected}`;
};

import { useMemo } from 'react';
import { useStore } from '../../store';
import { HobbyItem } from '../../types';
import { formatDistanceToNow } from 'date-fns';
import { MoreVertical } from 'lucide-react';

export default function ItemList() {
  const {
    categories,
    selectedCategoryTab,
    setSelectedCategoryTab,
    getFilteredItems,
    getCategoriesByFilters,
    openDetailPanel,
    toggleItemComplete,
    openImportModal,
  } = useStore();

  const filteredCategories = getCategoriesByFilters();
  const filteredItems = getFilteredItems();

  // Group items by time periods
  const groupedItems = useMemo(() => {
    const now = new Date();
    const oneWeekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
    const twoWeeksAgo = new Date(now.getTime() - 14 * 24 * 60 * 60 * 1000);

    const thisWeek: HobbyItem[] = [];
    const lastWeek: HobbyItem[] = [];
    const earlier: HobbyItem[] = [];

    filteredItems.forEach((item) => {
      const addedDate = new Date(item.addedAt);
      if (addedDate > oneWeekAgo) {
        thisWeek.push(item);
      } else if (addedDate > twoWeeksAgo) {
        lastWeek.push(item);
      } else {
        earlier.push(item);
      }
    });

    return { thisWeek, lastWeek, earlier };
  }, [filteredItems]);

  const getSourceIcon = (source: string) => {
    const icons: Record<string, string> = {
      youtube: '📺',
      instagram: '📷',
      twitter: '🐦',
      tiktok: '🎵',
      telegram: '✈️',
      web: '🌐',
      manual: '✍️',
      wikipedia: '📖',
    };
    return icons[source] || '📱';
  };

  const getCategoryById = (categoryId: string) => {
    return categories.find((cat) => cat.id === categoryId);
  };

  const renderItem = (item: HobbyItem) => {
    const category = getCategoryById(item.categoryId);

    return (
      <div
        key={item.id}
        className="group flex items-start gap-3 py-3 px-2 -mx-2 rounded hover:bg-gray-50 active:bg-gray-100 transition cursor-pointer border-b border-gray-100"
        onClick={() => openDetailPanel(item)}
      >
        {/* Checkbox */}
        <button
          onClick={(e) => {
            e.stopPropagation();
            toggleItemComplete(item.id);
          }}
          className={`mt-0.5 w-[22px] h-[22px] sm:w-[18px] sm:h-[18px] rounded-full border-2 flex items-center justify-center flex-shrink-0 transition ${
            item.isCompleted
              ? 'bg-primary-600 border-primary-600'
              : 'border-gray-300 hover:border-primary-600 active:border-primary-700'
          }`}
        >
          {item.isCompleted && (
            <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
            </svg>
          )}
        </button>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className={`text-sm mb-1 ${item.isCompleted ? 'line-through text-gray-400' : 'text-gray-900'}`}>
            {item.title}
          </div>
          <div className="flex items-center gap-3 flex-wrap text-xs text-gray-500">
            {category && (
              <span className="flex items-center gap-1">
                {category.icon} {category.name}
              </span>
            )}
            <span className="flex items-center gap-1">
              📅 {formatDistanceToNow(new Date(item.addedAt), { addSuffix: true })}
            </span>
            <span className="flex items-center gap-1">
              {getSourceIcon(item.source)} From {item.source.charAt(0).toUpperCase() + item.source.slice(1)}
            </span>
            {item.tags.slice(0, 2).map((tag) => (
              <span key={tag} className="px-2 py-0.5 bg-primary-50 text-primary-700 rounded-full text-[11px]">
                {tag}
              </span>
            ))}
          </div>
        </div>

        {/* Actions */}
        <button
          onClick={(e) => {
            e.stopPropagation();
            // Handle more options
          }}
          className="opacity-0 group-hover:opacity-100 p-1.5 hover:bg-gray-100 rounded transition"
        >
          <MoreVertical className="w-4 h-4 text-gray-500" />
        </button>
      </div>
    );
  };

  const renderTimeGroup = (title: string, items: HobbyItem[]) => {
    if (items.length === 0) return null;

    return (
      <div className="mb-8">
        <h3 className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-3">
          {title}
        </h3>
        <div className="space-y-0.5">
          {items.map(renderItem)}
        </div>
      </div>
    );
  };

  const selectedCategory = selectedCategoryTab !== 'all'
    ? categories.find((cat) => cat.id === selectedCategoryTab)
    : null;

  return (
    <main className="flex-1 bg-white overflow-y-auto flex flex-col">
      {/* Category Tabs */}
      <div className="border-b border-gray-200 px-4 sm:px-6 lg:px-8">
        <div className="flex gap-1 overflow-x-auto scrollbar-hide -mb-px">
          <button
            onClick={() => setSelectedCategoryTab('all')}
            className={`flex-shrink-0 px-4 py-3 text-sm font-medium border-b-2 transition whitespace-nowrap ${
              selectedCategoryTab === 'all'
                ? 'border-primary-600 text-primary-700'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            All Items
          </button>
          {filteredCategories.map((category) => (
            <button
              key={category.id}
              onClick={() => setSelectedCategoryTab(category.id)}
              className={`flex-shrink-0 flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 transition whitespace-nowrap ${
                selectedCategoryTab === category.id
                  ? 'border-primary-600 text-primary-700'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              <span>{category.icon}</span>
              <span>{category.name}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Header */}
      <div className="border-b border-gray-200 px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
        <h2 className="text-xl sm:text-2xl font-normal text-gray-900 mb-2">
          {selectedCategory ? (
            <span className="flex items-center gap-2">
              <span>{selectedCategory.icon}</span>
              <span>{selectedCategory.name}</span>
            </span>
          ) : (
            'All Items'
          )}
        </h2>
        <p className="text-sm text-gray-500">
          {filteredItems.length} item{filteredItems.length !== 1 ? 's' : ''} • Last updated today
        </p>
      </div>

      {/* Add Item Section */}
      <div className="border-b border-gray-200 px-4 sm:px-6 lg:px-8 py-4 sm:py-5">
        <div
          onClick={openImportModal}
          className="flex items-center gap-3 px-3 sm:px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 hover:bg-white hover:shadow-sm active:bg-white transition cursor-pointer"
        >
          <span className="text-primary-600 text-xl">+</span>
          <span className="flex-1 text-sm text-gray-500">Add a new item or paste a link...</span>
          <button className="hidden sm:block px-3 py-1.5 text-xs border border-gray-300 rounded bg-white text-gray-600 hover:border-primary-600 hover:text-primary-600 transition">
            Import
          </button>
        </div>
      </div>

      {/* Items List */}
      <div className="flex-1 px-4 sm:px-6 lg:px-8 py-4 sm:py-6 overflow-y-auto">
        {filteredItems.length === 0 ? (
          <div className="text-center py-16 text-gray-400">
            <div className="text-5xl mb-4 opacity-30">📋</div>
            <div className="text-base mb-2">No items yet</div>
            <div className="text-sm">Add your first hobby item to get started!</div>
          </div>
        ) : (
          <>
            {renderTimeGroup('This Week', groupedItems.thisWeek)}
            {renderTimeGroup('Last Week', groupedItems.lastWeek)}
            {renderTimeGroup('Earlier', groupedItems.earlier)}
          </>
        )}
      </div>
    </main>
  );
}

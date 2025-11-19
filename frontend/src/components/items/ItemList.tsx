import { useMemo } from 'react';
import { useStore } from '../../store';
import { HobbyItem } from '../../types';
import { MoreVertical } from 'lucide-react';
import { getBadgeClasses } from '../../config/colors';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faInstagram, faYoutube, faXTwitter, faTiktok, faTelegram, faWikipediaW } from '@fortawesome/free-brands-svg-icons';
import { faGlobe, faPen } from '@fortawesome/free-solid-svg-icons';

export default function ItemList() {
  const {
    categories,
    circles,
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

  // Check if there are duplicate category names to show circle badges
  const categoryNameCounts = useMemo(() => {
    const counts: Record<string, number> = {};
    filteredCategories.forEach(cat => {
      counts[cat.name] = (counts[cat.name] || 0) + 1;
    });
    return counts;
  }, [filteredCategories]);

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
    const icons: Record<string, any> = {
      youtube: faYoutube,
      instagram: faInstagram,
      twitter: faXTwitter,
      tiktok: faTiktok,
      telegram: faTelegram,
      web: faGlobe,
      manual: faPen,
      wikipedia: faWikipediaW,
    };
    return icons[source] || faGlobe;
  };

  const getCategoryById = (categoryId: string) => {
    return categories.find((cat) => cat.id === categoryId);
  };

  const renderItem = (item: HobbyItem) => {
    // Get the category for this item
    const category = getCategoryById(item.categoryId);

    // Get the circle for this item
    const circle = category?.circleId
      ? circles.find(c => c.id === category.circleId)
      : null;

    return (
      <div
        key={item.id}
        className="group flex items-start gap-3 py-3 px-2 -mx-2 rounded hover:bg-gray-50 dark:hover:bg-gray-800 active:bg-gray-100 dark:active:bg-gray-700 transition cursor-pointer border-b border-gray-100 dark:border-gray-800"
        onClick={() => openDetailPanel(item)}
      >
        {/* Checkbox */}
        <button
          onClick={(e) => {
            e.stopPropagation();
            toggleItemComplete(item.id);
          }}
          className={`mt-0.5 w-[22px] h-[22px] sm:w-[18px] sm:h-[18px] rounded-full border-2 flex items-center justify-center flex-shrink-0 transition ${item.isCompleted
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
          <div className={`text-sm mb-1.5 ${item.isCompleted ? 'line-through text-gray-400 dark:text-gray-500' : 'text-gray-900 dark:text-gray-100'}`}>
            {item.title}
          </div>

          {/* Aligned Badge Grid */}
          <div className="grid grid-cols-[auto,auto,auto,1fr] gap-x-2 gap-y-1 items-center text-[11px]">
            {/* Category Badge - shown in All Items view */}
            {category && selectedCategoryTab === 'all' && (
              <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md font-medium whitespace-nowrap ${getBadgeClasses('category')}`}>
                <span>{category.icon}</span>
                <span>{category.name}</span>
              </span>
            )}

            {/* Circle Badge */}
            {circle && (
              <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md font-medium whitespace-nowrap ${getBadgeClasses('circle')}`}>
                <span>{circle.icon}</span>
                <span>{circle.name}</span>
              </span>
            )}

            {/* Source Badge */}
            <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md font-medium whitespace-nowrap ${getBadgeClasses('source')}`}>
              <FontAwesomeIcon icon={getSourceIcon(item.source)} className="w-3 h-3" />
              <span>{item.source.charAt(0).toUpperCase() + item.source.slice(1)}</span>
            </span>

            {/* Tags */}
            <div className="flex items-center gap-1 flex-wrap">
              {item.tags.slice(0, 2).map((tag) => (
                <span key={tag} className={`px-2 py-0.5 rounded-full whitespace-nowrap ${getBadgeClasses('tag')}`}>
                  {tag}
                </span>
              ))}
            </div>
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
        <h3 className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-3">
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
    <main className="flex-1 bg-white dark:bg-gray-900 overflow-y-auto flex flex-col">
      {/* Category Tabs */}
      <div className="border-b border-gray-200 dark:border-gray-700 px-4 sm:px-6 lg:px-8">
        <div className="flex gap-1 overflow-x-auto scrollbar-hide -mb-px">
          <button
            onClick={() => setSelectedCategoryTab('all')}
            className={`flex-shrink-0 px-4 py-3 text-sm font-medium border-b-2 transition whitespace-nowrap ${selectedCategoryTab === 'all'
              ? 'border-primary-600 text-primary-700 dark:text-primary-400'
              : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
              }`}
          >
            All Items
          </button>
          {filteredCategories.map((category) => {
            const hasDuplicateName = categoryNameCounts[category.name] > 1;
            const circle = category.circleId ? circles.find(c => c.id === category.circleId) : null;

            return (
              <button
                key={category.id}
                onClick={() => setSelectedCategoryTab(category.id)}
                className={`flex-shrink-0 flex items-center gap-2 px-4 py-3 text-sm font-medium border-b-2 transition whitespace-nowrap ${selectedCategoryTab === category.id
                  ? 'border-primary-600 text-primary-700 dark:text-primary-400'
                  : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
                  }`}
              >
                <span>{category.icon}</span>
                <span>{category.name}</span>
                {hasDuplicateName && (
                  <span className="inline-flex items-center gap-1 px-1.5 py-0.5 bg-blue-50 dark:bg-blue-900 text-blue-700 dark:text-blue-300 rounded text-[10px] font-medium">
                    {circle ? circle.icon : '👤'}
                  </span>
                )}
              </button>
            );
          })}
        </div>
      </div>

      {/* Header */}
      <div className="border-b border-gray-200 dark:border-gray-700 px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
        <h2 className="text-xl sm:text-2xl font-normal text-gray-900 dark:text-gray-100 mb-2">
          {selectedCategory ? (
            <span className="flex items-center gap-2">
              <span>{selectedCategory.icon}</span>
              <span>{selectedCategory.name}</span>
            </span>
          ) : (
            'All Items'
          )}
        </h2>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          {filteredItems.length} item{filteredItems.length !== 1 ? 's' : ''} • Last updated today
        </p>
      </div>

      {/* Add Item Section */}
      <div className="border-b border-gray-200 dark:border-gray-700 px-4 sm:px-6 lg:px-8 py-4 sm:py-5">
        <div
          onClick={openImportModal}
          className="flex items-center gap-3 px-3 sm:px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-800 hover:bg-white dark:hover:bg-gray-700 hover:shadow-sm active:bg-white dark:active:bg-gray-700 transition cursor-pointer"
        >
          <span className="text-primary-600 dark:text-primary-400 text-xl">+</span>
          <span className="flex-1 text-sm text-gray-500 dark:text-gray-400">Add a new item or paste a link...</span>
          <button className="hidden sm:block px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 hover:border-primary-600 dark:hover:border-primary-500 hover:text-primary-600 dark:hover:text-primary-400 transition">
            Import
          </button>
        </div>
      </div>

      {/* Items List */}
      <div className="flex-1 px-4 sm:px-6 lg:px-8 py-4 sm:py-6 overflow-y-auto">
        {filteredItems.length === 0 ? (
          <div className="text-center py-16 text-gray-400 dark:text-gray-500">
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

import { X, ChevronDown, ChevronRight } from 'lucide-react';
import { useStore } from '../../store';
import { DataSource } from '../../types';

export default function Sidebar() {
  const {
    circles,
    items,
    selectedSources,
    selectedCircles,
    isSourcesExpanded,
    isCirclesExpanded,
    toggleSourceFilter,
    toggleCircleFilter,
    toggleSourcesExpanded,
    toggleCirclesExpanded,
    isMobileMenuOpen,
    closeMobileMenu,
  } = useStore();

  const allItemsCount = items.length;

  const sources: { value: DataSource; label: string; icon: string }[] = [
    { value: 'instagram', label: 'Instagram', icon: '📷' },
    { value: 'youtube', label: 'YouTube', icon: '📺' },
    { value: 'twitter', label: 'X (Twitter)', icon: '🐦' },
    { value: 'tiktok', label: 'TikTok', icon: '🎵' },
    { value: 'telegram', label: 'Telegram', icon: '✈️' },
    { value: 'web', label: 'Web', icon: '🌐' },
    { value: 'manual', label: 'Manual', icon: '✍️' },
    { value: 'wikipedia', label: 'Wikipedia', icon: '📖' },
  ];

  const handleSourceToggle = (source: DataSource) => {
    toggleSourceFilter(source);
  };

  const handleCircleToggle = (circleId: string) => {
    toggleCircleFilter(circleId);
    closeMobileMenu();
  };

  return (
    <>
      {/* Overlay for mobile */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={closeMobileMenu}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`
          fixed lg:static inset-y-0 left-0 z-50
          w-[280px] bg-white border-r border-gray-200 overflow-y-auto
          transform transition-transform duration-300 ease-in-out
          lg:transform-none
          ${isMobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
        `}
      >
        {/* Mobile close button */}
        <div className="lg:hidden flex justify-end p-4 border-b border-gray-200">
          <button
            onClick={closeMobileMenu}
            className="p-2 hover:bg-gray-100 rounded-lg transition"
            aria-label="Close menu"
          >
            <X className="w-5 h-5 text-gray-600" />
          </button>
        </div>

        <div className="p-6 space-y-6">
          {/* All Items Summary */}
          <div className="px-3 py-2 bg-gray-50 rounded-lg">
            <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-1">
              Total Items
            </div>
            <div className="text-2xl font-normal text-gray-900">{allItemsCount}</div>
          </div>

          {/* Sources Section */}
          <div>
            <button
              onClick={toggleSourcesExpanded}
              className="w-full flex items-center gap-2 px-3 py-2 hover:bg-gray-50 rounded-lg transition"
            >
              {isSourcesExpanded ? (
                <ChevronDown className="w-4 h-4 text-gray-500" />
              ) : (
                <ChevronRight className="w-4 h-4 text-gray-500" />
              )}
              <h2 className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                Sources
              </h2>
              {selectedSources.length > 0 && (
                <span className="ml-auto text-xs bg-primary-100 text-primary-700 px-2 py-0.5 rounded-full">
                  {selectedSources.length}
                </span>
              )}
            </button>

            {isSourcesExpanded && (
              <nav className="mt-2 space-y-1 px-3">
                {sources.map((source) => (
                  <label
                    key={source.value}
                    className="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-gray-50 cursor-pointer transition"
                  >
                    <input
                      type="checkbox"
                      checked={selectedSources.includes(source.value)}
                      onChange={() => handleSourceToggle(source.value)}
                      className="w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-2 focus:ring-primary-500"
                    />
                    <span className="text-base">{source.icon}</span>
                    <span className="flex-1 text-sm text-gray-700">{source.label}</span>
                  </label>
                ))}
              </nav>
            )}
          </div>

          {/* Circles Section */}
          <div>
            <button
              onClick={toggleCirclesExpanded}
              className="w-full flex items-center gap-2 px-3 py-2 hover:bg-gray-50 rounded-lg transition"
            >
              {isCirclesExpanded ? (
                <ChevronDown className="w-4 h-4 text-gray-500" />
              ) : (
                <ChevronRight className="w-4 h-4 text-gray-500" />
              )}
              <h2 className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                Circles
              </h2>
              {selectedCircles.length > 0 && (
                <span className="ml-auto text-xs bg-primary-100 text-primary-700 px-2 py-0.5 rounded-full">
                  {selectedCircles.length}
                </span>
              )}
            </button>

            {isCirclesExpanded && (
              <nav className="mt-2 space-y-1 px-3">
                {circles.length === 0 ? (
                  <div className="px-2 py-3 text-sm text-gray-400 text-center">
                    No circles yet
                  </div>
                ) : (
                  circles.map((circle) => (
                    <label
                      key={circle.id}
                      className="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-gray-50 cursor-pointer transition"
                    >
                      <input
                        type="checkbox"
                        checked={selectedCircles.includes(circle.id)}
                        onChange={() => handleCircleToggle(circle.id)}
                        className="w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-2 focus:ring-primary-500"
                      />
                      <span className="text-base">{circle.icon}</span>
                      <span className="flex-1 text-sm text-gray-700">{circle.name}</span>
                    </label>
                  ))
                )}
              </nav>
            )}
          </div>
        </div>
      </aside>
    </>
  );
}

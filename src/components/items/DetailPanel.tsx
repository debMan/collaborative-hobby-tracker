import { X } from 'lucide-react';
import { useStore } from '../../store';
import { formatDistanceToNow, format } from 'date-fns';

export default function DetailPanel() {
  const { selectedItem, isDetailPanelOpen, closeDetailPanel, updateItem, deleteItem, categories, circles } = useStore();

  if (!isDetailPanelOpen || !selectedItem) {
    return null;
  }

  const category = categories.find(cat => cat.id === selectedItem.categoryId);
  const itemCircles = circles.filter(circle => selectedItem.circleIds.includes(circle.id));

  const handleDelete = async () => {
    if (confirm('Are you sure you want to delete this item?')) {
      await deleteItem(selectedItem.id);
      closeDetailPanel();
    }
  };

  return (
    <>
      {/* Mobile overlay */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
        onClick={closeDetailPanel}
      />

      {/* Detail Panel */}
      <aside className="fixed lg:static inset-0 lg:w-[360px] bg-white lg:border-l border-gray-200 overflow-y-auto z-50 lg:z-auto">
        <div className="p-4 sm:p-6">
        {/* Close Button */}
        <div className="flex justify-end mb-4">
          <button
            onClick={closeDetailPanel}
            className="p-1.5 hover:bg-gray-100 rounded transition"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Image */}
        {selectedItem.imageUrl ? (
          <img
            src={selectedItem.imageUrl}
            alt={selectedItem.title}
            className="w-full h-48 object-cover rounded-lg mb-5"
            onError={(e) => {
              e.currentTarget.style.display = 'none';
            }}
          />
        ) : (
          <div className="w-full h-48 bg-gray-100 rounded-lg mb-5 flex items-center justify-center text-gray-400 text-sm">
            No image available
          </div>
        )}

        {/* Title */}
        <h3 className="text-lg font-normal text-gray-900 mb-5">
          {selectedItem.title}
        </h3>

        {/* Details */}
        <div className="space-y-5">
          {/* Category */}
          {category && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Category
              </div>
              <div className="text-sm text-gray-900">
                {category.icon} {category.name}
              </div>
            </div>
          )}

          {/* Circles */}
          {itemCircles.length > 0 && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Shared With
              </div>
              <div className="flex flex-wrap gap-2">
                {itemCircles.map((circle) => (
                  <span
                    key={circle.id}
                    className="px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded-full"
                  >
                    {circle.icon} {circle.name}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Source */}
          {selectedItem.source && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Source
              </div>
              <div className="text-sm text-gray-900">
                {selectedItem.source.charAt(0).toUpperCase() + selectedItem.source.slice(1)}
                {selectedItem.sourceUrl && (
                  <a
                    href={selectedItem.sourceUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="ml-2 text-primary-600 hover:underline"
                  >
                    View original
                  </a>
                )}
              </div>
            </div>
          )}

          {/* Added */}
          <div>
            <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
              Added
            </div>
            <div className="text-sm text-gray-900">
              {formatDistanceToNow(new Date(selectedItem.addedAt), { addSuffix: true })}
            </div>
          </div>

          {/* Due Date */}
          {selectedItem.dueDate && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Due Date
              </div>
              <div className="text-sm text-gray-900">
                {format(new Date(selectedItem.dueDate), 'MMM d, yyyy')}
              </div>
            </div>
          )}

          {/* Completed */}
          {selectedItem.isCompleted && selectedItem.completedAt && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Completed
              </div>
              <div className="text-sm text-gray-900">
                {format(new Date(selectedItem.completedAt), 'MMM d, yyyy')}
              </div>
            </div>
          )}

          {/* Tags */}
          {selectedItem.tags && selectedItem.tags.length > 0 && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Tags
              </div>
              <div className="flex flex-wrap gap-2">
                {selectedItem.tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-2 py-1 text-xs bg-primary-50 text-primary-700 rounded-full"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Description */}
          {selectedItem.description && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Description
              </div>
              <div className="text-sm text-gray-700 leading-relaxed">
                {selectedItem.description}
              </div>
            </div>
          )}

          {/* Metadata */}
          {selectedItem.metadata?.rating && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Rating
              </div>
              <div className="text-sm text-gray-900">
                ⭐ {selectedItem.metadata.rating}/10
              </div>
            </div>
          )}

          {selectedItem.metadata?.location && (
            <div>
              <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                Location
              </div>
              <div className="text-sm text-gray-900">
                📍 {selectedItem.metadata.location.address}
              </div>
            </div>
          )}
        </div>

        {/* Actions */}
        <div className="mt-8 pt-6 border-t border-gray-200 space-y-2">
          <button
            onClick={handleDelete}
            className="w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 rounded-lg transition"
          >
            Delete Item
          </button>
        </div>
      </div>
    </aside>
    </>
  );
}

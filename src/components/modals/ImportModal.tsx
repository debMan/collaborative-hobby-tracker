import { useState, useMemo, useEffect } from 'react';
import { X } from 'lucide-react';
import { useStore } from '../../store';
import { DataSource } from '../../types';

export default function ImportModal() {
  const {
    isImportModalOpen,
    closeImportModal,
    createItem,
    importItem,
    categories,
    circles,
    findOrCreateCategory,
  } = useStore();

  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [suggestions, setSuggestions] = useState<{
    categoryName: string;
    categoryIcon: string;
    confidence: number;
    tags: string[];
  } | null>(null);

  // New hierarchical state
  const [selectedCircles, setSelectedCircles] = useState<string[]>(['circle-personal']); // Default to personal
  const [categoryByCircle, setCategoryByCircle] = useState<Record<string, {
    categoryId?: string; // Existing category ID
    categoryName: string; // Name (for new or existing)
    categoryIcon: string; // Icon (for new or existing)
  }>>({});

  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [newTag, setNewTag] = useState('');

  // Close modal on Esc key press
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        closeImportModal();
      }
    };

    if (isImportModalOpen) {
      document.addEventListener('keydown', handleEscape);
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
    };
  }, [isImportModalOpen, closeImportModal]);

  if (!isImportModalOpen) return null;

  const detectSource = (text: string): DataSource => {
    if (text.includes('youtube.com') || text.includes('youtu.be')) return 'youtube';
    if (text.includes('instagram.com')) return 'instagram';
    if (text.includes('twitter.com') || text.includes('x.com')) return 'twitter';
    if (text.includes('tiktok.com')) return 'tiktok';
    if (text.includes('t.me')) return 'telegram';
    if (text.includes('wikipedia.org')) return 'wikipedia';
    if (text.startsWith('http')) return 'web';
    return 'manual';
  };

  const handleImport = async () => {
    if (!input.trim()) return;

    setIsLoading(true);
    try {
      const source = detectSource(input);
      const isUrl = input.startsWith('http');

      const result = await importItem({
        source,
        url: isUrl ? input : undefined,
        text: !isUrl ? input : undefined,
      });

      if (result.success && result.suggestions) {
        setSuggestions(result.suggestions);
        setSelectedTags(result.suggestions.tags || []);

        // Initialize category for selected circles with AI suggestion
        const initialCategories: Record<string, { categoryName: string; categoryIcon: string }> = {};
        selectedCircles.forEach(circleId => {
          initialCategories[circleId] = {
            categoryName: result.suggestions.categoryName,
            categoryIcon: result.suggestions.categoryIcon,
          };
        });
        setCategoryByCircle(initialCategories);
      }
    } catch (error) {
      console.error('Import failed:', error);
      alert('Failed to import item. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddTag = () => {
    if (newTag.trim() && !selectedTags.includes(newTag.trim())) {
      setSelectedTags([...selectedTags, newTag.trim()]);
      setNewTag('');
    }
  };

  const handleRemoveTag = (tag: string) => {
    setSelectedTags(selectedTags.filter((t) => t !== tag));
  };

  const toggleCircle = (circleId: string) => {
    setSelectedCircles((prev) => {
      const newCircles = prev.includes(circleId)
        ? prev.filter((id) => id !== circleId)
        : [...prev, circleId];

      // Initialize category for newly selected circle
      if (!prev.includes(circleId) && suggestions) {
        setCategoryByCircle(prevCats => ({
          ...prevCats,
          [circleId]: {
            categoryName: suggestions.categoryName,
            categoryIcon: suggestions.categoryIcon,
          },
        }));
      }

      return newCircles;
    });
  };

  const handleCategorySelect = (circleId: string, categoryId: string, name: string, icon: string) => {
    setCategoryByCircle(prev => ({
      ...prev,
      [circleId]: { categoryId, categoryName: name, categoryIcon: icon },
    }));
  };

  const handleNewCategoryInput = (circleId: string, name: string, icon?: string) => {
    setCategoryByCircle(prev => ({
      ...prev,
      [circleId]: {
        categoryName: name,
        categoryIcon: icon || prev[circleId]?.categoryIcon || '📋',
        categoryId: undefined, // Clear existing category selection
      },
    }));
  };

  const getCategoriesForCircle = (circleId: string) => {
    return categories.filter(cat => cat.circleId === circleId);
  };

  const handleSave = async () => {
    if (!input.trim()) return;

    // Validate that all selected circles have a category
    const missingCategory = selectedCircles.some(circleId => !categoryByCircle[circleId]?.categoryName);
    if (missingCategory) {
      alert('Please select or create a category for each selected circle.');
      return;
    }

    setIsLoading(true);
    try {
      const source = detectSource(input);
      const isUrl = input.startsWith('http');

      // Create separate items for each selected circle
      // Each item gets one categoryId from that circle
      for (const circleId of selectedCircles) {
        const circleCategory = categoryByCircle[circleId];
        if (!circleCategory) continue;

        let categoryId: string;

        if (circleCategory.categoryId) {
          // Use existing category
          const category = categories.find(c => c.id === circleCategory.categoryId);
          if (!category) throw new Error('Selected category not found');
          categoryId = category.id;
        } else {
          // Create new category in this circle
          const newCategory = await findOrCreateCategory(
            circleCategory.categoryName.trim(),
            circleCategory.categoryIcon,
            circleId
          );
          categoryId = newCategory.id;
        }

        // Create a separate item for this circle
        await createItem({
          title: input.trim(),
          categoryId,
          source,
          sourceUrl: isUrl ? input : undefined,
          tags: selectedTags,
        });
      }

      // Reset form
      setInput('');
      setSuggestions(null);
      setSelectedTags([]);
      setSelectedCircles(['circle-personal']);
      setCategoryByCircle({});
      closeImportModal();
    } catch (error) {
      console.error('Failed to create item:', error);
      alert('Failed to create item. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const emojiOptions = ['🎬', '🍽️', '✈️', '🎵', '🎯', '🍕', '📚', '📋', '🎨', '🏃', '🎮', '📺'];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 dark:border-gray-700 p-4 sm:p-6">
          <h2 className="text-lg sm:text-xl font-normal text-gray-900 dark:text-gray-100">Add New Item</h2>
          <button
            onClick={closeImportModal}
            className="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition"
          >
            <X className="w-5 h-5 text-gray-500 dark:text-gray-400" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
          {/* Input */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Enter text or paste a link
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && !suggestions && handleImport()}
                placeholder="Paste YouTube link, Instagram post, or type manually..."
                className="flex-1 px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                disabled={isLoading}
              />
              {!suggestions && (
                <button
                  onClick={handleImport}
                  disabled={!input.trim() || isLoading}
                  className="px-6 py-2.5 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition"
                >
                  {isLoading ? 'Analyzing...' : 'Analyze'}
                </button>
              )}
            </div>
          </div>

          {/* AI Suggestions */}
          {suggestions && (
            <div className="p-4 bg-blue-50 dark:bg-blue-900 border border-blue-200 dark:border-blue-800 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-sm font-medium text-blue-900 dark:text-blue-100">AI Suggestion</span>
                <span className="text-xs text-blue-600 dark:text-blue-300">
                  {Math.round(suggestions.confidence * 100)}% confident
                </span>
              </div>
              <p className="text-sm text-blue-800 dark:text-blue-200">
                Suggested category: <strong>{suggestions.categoryName}</strong> {suggestions.categoryIcon}
              </p>
            </div>
          )}

          {/* Circle Selection */}
          {suggestions && (
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Select Circles
              </label>
              <div className="flex flex-wrap gap-2">
                {circles.map((circle) => (
                  <button
                    key={circle.id}
                    onClick={() => toggleCircle(circle.id)}
                    className={`px-3 py-2 border rounded-lg text-sm transition ${
                      selectedCircles.includes(circle.id)
                        ? 'border-primary-600 bg-primary-50 dark:bg-primary-900 text-primary-700 dark:text-primary-300'
                        : 'border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-500'
                    }`}
                  >
                    {circle.icon} {circle.name}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Category Selection Per Circle */}
          {suggestions && selectedCircles.map(circleId => {
            const circle = circles.find(c => c.id === circleId);
            const circleName = circle?.name || 'Unknown';
            const circleIcon = circle?.icon || '❓';
            const circleCategories = getCategoriesForCircle(circleId);
            const selectedCategory = categoryByCircle[circleId];

            return (
              <div key={circleId} className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800">
                <h3 className="text-sm font-medium text-gray-900 dark:text-gray-100 mb-3 flex items-center gap-2">
                  {circleIcon} {circleName} → Select Category
                </h3>

                {/* Existing categories */}
                {circleCategories.length > 0 && (
                  <div className="mb-3">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">Choose existing:</p>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                      {circleCategories.map((cat) => (
                        <button
                          key={cat.id}
                          onClick={() => handleCategorySelect(circleId, cat.id, cat.name, cat.icon)}
                          className={`p-2 border rounded-lg text-left transition hover:border-primary-600 dark:hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900 ${
                            selectedCategory?.categoryId === cat.id
                              ? 'border-primary-600 dark:border-primary-500 bg-primary-50 dark:bg-primary-900 text-primary-700 dark:text-primary-300'
                              : 'border-gray-300 dark:border-gray-600'
                          }`}
                        >
                          <span className="text-xl block mb-1">{cat.icon}</span>
                          <div className="text-xs text-gray-700 dark:text-gray-300 truncate">{cat.name}</div>
                        </button>
                      ))}
                    </div>
                  </div>
                )}

                {/* Create new category */}
                <div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">Or create new:</p>
                  <div className="flex gap-2">
                    <select
                      value={selectedCategory?.categoryIcon || '📋'}
                      onChange={(e) => handleNewCategoryInput(circleId, selectedCategory?.categoryName || '', e.target.value)}
                      className="w-16 px-2 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-primary-500 text-center text-xl"
                    >
                      {emojiOptions.map((emoji) => (
                        <option key={emoji} value={emoji}>
                          {emoji}
                        </option>
                      ))}
                    </select>
                    <input
                      type="text"
                      value={selectedCategory?.categoryName || ''}
                      onChange={(e) => handleNewCategoryInput(circleId, e.target.value)}
                      placeholder="Enter category name..."
                      className="flex-1 px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500"
                    />
                  </div>
                </div>
              </div>
            );
          })}

          {/* Tags */}
          {suggestions && (
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Tags
              </label>
              <div className="flex flex-wrap gap-2 mb-2">
                {selectedTags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-3 py-1 bg-primary-50 dark:bg-primary-900 text-primary-700 dark:text-primary-300 rounded-full text-sm"
                  >
                    {tag}
                    <button
                      onClick={() => handleRemoveTag(tag)}
                      className="hover:text-primary-900 dark:hover:text-primary-100"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
              <div className="flex gap-2">
                <input
                  type="text"
                  value={newTag}
                  onChange={(e) => setNewTag(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && handleAddTag()}
                  placeholder="Add a tag..."
                  className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 placeholder-gray-400 dark:placeholder-gray-500 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
                <button
                  onClick={handleAddTag}
                  className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 text-sm hover:bg-gray-50 dark:hover:bg-gray-600 transition"
                >
                  Add
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex flex-col-reverse sm:flex-row items-stretch sm:items-center justify-end gap-2 sm:gap-3 border-t border-gray-200 dark:border-gray-700 p-4 sm:p-6">
          <button
            onClick={closeImportModal}
            className="px-4 py-2.5 sm:py-2 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 active:bg-gray-200 dark:active:bg-gray-600 rounded-lg transition"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={!suggestions || selectedCircles.length === 0 || isLoading}
            className="px-6 py-2.5 sm:py-2 bg-primary-600 dark:bg-primary-700 text-white rounded-lg hover:bg-primary-700 dark:hover:bg-primary-600 active:bg-primary-800 dark:active:bg-primary-800 disabled:opacity-50 disabled:cursor-not-allowed transition"
          >
            {isLoading ? 'Saving...' : 'Add Item'}
          </button>
        </div>
      </div>
    </div>
  );
}

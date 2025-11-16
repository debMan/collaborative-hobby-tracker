import { useState, useMemo } from 'react';
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
    selectedCircles,
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
  const [categoryInput, setCategoryInput] = useState('');
  const [categoryIcon, setCategoryIcon] = useState('📋');
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [newTag, setNewTag] = useState('');
  const [selectedItemCircles, setSelectedItemCircles] = useState<string[]>([]);

  // Filter categories that match the input
  const filteredCategories = useMemo(() => {
    if (!categoryInput.trim()) return categories.slice(0, 6); // Show recent categories
    return categories
      .filter((cat) =>
        cat.name.toLowerCase().includes(categoryInput.toLowerCase())
      )
      .slice(0, 6);
  }, [categoryInput, categories]);

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
        setCategoryInput(result.suggestions.categoryName);
        setCategoryIcon(result.suggestions.categoryIcon);
        setSelectedTags(result.suggestions.tags || []);
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

  const handleCategorySelect = (name: string, icon: string) => {
    setCategoryInput(name);
    setCategoryIcon(icon);
  };

  const toggleCircle = (circleId: string) => {
    setSelectedItemCircles((prev) =>
      prev.includes(circleId)
        ? prev.filter((id) => id !== circleId)
        : [...prev, circleId]
    );
  };

  const handleSave = async () => {
    if (!input.trim() || !categoryInput.trim()) return;

    setIsLoading(true);
    try {
      const source = detectSource(input);
      const isUrl = input.startsWith('http');

      // Determine which circle to assign the category to
      // If one circle is selected in filters, use that; otherwise null (personal)
      const categoryCircleId = selectedCircles.length === 1 ? selectedCircles[0] : null;

      // Find or create the category
      const category = await findOrCreateCategory(
        categoryInput.trim(),
        categoryIcon,
        categoryCircleId
      );

      await createItem({
        title: input.trim(),
        categoryId: category.id,
        source,
        sourceUrl: isUrl ? input : undefined,
        tags: selectedTags,
        circleIds: selectedItemCircles,
      });

      // Reset form
      setInput('');
      setSuggestions(null);
      setSelectedTags([]);
      setCategoryInput('');
      setCategoryIcon('📋');
      setSelectedItemCircles([]);
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
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 p-4 sm:p-6">
          <h2 className="text-lg sm:text-xl font-normal text-gray-900">Add New Item</h2>
          <button
            onClick={closeImportModal}
            className="p-1.5 hover:bg-gray-100 rounded transition"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Content */}
        <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
          {/* Input */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Enter text or paste a link
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && !suggestions && handleImport()}
                placeholder="Paste YouTube link, Instagram post, or type manually..."
                className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
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
            <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-sm font-medium text-blue-900">AI Suggestion</span>
                <span className="text-xs text-blue-600">
                  {Math.round(suggestions.confidence * 100)}% confident
                </span>
              </div>
              <p className="text-sm text-blue-800">
                Suggested category: <strong>{suggestions.categoryName}</strong> {suggestions.categoryIcon}
              </p>
            </div>
          )}

          {/* Category Selection */}
          {suggestions && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Category
              </label>
              <div className="space-y-3">
                {/* Existing categories */}
                {categories.length > 0 && (
                  <div>
                    <p className="text-xs text-gray-500 mb-2">Select an existing category:</p>
                    <div className="grid grid-cols-2 sm:grid-cols-3 gap-2 mb-3">
                      {categories.slice(0, 6).map((cat) => (
                        <button
                          key={cat.id}
                          onClick={() => handleCategorySelect(cat.name, cat.icon)}
                          className={`p-2 border rounded-lg text-left transition hover:border-primary-600 hover:bg-primary-50 ${
                            categoryInput === cat.name
                              ? 'border-primary-600 bg-primary-50 text-primary-700'
                              : 'border-gray-300'
                          }`}
                        >
                          <div className="text-xl mb-1">{cat.icon}</div>
                          <div className="text-xs text-gray-700 truncate">{cat.name}</div>
                        </button>
                      ))}
                    </div>
                  </div>
                )}

                {/* Or create new category */}
                <div>
                  <p className="text-xs text-gray-500 mb-2">Or create a new category:</p>
                  <div className="flex gap-2">
                    <select
                      value={categoryIcon}
                      onChange={(e) => setCategoryIcon(e.target.value)}
                      className="w-16 px-2 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-center text-xl"
                    >
                      {emojiOptions.map((emoji) => (
                        <option key={emoji} value={emoji}>
                          {emoji}
                        </option>
                      ))}
                    </select>
                    <input
                      type="text"
                      value={categoryInput}
                      onChange={(e) => setCategoryInput(e.target.value)}
                      placeholder="Enter new category name..."
                      className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                    />
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Circle Selection */}
          {suggestions && circles.length > 0 && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Share with Circles (optional)
              </label>
              <div className="flex flex-wrap gap-2">
                {circles.map((circle) => (
                  <button
                    key={circle.id}
                    onClick={() => toggleCircle(circle.id)}
                    className={`px-3 py-2 border rounded-lg text-sm transition ${
                      selectedItemCircles.includes(circle.id)
                        ? 'border-primary-600 bg-primary-50 text-primary-700'
                        : 'border-gray-300 text-gray-700 hover:border-gray-400'
                    }`}
                  >
                    {circle.icon} {circle.name}
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Tags */}
          {suggestions && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tags
              </label>
              <div className="flex flex-wrap gap-2 mb-2">
                {selectedTags.map((tag) => (
                  <span
                    key={tag}
                    className="inline-flex items-center gap-1 px-3 py-1 bg-primary-50 text-primary-700 rounded-full text-sm"
                  >
                    {tag}
                    <button
                      onClick={() => handleRemoveTag(tag)}
                      className="hover:text-primary-900"
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
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
                <button
                  onClick={handleAddTag}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 transition"
                >
                  Add
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex flex-col-reverse sm:flex-row items-stretch sm:items-center justify-end gap-2 sm:gap-3 border-t border-gray-200 p-4 sm:p-6">
          <button
            onClick={closeImportModal}
            className="px-4 py-2.5 sm:py-2 text-gray-700 hover:bg-gray-100 active:bg-gray-200 rounded-lg transition"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={!suggestions || !categoryInput.trim() || isLoading}
            className="px-6 py-2.5 sm:py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 active:bg-primary-800 disabled:opacity-50 disabled:cursor-not-allowed transition"
          >
            {isLoading ? 'Saving...' : 'Add Item'}
          </button>
        </div>
      </div>
    </div>
  );
}

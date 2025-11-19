import { create } from 'zustand';
import { HobbyItem, Category, Circle, User, Tag, DataSource } from '../types';
import {
  authService,
  itemsService,
  categoriesService,
  circlesService,
  tagsService,
  importService
} from '../services/api';
import type { ImportRequest } from '../types';

interface AppStore {
  // Auth state
  user: User | null;
  isAuthenticated: boolean;
  isAuthLoading: boolean;

  // Data state
  items: HobbyItem[];
  categories: Category[];
  circles: Circle[];
  tags: Tag[];

  // UI state
  selectedItem: HobbyItem | null;
  selectedCategoryTab: string | 'all'; // 'all' or category ID
  selectedSources: DataSource[];
  selectedCircles: string[]; // Circle IDs (can select multiple)
  isDetailPanelOpen: boolean;
  isImportModalOpen: boolean;
  isMobileMenuOpen: boolean;
  isSourcesExpanded: boolean;
  isCirclesExpanded: boolean;
  isLoading: boolean;
  isDarkMode: boolean;

  // Auth actions
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, name: string) => Promise<void>;
  oauthLogin: (provider: 'google' | 'apple' | 'twitter') => Promise<void>;
  logout: () => Promise<void>;
  checkAuth: () => Promise<void>;

  // Item actions
  fetchItems: () => Promise<void>;
  getFilteredItems: () => HobbyItem[]; // Get items filtered by current filters
  createItem: (item: Partial<HobbyItem>) => Promise<void>;
  updateItem: (id: string, updates: Partial<HobbyItem>) => Promise<void>;
  deleteItem: (id: string) => Promise<void>;
  toggleItemComplete: (id: string) => Promise<void>;

  // Import actions
  importItem: (request: ImportRequest) => Promise<any>;

  // Category actions
  fetchCategories: () => Promise<void>;
  getCategoriesByFilters: () => Category[]; // Get categories filtered by selected circles
  createCategory: (category: Partial<Category>) => Promise<Category>;
  updateCategory: (id: string, updates: Partial<Category>) => Promise<void>;
  deleteCategory: (id: string) => Promise<void>;
  findOrCreateCategory: (name: string, icon: string, circleId: string | null) => Promise<Category>;

  // Circle actions
  fetchCircles: () => Promise<void>;
  createCircle: (circle: Partial<Circle>) => Promise<void>;
  updateCircle: (id: string, updates: Partial<Circle>) => Promise<void>;
  deleteCircle: (id: string) => Promise<void>;

  // Tag actions
  fetchTags: () => Promise<void>;
  suggestTags: (text: string) => Promise<string[]>;

  // UI actions
  setSelectedItem: (item: HobbyItem | null) => void;
  setSelectedCategoryTab: (tabId: string | 'all') => void;
  toggleSourceFilter: (source: DataSource) => void;
  toggleCircleFilter: (circleId: string) => void;
  setSelectedCircles: (circleIds: string[]) => void;
  toggleSourcesExpanded: () => void;
  toggleCirclesExpanded: () => void;
  openDetailPanel: (item: HobbyItem) => void;
  closeDetailPanel: () => void;
  openImportModal: () => void;
  closeImportModal: () => void;
  toggleMobileMenu: () => void;
  closeMobileMenu: () => void;
  toggleDarkMode: () => void;
}

export const useStore = create<AppStore>((set, get) => ({
  // Initial state
  user: null,
  isAuthenticated: false,
  isAuthLoading: true,
  items: [],
  categories: [],
  circles: [],
  tags: [],
  selectedItem: null,
  selectedCategoryTab: 'all',
  selectedSources: [],
  selectedCircles: [],
  isDetailPanelOpen: false,
  isImportModalOpen: false,
  isMobileMenuOpen: false,
  isSourcesExpanded: false,
  isCirclesExpanded: true,
  isLoading: false,
  isDarkMode: (() => {
    const stored = localStorage.getItem('darkMode');
    if (stored !== null) {
      return stored === 'true';
    }
    // Default to system preference if no stored preference
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  })(),

  // Auth actions
  login: async (email: string, password: string) => {
    try {
      const { user, token } = await authService.login(email, password);
      set({ user, isAuthenticated: true });
      // Load user data after login
      get().fetchItems();
      get().fetchCategories();
      get().fetchCircles();
      get().fetchTags();
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  },

  register: async (email: string, password: string, name: string) => {
    try {
      const { user, token } = await authService.register(email, password, name);
      set({ user, isAuthenticated: true });
      // Load user data after registration
      get().fetchCategories();
      get().fetchCircles();
      get().fetchTags();
    } catch (error) {
      console.error('Registration failed:', error);
      throw error;
    }
  },

  oauthLogin: async (provider: 'google' | 'apple' | 'twitter') => {
    try {
      const { user, token } = await authService.oauthLogin(provider);
      set({ user, isAuthenticated: true });
      // Load user data after OAuth login
      get().fetchItems();
      get().fetchCategories();
      get().fetchCircles();
      get().fetchTags();
    } catch (error) {
      console.error('OAuth login failed:', error);
      throw error;
    }
  },

  logout: async () => {
    try {
      await authService.logout();
      set({
        user: null,
        isAuthenticated: false,
        items: [],
        categories: [],
        circles: [],
        selectedItem: null,
        selectedCategoryTab: 'all',
        selectedSources: [],
        selectedCircles: [],
      });
    } catch (error) {
      console.error('Logout failed:', error);
    }
  },

  checkAuth: async () => {
    try {
      set({ isAuthLoading: true });
      const user = await authService.getCurrentUser();
      if (user) {
        set({ user, isAuthenticated: true });
        // Load user data
        await Promise.all([
          get().fetchItems(),
          get().fetchCategories(),
          get().fetchCircles(),
          get().fetchTags(),
        ]);
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      set({ user: null, isAuthenticated: false });
    } finally {
      set({ isAuthLoading: false });
    }
  },

  // Item actions
  fetchItems: async () => {
    try {
      set({ isLoading: true });
      const items = await itemsService.getItems();
      set({ items });
    } catch (error) {
      console.error('Failed to fetch items:', error);
    } finally {
      set({ isLoading: false });
    }
  },

  getFilteredItems: () => {
    const { items, categories, selectedCategoryTab, selectedSources, selectedCircles } = get();

    return items.filter(item => {
      // Filter by category tab
      if (selectedCategoryTab !== 'all' && item.categoryId !== selectedCategoryTab) {
        return false;
      }

      // Filter by sources
      if (selectedSources.length > 0 && !selectedSources.includes(item.source)) {
        return false;
      }

      // Filter by circles - check if item's category is in selected circles
      if (selectedCircles.length > 0) {
        const category = categories.find(c => c.id === item.categoryId);
        if (!category) return false;

        // Check if category's circle matches selected circles
        if (!selectedCircles.includes(category.circleId)) return false;
      }

      return true;
    });
  },

  createItem: async (item: Partial<HobbyItem>) => {
    try {
      const newItem = await itemsService.createItem(item);
      set({ items: [...get().items, newItem] });
    } catch (error) {
      console.error('Failed to create item:', error);
      throw error;
    }
  },

  updateItem: async (id: string, updates: Partial<HobbyItem>) => {
    try {
      const updatedItem = await itemsService.updateItem(id, updates);
      set({
        items: get().items.map(item => item.id === id ? updatedItem : item),
        selectedItem: get().selectedItem?.id === id ? updatedItem : get().selectedItem,
      });
    } catch (error) {
      console.error('Failed to update item:', error);
      throw error;
    }
  },

  deleteItem: async (id: string) => {
    try {
      await itemsService.deleteItem(id);
      set({
        items: get().items.filter(item => item.id !== id),
        selectedItem: get().selectedItem?.id === id ? null : get().selectedItem,
        isDetailPanelOpen: get().selectedItem?.id === id ? false : get().isDetailPanelOpen,
      });
    } catch (error) {
      console.error('Failed to delete item:', error);
      throw error;
    }
  },

  toggleItemComplete: async (id: string) => {
    try {
      const updatedItem = await itemsService.toggleComplete(id);
      set({
        items: get().items.map(item => item.id === id ? updatedItem : item),
        selectedItem: get().selectedItem?.id === id ? updatedItem : get().selectedItem,
      });
    } catch (error) {
      console.error('Failed to toggle item:', error);
      throw error;
    }
  },

  // Import actions
  importItem: async (request: ImportRequest) => {
    try {
      const result = await importService.importItem(request);
      return result;
    } catch (error) {
      console.error('Failed to import item:', error);
      throw error;
    }
  },

  // Category actions
  fetchCategories: async () => {
    try {
      const categories = await categoriesService.getCategories();
      set({ categories });
    } catch (error) {
      console.error('Failed to fetch categories:', error);
    }
  },

  getCategoriesByFilters: () => {
    const { categories, selectedCircles, items } = get();

    // If no circles selected, show all categories
    if (selectedCircles.length === 0) {
      return categories;
    }

    // Get unique category IDs from items that match the selected circles
    const categoryIds = new Set<string>();
    items.forEach(item => {
      const category = categories.find(c => c.id === item.categoryId);
      if (!category) return;

      // Only include if category belongs to one of the selected circles
      if (selectedCircles.includes(category.circleId)) {
        categoryIds.add(item.categoryId);
      }
    });

    // Return categories that have items in the selected circles
    return categories.filter(category => categoryIds.has(category.id));
  },

  createCategory: async (category: Partial<Category>) => {
    try {
      const newCategory = await categoriesService.createCategory(category);
      set({ categories: [...get().categories, newCategory] });
      return newCategory;
    } catch (error) {
      console.error('Failed to create category:', error);
      throw error;
    }
  },

  updateCategory: async (id: string, updates: Partial<Category>) => {
    try {
      const updatedCategory = await categoriesService.updateCategory(id, updates);
      set({ categories: get().categories.map(cat => cat.id === id ? updatedCategory : cat) });
    } catch (error) {
      console.error('Failed to update category:', error);
      throw error;
    }
  },

  deleteCategory: async (id: string) => {
    try {
      await categoriesService.deleteCategory(id);
      set({ categories: get().categories.filter(cat => cat.id !== id) });
    } catch (error) {
      console.error('Failed to delete category:', error);
      throw error;
    }
  },

  findOrCreateCategory: async (name: string, icon: string, circleId: string | null) => {
    const { categories } = get();

    // Try to find existing category with same name and circle
    const existing = categories.find(cat =>
      cat.name.toLowerCase() === name.toLowerCase() && cat.circleId === circleId
    );

    if (existing) {
      return existing;
    }

    // Create new category
    const newCategory = await get().createCategory({
      name,
      icon,
      circleId,
      itemCount: 0,
    });

    return newCategory;
  },

  // Circle actions
  fetchCircles: async () => {
    try {
      const circles = await circlesService.getCircles();
      set({ circles });
    } catch (error) {
      console.error('Failed to fetch circles:', error);
    }
  },

  createCircle: async (circle: Partial<Circle>) => {
    try {
      const newCircle = await circlesService.createCircle(circle);
      set({ circles: [...get().circles, newCircle] });
    } catch (error) {
      console.error('Failed to create circle:', error);
      throw error;
    }
  },

  updateCircle: async (id: string, updates: Partial<Circle>) => {
    try {
      const updatedCircle = await circlesService.updateCircle(id, updates);
      set({ circles: get().circles.map(circle => circle.id === id ? updatedCircle : circle) });
    } catch (error) {
      console.error('Failed to update circle:', error);
      throw error;
    }
  },

  deleteCircle: async (id: string) => {
    try {
      await circlesService.deleteCircle(id);
      set({ circles: get().circles.filter(circle => circle.id !== id) });
    } catch (error) {
      console.error('Failed to delete circle:', error);
      throw error;
    }
  },

  // Tag actions
  fetchTags: async () => {
    try {
      const tags = await tagsService.getTags();
      set({ tags });
    } catch (error) {
      console.error('Failed to fetch tags:', error);
    }
  },

  suggestTags: async (text: string) => {
    try {
      return await tagsService.suggestTags(text);
    } catch (error) {
      console.error('Failed to suggest tags:', error);
      return [];
    }
  },

  // UI actions
  setSelectedItem: (item) => set({ selectedItem: item }),

  setSelectedCategoryTab: (tabId) => set({ selectedCategoryTab: tabId }),

  toggleSourceFilter: (source) => {
    const { selectedSources } = get();
    set({
      selectedSources: selectedSources.includes(source)
        ? selectedSources.filter(s => s !== source)
        : [...selectedSources, source]
    });
  },

  toggleCircleFilter: (circleId) => {
    const { selectedCircles } = get();
    set({
      selectedCircles: selectedCircles.includes(circleId)
        ? selectedCircles.filter(id => id !== circleId)
        : [...selectedCircles, circleId]
    });
  },

  setSelectedCircles: (circleIds) => set({ selectedCircles: circleIds }),

  toggleSourcesExpanded: () => set((state) => ({ isSourcesExpanded: !state.isSourcesExpanded })),

  toggleCirclesExpanded: () => set((state) => ({ isCirclesExpanded: !state.isCirclesExpanded })),

  openDetailPanel: (item) => set({ selectedItem: item, isDetailPanelOpen: true }),

  closeDetailPanel: () => set({ selectedItem: null, isDetailPanelOpen: false }),

  openImportModal: () => set({ isImportModalOpen: true }),

  closeImportModal: () => set({ isImportModalOpen: false }),

  toggleMobileMenu: () => set((state) => ({ isMobileMenuOpen: !state.isMobileMenuOpen })),

  closeMobileMenu: () => set({ isMobileMenuOpen: false }),

  toggleDarkMode: () => {
    const newDarkMode = !get().isDarkMode;
    set({ isDarkMode: newDarkMode });
    localStorage.setItem('darkMode', String(newDarkMode));

    // Apply dark mode to document root
    if (newDarkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  },
}));

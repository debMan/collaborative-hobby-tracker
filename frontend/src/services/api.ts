import { HobbyItem, Category, Circle, User, ImportRequest, ImportResult, Tag } from '../types';
import { mockUser, mockItems, mockCategories, mockCircles, mockTags } from '../utils/mockData';

// Backend API base URL - configure this in production
// @ts-expect-error - Will be used when backend is ready
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

// Helper function for API calls (to be used when backend is ready)
// @ts-expect-error - Will be used when backend is ready
async function apiCall<T>(_endpoint: string, _options?: RequestInit): Promise<T> {
  /* BACKEND CONNECTION - Uncomment when backend is ready
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${getAuthToken()}`,
      ...options?.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API Error: ${response.statusText}`);
  }

  return response.json();
  */

  // Mock implementation - remove when backend is ready
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({} as T);
    }, 300); // Simulate network delay
  });
}

// ============================================================================
// AUTH SERVICE
// ============================================================================

export const authService = {
  // Get current user
  async getCurrentUser(): Promise<User | null> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    try {
      return await apiCall<User>('/auth/me');
    } catch (error) {
      return null;
    }
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const token = localStorage.getItem('auth_token');
        resolve(token ? mockUser : null);
      }, 300);
    });
  },

  // Login with email/password
  // @ts-expect-error - Parameters will be used when backend is ready
  async login(email: string, password: string): Promise<{ user: User; token: string }> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<{ user: User; token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const token = 'mock-jwt-token-' + Date.now();
        localStorage.setItem('auth_token', token);
        resolve({ user: mockUser, token });
      }, 500);
    });
  },

  // Register new user
  // @ts-expect-error - Parameters will be used when backend is ready
  async register(email: string, password: string, name: string): Promise<{ user: User; token: string }> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<{ user: User; token: string }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const token = 'mock-jwt-token-' + Date.now();
        localStorage.setItem('auth_token', token);
        const newUser = { ...mockUser, email, name, id: 'user-' + Date.now() };
        resolve({ user: newUser, token });
      }, 500);
    });
  },

  // OAuth login (Google, Apple, Twitter)
  async oauthLogin(provider: 'google' | 'apple' | 'twitter'): Promise<{ user: User; token: string }> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    // This would typically redirect to OAuth provider
    window.location.href = `${API_BASE_URL}/auth/oauth/${provider}`;
    return Promise.reject('Redirecting to OAuth provider');
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const token = 'mock-oauth-token-' + Date.now();
        localStorage.setItem('auth_token', token);
        resolve({ user: { ...mockUser, provider }, token });
      }, 500);
    });
  },

  // Logout
  async logout(): Promise<void> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    await apiCall('/auth/logout', { method: 'POST' });
    */

    // Mock implementation
    localStorage.removeItem('auth_token');
  },
};

// ============================================================================
// ITEMS SERVICE
// ============================================================================

export const itemsService = {
  // Get all items
  async getItems(): Promise<HobbyItem[]> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<HobbyItem[]>('/items');
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const items = JSON.parse(localStorage.getItem('hobby_items') || JSON.stringify(mockItems));
        resolve(items);
      }, 300);
    });
  },

  // Create new item
  async createItem(item: Partial<HobbyItem>): Promise<HobbyItem> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<HobbyItem>('/items', {
      method: 'POST',
      body: JSON.stringify(item),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const items = JSON.parse(localStorage.getItem('hobby_items') || JSON.stringify(mockItems));
        const newItem: HobbyItem = {
          id: 'item-' + Date.now(),
          title: item.title || 'Untitled',
          categoryId: item.categoryId || '',
          isCompleted: false,
          source: item.source || 'manual',
          addedBy: 'user-1',
          addedAt: new Date(),
          tags: item.tags || [],
          ...item,
        } as HobbyItem;
        items.push(newItem);
        localStorage.setItem('hobby_items', JSON.stringify(items));
        resolve(newItem);
      }, 300);
    });
  },

  // Update item
  async updateItem(id: string, updates: Partial<HobbyItem>): Promise<HobbyItem> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<HobbyItem>(`/items/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const items = JSON.parse(localStorage.getItem('hobby_items') || JSON.stringify(mockItems));
        const index = items.findIndex((item: HobbyItem) => item.id === id);
        if (index !== -1) {
          items[index] = { ...items[index], ...updates };
          localStorage.setItem('hobby_items', JSON.stringify(items));
          resolve(items[index]);
        }
      }, 300);
    });
  },

  // Delete item
  async deleteItem(id: string): Promise<void> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    await apiCall(`/items/${id}`, { method: 'DELETE' });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const items = JSON.parse(localStorage.getItem('hobby_items') || JSON.stringify(mockItems));
        const filtered = items.filter((item: HobbyItem) => item.id !== id);
        localStorage.setItem('hobby_items', JSON.stringify(filtered));
        resolve();
      }, 300);
    });
  },

  // Toggle item completion
  async toggleComplete(id: string): Promise<HobbyItem> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<HobbyItem>(`/items/${id}/toggle`, { method: 'POST' });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const items = JSON.parse(localStorage.getItem('hobby_items') || JSON.stringify(mockItems));
        const index = items.findIndex((item: HobbyItem) => item.id === id);
        if (index !== -1) {
          items[index].isCompleted = !items[index].isCompleted;
          items[index].completedAt = items[index].isCompleted ? new Date() : undefined;
          localStorage.setItem('hobby_items', JSON.stringify(items));
          resolve(items[index]);
        }
      }, 300);
    });
  },
};

// ============================================================================
// IMPORT SERVICE
// ============================================================================

export const importService = {
  // Import item from URL or text with AI categorization
  async importItem(request: ImportRequest): Promise<ImportResult> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<ImportResult>('/import', {
      method: 'POST',
      body: JSON.stringify(request),
    });
    */

    // Mock implementation with simulated AI categorization
    return new Promise((resolve) => {
      setTimeout(() => {
        // Simulate AI category detection
        let categoryName = request.suggestedCategory || 'Other';
        let categoryIcon = '📋';
        let confidence = 0.85;
        const tags: string[] = request.suggestedTags || [];

        // Simple keyword-based category detection (mock AI)
        const text = (request.text || request.url || '').toLowerCase();
        if (text.includes('movie') || text.includes('film') || text.includes('imdb')) {
          categoryName = 'Movies';
          categoryIcon = '🎬';
          confidence = 0.95;
          tags.push('To Watch');
        } else if (text.includes('restaurant') || text.includes('food') || text.includes('cafe')) {
          categoryName = 'Restaurants';
          categoryIcon = '🍽️';
          confidence = 0.92;
        } else if (text.includes('travel') || text.includes('visit') || text.includes('destination')) {
          categoryName = 'Travel';
          categoryIcon = '✈️';
          confidence = 0.88;
        } else if (text.includes('music') || text.includes('album') || text.includes('song')) {
          categoryName = 'Music';
          categoryIcon = '🎵';
          confidence = 0.90;
        } else if (text.includes('book') || text.includes('read') || text.includes('novel')) {
          categoryName = 'Books';
          categoryIcon = '📚';
          confidence = 0.93;
        } else if (text.includes('activity') || text.includes('learn') || text.includes('do')) {
          categoryName = 'Activities';
          categoryIcon = '🎯';
          confidence = 0.87;
        }

        resolve({
          success: true,
          item: {
            title: request.text || 'Imported Item',
            source: request.source,
            sourceUrl: request.url,
          },
          suggestions: {
            categoryName,
            categoryIcon,
            confidence,
            tags,
          },
        });
      }, 800); // Longer delay to simulate AI processing
    });
  },
};

// ============================================================================
// CATEGORIES SERVICE
// ============================================================================

export const categoriesService = {
  // Get all categories
  async getCategories(): Promise<Category[]> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Category[]>('/categories');
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const categories = JSON.parse(localStorage.getItem('hobby_categories') || JSON.stringify(mockCategories));
        resolve(categories);
      }, 300);
    });
  },

  // Create new category
  async createCategory(category: Partial<Category>): Promise<Category> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Category>('/categories', {
      method: 'POST',
      body: JSON.stringify(category),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const categories = JSON.parse(localStorage.getItem('hobby_categories') || JSON.stringify(mockCategories));
        const newCategory: Category = {
          id: 'cat-' + Date.now(),
          name: category.name || 'Untitled Category',
          icon: category.icon || '📋',
          circleId: category.circleId || 'circle-personal',
          ownerId: 'user-1',
          itemCount: 0,
          createdAt: new Date(),
          updatedAt: new Date(),
          ...category,
        } as Category;
        categories.push(newCategory);
        localStorage.setItem('hobby_categories', JSON.stringify(categories));
        resolve(newCategory);
      }, 300);
    });
  },

  // Update category
  async updateCategory(id: string, updates: Partial<Category>): Promise<Category> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Category>(`/categories/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const categories = JSON.parse(localStorage.getItem('hobby_categories') || JSON.stringify(mockCategories));
        const index = categories.findIndex((cat: Category) => cat.id === id);
        if (index !== -1) {
          categories[index] = { ...categories[index], ...updates, updatedAt: new Date() };
          localStorage.setItem('hobby_categories', JSON.stringify(categories));
          resolve(categories[index]);
        }
      }, 300);
    });
  },

  // Delete category
  async deleteCategory(id: string): Promise<void> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    await apiCall(`/categories/${id}`, { method: 'DELETE' });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const categories = JSON.parse(localStorage.getItem('hobby_categories') || JSON.stringify(mockCategories));
        const filtered = categories.filter((cat: Category) => cat.id !== id);
        localStorage.setItem('hobby_categories', JSON.stringify(filtered));
        resolve();
      }, 300);
    });
  },
};

// ============================================================================
// CIRCLES SERVICE
// ============================================================================

export const circlesService = {
  // Get all circles
  async getCircles(): Promise<Circle[]> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Circle[]>('/circles');
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const circles = JSON.parse(localStorage.getItem('hobby_circles') || JSON.stringify(mockCircles));
        resolve(circles);
      }, 300);
    });
  },

  // Create new circle
  async createCircle(circle: Partial<Circle>): Promise<Circle> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Circle>('/circles', {
      method: 'POST',
      body: JSON.stringify(circle),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const circles = JSON.parse(localStorage.getItem('hobby_circles') || JSON.stringify(mockCircles));
        const newCircle: Circle = {
          id: 'circle-' + Date.now(),
          name: circle.name || 'Untitled Circle',
          icon: circle.icon || '👥',
          ownerId: 'user-1',
          createdAt: new Date(),
          members: [
            { userId: 'user-1', accessLevel: 'admin', joinedAt: new Date() },
          ],
          ...circle,
        } as Circle;
        circles.push(newCircle);
        localStorage.setItem('hobby_circles', JSON.stringify(circles));
        resolve(newCircle);
      }, 300);
    });
  },

  // Update circle
  async updateCircle(id: string, updates: Partial<Circle>): Promise<Circle> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Circle>(`/circles/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(updates),
    });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const circles = JSON.parse(localStorage.getItem('hobby_circles') || JSON.stringify(mockCircles));
        const index = circles.findIndex((circle: Circle) => circle.id === id);
        if (index !== -1) {
          circles[index] = { ...circles[index], ...updates };
          localStorage.setItem('hobby_circles', JSON.stringify(circles));
          resolve(circles[index]);
        }
      }, 300);
    });
  },

  // Delete circle
  async deleteCircle(id: string): Promise<void> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    await apiCall(`/circles/${id}`, { method: 'DELETE' });
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const circles = JSON.parse(localStorage.getItem('hobby_circles') || JSON.stringify(mockCircles));
        const filtered = circles.filter((circle: Circle) => circle.id !== id);
        localStorage.setItem('hobby_circles', JSON.stringify(filtered));
        resolve();
      }, 300);
    });
  },
};

// ============================================================================
// TAGS SERVICE
// ============================================================================

export const tagsService = {
  // Get all tags
  async getTags(): Promise<Tag[]> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<Tag[]>('/tags');
    */

    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(mockTags);
      }, 300);
    });
  },

  // Suggest tags based on user history and AI (mock AI learning)
  async suggestTags(text: string): Promise<string[]> {
    /* BACKEND CONNECTION - Uncomment when backend is ready
    return await apiCall<string[]>('/tags/suggest', {
      method: 'POST',
      body: JSON.stringify({ text }),
    });
    */

    // Mock implementation with simple keyword matching
    return new Promise((resolve) => {
      setTimeout(() => {
        const suggestions: string[] = [];
        const lowerText = text.toLowerCase();
        
        mockTags.forEach(tag => {
          if (lowerText.includes(tag.name.toLowerCase())) {
            suggestions.push(tag.name);
          }
        });

        resolve(suggestions.slice(0, 5)); // Return top 5 suggestions
      }, 400);
    });
  },
};

// Helper function to get auth token (to be implemented)
// @ts-expect-error - Will be used when backend is ready
function getAuthToken(): string {
  return localStorage.getItem('auth_token') || '';
}

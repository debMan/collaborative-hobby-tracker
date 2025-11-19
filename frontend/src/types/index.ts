// Type definitions for the Hobby Tracker application

// Note: ItemCategory kept for backward compatibility during transition
// Categories are now dynamic and user-created
export type ItemCategory =
    | 'movie'
    | 'restaurant'
    | 'travel'
    | 'music'
    | 'activity'
    | 'food'
    | 'book'
    | 'other';

export type DataSource =
    | 'instagram'
    | 'youtube'
    | 'twitter'
    | 'tiktok'
    | 'telegram'
    | 'web'
    | 'manual'
    | 'wikipedia';

export type AccessLevel = 'private' | 'view' | 'edit' | 'admin';

export interface User {
    id: string;
    email: string;
    name: string;
    avatarUrl?: string;
    provider?: 'google' | 'apple' | 'twitter' | 'email';
    createdAt: Date;
}

export interface Tag {
    id: string;
    name: string;
    color?: string;
    usageCount: number;
}

export interface Category {
    id: string;
    name: string;
    icon: string;
    circleId: string; // Circle ID (including 'circle-personal' for personal items)
    ownerId: string;
    itemCount: number;
    createdAt: Date;
    updatedAt: Date;
}

export interface HobbyItem {
    id: string;
    title: string;
    description?: string;
    categoryId: string; // Single category ID - item belongs to one category in one circle
    categoryConfidence?: number; // AI confidence score (0-1)
    isCompleted: boolean;
    imageUrl?: string;
    source: DataSource;
    sourceUrl?: string;
    addedBy: string; // User ID
    addedAt: Date;
    completedAt?: Date;
    dueDate?: Date;
    tags: string[];
    metadata?: {
        imdbId?: string;
        rating?: number;
        location?: {
            lat: number;
            lng: number;
            address: string;
        };
        externalData?: Record<string, any>;
    };
}

export interface List {
    id: string;
    name: string;
    icon: string;
    category?: ItemCategory;
    description?: string;
    ownerId: string;
    createdAt: Date;
    updatedAt: Date;
    itemCount: number;
    sharedWith: {
        userId: string;
        accessLevel: AccessLevel;
    }[];
}

export interface Circle {
    id: string;
    name: string;
    icon: string;
    description?: string;
    ownerId: string;
    createdAt: Date;
    members: {
        userId: string;
        accessLevel: AccessLevel;
        joinedAt: Date;
    }[];
}

export interface ImportRequest {
    source: DataSource;
    url?: string;
    text?: string;
    suggestedCategory?: string; // Category name (can be new or existing)
    suggestedTags?: string[];
}

export interface ImportResult {
    success: boolean;
    item?: Partial<HobbyItem>;
    suggestions: {
        categoryName: string; // Suggested category name
        categoryIcon: string; // Suggested icon for category
        confidence: number;
        tags: string[];
        metadata?: any;
    };
    error?: string;
}

export interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
}

export interface AppState {
    items: HobbyItem[];
    categories: Category[];
    circles: Circle[];
    tags: Tag[];
    selectedItem: HobbyItem | null;
    selectedCategoryTab: string | 'all'; // 'all' or category ID
    selectedSources: DataSource[];
    selectedCircles: string[]; // Circle IDs
    isDetailPanelOpen: boolean;
    isSourcesExpanded: boolean;
    isCirclesExpanded: boolean;
}

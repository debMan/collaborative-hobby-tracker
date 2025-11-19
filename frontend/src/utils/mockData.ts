import { HobbyItem, Category, Circle, User, Tag } from '../types';

export const mockUser: User = {
  id: 'user-1',
  email: 'john.doe@example.com',
  name: 'John Doe',
  avatarUrl: undefined,
  provider: 'google',
  createdAt: new Date('2024-01-01'),
};

export const mockTags: Tag[] = [
  { id: 'tag-1', name: 'Wes Anderson', color: '#1a73e8', usageCount: 3 },
  { id: 'tag-2', name: 'Japanese', color: '#34a853', usageCount: 5 },
  { id: 'tag-3', name: 'Ramen', color: '#fbbc04', usageCount: 2 },
  { id: 'tag-4', name: 'Europe', color: '#ea4335', usageCount: 7 },
  { id: 'tag-5', name: 'Summer 2026', color: '#9334e6', usageCount: 4 },
  { id: 'tag-6', name: 'Cooking', color: '#ff6d00', usageCount: 3 },
  { id: 'tag-7', name: 'R&B', color: '#00897b', usageCount: 6 },
  { id: 'tag-8', name: 'Paris', color: '#c2185b', usageCount: 2 },
  { id: 'tag-9', name: 'Art', color: '#5e35b1', usageCount: 8 },
];

export const mockCategories: Category[] = [
  {
    id: 'cat-1',
    name: 'Movies',
    icon: '🎬',
    circleId: 'circle-personal', // Personal category
    ownerId: 'user-1',
    itemCount: 2,
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date(),
  },
  {
    id: 'cat-2',
    name: 'Restaurants',
    icon: '🍽️',
    circleId: 'circle-1', // Partner circle
    ownerId: 'user-1',
    itemCount: 1,
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date(),
  },
  {
    id: 'cat-3',
    name: 'Travel',
    icon: '✈️',
    circleId: 'circle-1', // Partner circle
    ownerId: 'user-1',
    itemCount: 2,
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date(),
  },
  {
    id: 'cat-4',
    name: 'Music',
    icon: '🎵',
    circleId: 'circle-personal', // Personal category
    ownerId: 'user-1',
    itemCount: 1,
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date(),
  },
  {
    id: 'cat-5',
    name: 'Activities',
    icon: '🎯',
    circleId: 'circle-personal', // Personal category
    ownerId: 'user-1',
    itemCount: 2,
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date(),
  },
  {
    id: 'cat-6',
    name: 'Movies',
    icon: '🎬',
    circleId: 'circle-1', // Partner circle category
    ownerId: 'user-1',
    itemCount: 1,
    createdAt: new Date('2024-01-01'),
    updatedAt: new Date(),
  },

];

export const mockItems: HobbyItem[] = [
  {
    id: 'item-1',
    title: 'The Grand Budapest Hotel',
    description: 'The adventures of Gustave H, a legendary concierge at a famous European hotel between the wars.',
    categoryId: 'cat-1', // Personal Movies
    categoryConfidence: 0.98,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=400',
    source: 'youtube',
    sourceUrl: 'https://youtube.com/watch?v=example',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
    tags: ['Wes Anderson', 'Comedy', '2014'],
    metadata: {
      imdbId: 'tt2278388',
      rating: 8.1,
    },
  },
  {
    id: 'item-2',
    title: 'Ramen Tatsu-Ya - Austin, TX',
    description: 'Authentic Japanese ramen restaurant in Austin, Texas. Famous for their tonkotsu broth.',
    categoryId: 'cat-2', // Partner Restaurants
    categoryConfidence: 0.95,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=400',
    source: 'instagram',
    sourceUrl: 'https://instagram.com/p/example',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000), // 1 day ago
    tags: ['Japanese', 'Ramen', 'Austin'],
    metadata: {
      location: {
        lat: 30.2672,
        lng: -97.7431,
        address: '1600 E 6th St, Austin, TX 78702',
      },
      rating: 4.5,
    },
  },
  {
    id: 'item-3',
    title: 'Santorini, Greece',
    description: 'Beautiful Greek island known for stunning sunsets, white buildings, and blue domed churches.',
    categoryId: 'cat-3', // Partner Travel
    categoryConfidence: 0.99,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1613395877344-13d4a8e0d49e?w=400',
    source: 'web',
    sourceUrl: 'https://example.com/santorini',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000), // 3 days ago
    dueDate: new Date('2026-07-15'),
    tags: ['Europe', 'Summer 2026', 'Greece'],
    metadata: {
      location: {
        lat: 36.3932,
        lng: 25.4615,
        address: 'Santorini, Greece',
      },
    },
  },
  {
    id: 'item-4',
    title: 'Learn to make sourdough bread',
    description: 'Master the art of sourdough baking from starter to final loaf.',
    categoryId: 'cat-5', // Personal Activities
    categoryConfidence: 0.92,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400',
    source: 'tiktok',
    sourceUrl: 'https://tiktok.com/@example',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000), // 5 days ago
    tags: ['Cooking', 'Baking', 'DIY'],
  },
  {
    id: 'item-5',
    title: 'Blonde - Frank Ocean',
    description: 'Second studio album by Frank Ocean, critically acclaimed R&B masterpiece.',
    categoryId: 'cat-4', // Personal Music
    categoryConfidence: 0.97,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=400',
    source: 'twitter',
    sourceUrl: 'https://twitter.com/example',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 6 * 24 * 60 * 60 * 1000), // 6 days ago
    tags: ['R&B', 'Frank Ocean', '2016'],
    metadata: {
      externalData: {
        spotifyId: '3mH6qwIy9crq0I9YQbOuDf',
        releaseYear: 2016,
      },
    },
  },
  {
    id: 'item-6',
    title: 'Visit the Louvre Museum',
    description: 'Explore one of the world\'s largest and most visited museums in Paris.',
    categoryId: 'cat-5', // Personal Activities (shared with Friends but same category)
    categoryConfidence: 0.94,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1499856871958-5b9627545d1a?w=400',
    source: 'wikipedia',
    sourceUrl: 'https://en.wikipedia.org/wiki/Louvre',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 14 * 24 * 60 * 60 * 1000), // 2 weeks ago
    tags: ['Paris', 'Art', 'Museum'],
    metadata: {
      location: {
        lat: 48.8606,
        lng: 2.3376,
        address: 'Rue de Rivoli, 75001 Paris, France',
      },
    },
  },
  {
    id: 'item-7',
    title: 'Parasite (2019)',
    description: 'South Korean thriller film directed by Bong Joon-ho.',
    categoryId: 'cat-1', // Personal Movies
    categoryConfidence: 0.99,
    isCompleted: true,
    completedAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
    imageUrl: 'https://images.unsplash.com/photo-1536440136628-849c177e76a1?w=400',
    source: 'youtube',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000), // 1 month ago
    tags: ['Korean Cinema', 'Thriller', 'Oscar Winner'],
    metadata: {
      imdbId: 'tt6751668',
      rating: 8.6,
    },
  },
  {
    id: 'item-8',
    title: 'Kyoto Temple Tour',
    description: 'Visit the historic temples of Kyoto including Kinkaku-ji and Fushimi Inari.',
    categoryId: 'cat-3', // Partner Travel
    categoryConfidence: 0.96,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e?w=400',
    source: 'instagram',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000),
    tags: ['Japan', 'Temple', 'Culture'],
  },
  {
    id: 'item-9',
    title: 'The Pianist',
    description: 'A biographical war drama film about a Polish Jewish pianist.',
    categoryId: 'cat-6', // Partner Movies
    categoryConfidence: 0.98,
    isCompleted: false,
    imageUrl: 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=400',
    source: 'youtube',
    sourceUrl: 'https://youtube.com/watch?v=example',
    addedBy: 'user-1',
    addedAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000), // 2 days ago
    tags: ['Drama', 'War', 'Biography'],
    metadata: {
      imdbId: 'tt0253474',
      rating: 8.5,
    },
  }
];

export const mockCircles: Circle[] = [
  {
    id: 'circle-personal',
    name: 'Personal',
    icon: '👤',
    description: 'Personal items',
    ownerId: 'user-1',
    createdAt: new Date('2024-01-01'),
    members: [
      { userId: 'user-1', accessLevel: 'admin', joinedAt: new Date('2024-01-01') },
    ],
  },
  {
    id: 'circle-1',
    name: 'Partner',
    icon: '💑',
    description: 'Shared with my partner',
    ownerId: 'user-1',
    createdAt: new Date('2024-01-01'),
    members: [
      { userId: 'user-1', accessLevel: 'admin', joinedAt: new Date('2024-01-01') },
      { userId: 'user-2', accessLevel: 'edit', joinedAt: new Date('2024-01-15') },
    ],
  },
  {
    id: 'circle-2',
    name: 'Friends',
    icon: '👥',
    description: 'Shared with close friends',
    ownerId: 'user-1',
    createdAt: new Date('2024-01-01'),
    members: [
      { userId: 'user-1', accessLevel: 'admin', joinedAt: new Date('2024-01-01') },
      { userId: 'user-3', accessLevel: 'view', joinedAt: new Date('2024-02-01') },
      { userId: 'user-4', accessLevel: 'view', joinedAt: new Date('2024-02-15') },
    ],
  },
  {
    id: 'circle-3',
    name: 'Family',
    icon: '👨‍👩‍👧',
    description: 'Family activities and plans',
    ownerId: 'user-1',
    createdAt: new Date('2024-01-01'),
    members: [
      { userId: 'user-1', accessLevel: 'admin', joinedAt: new Date('2024-01-01') },
    ],
  },
];

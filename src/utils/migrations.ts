// Data migration utilities for handling schema changes in localStorage

import { Category, Circle } from '../types';
import { mockCircles } from './mockData';

/**
 * Migrate old localStorage data to new schema
 * This runs automatically when the app loads
 */
export function migrateLocalStorageData() {
  migrateCategoriesToPersonalCircle();
  ensurePersonalCircleExists();
}

/**
 * Migrate categories with null circleId to use 'circle-personal'
 */
function migrateCategoriesToPersonalCircle() {
  try {
    const categoriesJson = localStorage.getItem('hobby_categories');
    if (!categoriesJson) return;

    const categories: Category[] = JSON.parse(categoriesJson);
    let migrated = false;

    const updatedCategories = categories.map(cat => {
      if (cat.circleId === null || cat.circleId === undefined) {
        migrated = true;
        return { ...cat, circleId: 'circle-personal' };
      }
      return cat;
    });

    if (migrated) {
      localStorage.setItem('hobby_categories', JSON.stringify(updatedCategories));
      console.log('✅ Migrated categories to use circle-personal');
    }
  } catch (error) {
    console.error('Failed to migrate categories:', error);
  }
}

/**
 * Ensure Personal circle exists in circles data
 */
function ensurePersonalCircleExists() {
  try {
    const circlesJson = localStorage.getItem('hobby_circles');
    if (!circlesJson) return;

    const circles: Circle[] = JSON.parse(circlesJson);

    // Check if Personal circle already exists
    const hasPersonalCircle = circles.some(c => c.id === 'circle-personal');

    if (!hasPersonalCircle) {
      // Add Personal circle at the beginning
      const personalCircle = mockCircles.find(c => c.id === 'circle-personal');
      if (personalCircle) {
        circles.unshift(personalCircle);
        localStorage.setItem('hobby_circles', JSON.stringify(circles));
        console.log('✅ Added Personal circle to circles data');
      }
    }
  } catch (error) {
    console.error('Failed to ensure Personal circle:', error);
  }
}

/**
 * Clear all localStorage data (useful for debugging)
 */
export function clearAllData() {
  localStorage.removeItem('hobby_items');
  localStorage.removeItem('hobby_categories');
  localStorage.removeItem('hobby_circles');
  localStorage.removeItem('hobby_tags');
  console.log('🗑️ Cleared all hobby tracker data from localStorage');
}

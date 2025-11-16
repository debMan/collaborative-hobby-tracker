## Instructions

I want to develop a cross-platform web application that tracks and organizes hobbies.

## Features

- The app should be able to decide which kind of hobby the user is adding to the category. It can be a travel destination, a restaurant to visit, food to taste, a movie to watch, an activity to do, music to listen to, etc
- Categorizing the hobby should be decided with the assistance of an AI agent and suggested to the user when importing. Also, the user can change and override the suggested category to add it to their desired category list.
- The incoming data can be imported from Instagram, YouTube, X (former Twitter) post, Wikipedia, TikTok, Telegram channel post, a web link, plain text, etc
- An image of that item should be fetched and attached to the item from the source (Nice to have)
- Each item should have a check box beside it to be tracked, enriched with metadata like date added, date done, imported by which user, due date for planning, integrating with personal calendar (on phone's calendar, or third-party accounts like Google Calendar or Apple Calendar), data source (Instagram, YouTube, etc), tags
- The AI assistant should learn and decide which tags are used by user, and tag the next addind items based on the user preferences
- It should be able to fetch details for each item from the official source. For example, if a user adds an introduction about a movie, the information about the movie should be fetched from IMDB, or if the incoming data is an influencer video about a travel location, the information should be fetched from Google Maps
- Each category can be shared, with access level control
- Also, the user should be able to package some categories together inside a "Circle" concept. For example, users have "Partner" circle, "Friends" circle, "Parents" circle, "Colleagues" circle, and etc.
- It should be possible to share a whole circle which has some categories, with others, in a collabrative manner.
- The app should persist users' data on the servers, and users can view and edit their list from a phone, laptop, etc. For the first step, we focus on only web applications with front-end and back-end, ignoring the native Android or iOS apps.
- It should have Register and Login features
- Users can log in with their ID providers like Google, Apple, or X, etc
- The app should be something like Google Tasks, simple without distractions, material design.

## Memory

## Purpose & context

Developer is developing a cross-platform web application for tracking and organizing hobbies with AI-powered categorization and multi-source import capabilities. The project emphasizes creating a clean, user-friendly experience that avoids the visual clutter common in modern applications. Success is measured by achieving a distraction-free interface that allows users to focus on their hobby content without unnecessary visual noise.

### Current state

The project has moved into focused visual design iteration after initially exploring comprehensive documentation approaches. Developer has established a streamlined workflow that prioritizes rapid design iteration over full documentation generation, allowing for efficient customization and feedback cycles. The current design direction features a professional blue color scheme, list-based layouts grouped by time periods, and minimal visual elements with strategic use of whitespace.

### Key learnings & principles

Several core design principles have emerged: simplicity over feature density, with interfaces designed to minimize distractions; content-first approach where images and thumbnails are hidden in list views and only appear when accessing individual items; list views as the default presentation method rather than card-based layouts; and restrained color application using minimal accents rather than heavy theming throughout the interface.

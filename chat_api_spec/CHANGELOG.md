# VOOBIZE Chat API Specification - Changelog

## [Version 1.1.0] - 2025-01-07

### 🆕 Added - New Endpoints (4 endpoints)

#### Phase 1 (Must Have)
- **GET /chat/messages/:id/context** - Jump to message with context
  - Returns target message + 20 messages before + 20 messages after
  - Use case: Jump from search results, media tab, links tab
  - Includes cursors for loading more in both directions

#### Phase 2 (Nice to Have - Telegram-style features)
- **GET /chat/conversations/:id/media** - List all media in conversation
  - Filter by type: image, video, or all
  - Paginated with cursor
  - Use case: Media gallery view, quick access to shared photos/videos

- **GET /chat/conversations/:id/links** - List all links in conversation
  - Automatically extract URLs from messages
  - Fetch Open Graph metadata (title, description, image)
  - Use case: Quick access to shared URLs

- **GET /chat/conversations/:id/files** - List all files in conversation
  - Filter by file type/MIME type
  - Support documents, archives, text files, code files
  - Use case: File browser view

### 📝 Changes
- Updated total endpoint count from 10 to 14
- Reorganized Messages endpoints section (4 → 8 endpoints)
- Added use case descriptions for new endpoints
- Updated README.md to reflect new endpoint count

### 🎯 Rationale
- **Jump to Message**: Essential for modern chat UX (search, replies, pinned messages)
- **Telegram-style Features**: User expectations based on popular chat apps
- **Phase separation**: Core feature (Phase 1) vs Nice-to-have (Phase 2)

### 📦 Implementation Impact
- Jump to Message: Requires new repository methods (`GetMessagesBeforeTimestamp`, `GetMessagesAfterTimestamp`)
- Media/Links/Files: Requires relationship with existing media system or new tables
- All features: Follow existing Clean Architecture pattern

---

## [Version 1.0.0] - 2024-01-XX (Initial)

### Initial Specification
- 10 REST API endpoints
- 8 WebSocket events
- 3 database tables (conversations, messages, blocks)
- Cursor-based pagination
- Redis caching strategy
- Complete implementation plan (14 weeks)

---

## Notes

### Breaking Changes
- None - All changes are additive

### Deprecations
- None

### Migration Guide
- No migration needed - New endpoints are optional Phase 2 features
- Jump to Message (Phase 1) is backward compatible

---

**Last Updated**: 2025-01-07
**Status**: Specification Complete - Ready for Implementation

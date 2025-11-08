# Chat API Specification - Pagination & Infinite Scroll

## Overview

ระบบ Chat ใช้ **Cursor-based Pagination** เพื่อรองรับ Infinite Scroll ที่มีประสิทธิภาพและสม่ำเสมอ

### ทำไมไม่ใช้ Offset-based Pagination?

**ปัญหาของ Offset**:
```sql
-- หน้า 1
SELECT * FROM messages WHERE conversation_id = 'conv-1'
ORDER BY created_at DESC LIMIT 20 OFFSET 0;

-- หน้า 2
SELECT * FROM messages WHERE conversation_id = 'conv-1'
ORDER BY created_at DESC LIMIT 20 OFFSET 20;

-- ปัญหา: ถ้ามี message ใหม่เข้ามาระหว่างการดึงหน้า 1 และ 2
-- จะมี messages ซ้ำซ้อนหรือหายไปได้
```

**ข้อดีของ Cursor**:
- Consistent results แม้มีข้อมูลใหม่เข้ามา
- Performance ดีกว่า (ไม่ต้อง skip rows)
- Suitable สำหรับ real-time data

---

## Cursor Design

### Cursor Structure

Cursor เป็น Base64-encoded JSON object:

```typescript
interface Cursor {
  created_at: string;   // ISO 8601 timestamp
  id: string;           // UUID (tie-breaker)
}

// Example
const cursor = {
  created_at: "2024-01-01T10:00:00Z",
  id: "msg-050"
};

const encoded = btoa(JSON.stringify(cursor));
// "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0wMVQxMDowMDowMFoiLCJpZCI6Im1zZy0wNTAifQ=="
```

**ทำไมต้องมี `id` ด้วย?**
- เพื่อ handle กรณีที่มี `created_at` ซ้ำกัน (tie-breaker)
- ทำให้ pagination มั่นใจได้ว่าไม่มี duplicates

---

## Conversations Pagination

### API Endpoint
```
GET /chat/conversations?cursor={cursor}&limit={limit}
```

### SQL Query

#### Page 1 (No cursor)
```sql
SELECT
    c.id,
    c.user1_id,
    c.user2_id,
    c.last_message_at,
    c.updated_at,
    -- จอิน user และ message data
FROM conversations c
WHERE (c.user1_id = $user_id OR c.user2_id = $user_id)
ORDER BY c.updated_at DESC, c.id DESC
LIMIT 20;
```

#### Page 2+ (With cursor)
```sql
SELECT
    c.id,
    c.user1_id,
    c.user2_id,
    c.last_message_at,
    c.updated_at
FROM conversations c
WHERE (c.user1_id = $user_id OR c.user2_id = $user_id)
  AND (
    c.updated_at < $cursor_updated_at
    OR (c.updated_at = $cursor_updated_at AND c.id < $cursor_id)
  )
ORDER BY c.updated_at DESC, c.id DESC
LIMIT 20;
```

**Explanation**:
- `c.updated_at < $cursor_updated_at`: ดึง conversations ที่เก่ากว่า cursor
- `OR (c.updated_at = ... AND c.id < ...)`: Handle tie-breaking

### Response Format
```json
{
  "success": true,
  "data": {
    "conversations": [ /* ... */ ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJ1cGRhdGVkX2F0IjoiMjAyNC0wMS0wMVQwOTowMDowMFoiLCJpZCI6ImNvbnYtMDIwIn0"
    }
  }
}
```

### Frontend Implementation

```typescript
import { useInfiniteQuery } from '@tanstack/react-query';

interface ConversationsCursor {
  updated_at: string;
  id: string;
}

export function useConversations() {
  return useInfiniteQuery({
    queryKey: ['conversations'],
    queryFn: async ({ pageParam }) => {
      const params = new URLSearchParams();
      if (pageParam) {
        params.append('cursor', pageParam);
      }
      params.append('limit', '20');

      const response = await fetch(`/api/chat/conversations?${params}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      return response.json();
    },
    getNextPageParam: (lastPage) => {
      return lastPage.data.meta.hasMore
        ? lastPage.data.meta.nextCursor
        : undefined;
    },
    initialPageParam: undefined,
  });
}

// Component
function ConversationsList() {
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useConversations();

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop, scrollHeight, clientHeight } = e.currentTarget;
    const threshold = 100; // pixels from bottom

    if (
      scrollHeight - scrollTop - clientHeight < threshold &&
      hasNextPage &&
      !isFetchingNextPage
    ) {
      fetchNextPage();
    }
  };

  return (
    <ScrollArea onScroll={handleScroll}>
      {data?.pages.map((page) =>
        page.data.conversations.map((conv) => (
          <ConversationItem key={conv.id} conversation={conv} />
        ))
      )}
      {isFetchingNextPage && <Loader />}
    </ScrollArea>
  );
}
```

---

## Messages Pagination

### API Endpoint
```
GET /chat/conversations/:conversationId/messages?cursor={cursor}&limit={limit}
```

### SQL Query

#### Page 1 (Most recent messages - no cursor)
```sql
SELECT
    m.id,
    m.conversation_id,
    m.sender_id,
    m.content,
    m.is_read,
    m.read_at,
    m.created_at,
    m.updated_at
FROM messages m
WHERE m.conversation_id = $conversation_id
  AND m.deleted_at IS NULL
ORDER BY m.created_at DESC, m.id DESC
LIMIT 50;
```

#### Page 2+ (Load older messages - with cursor)
```sql
SELECT
    m.id,
    m.conversation_id,
    m.sender_id,
    m.content,
    m.is_read,
    m.read_at,
    m.created_at,
    m.updated_at
FROM messages m
WHERE m.conversation_id = $conversation_id
  AND m.deleted_at IS NULL
  AND (
    m.created_at < $cursor_created_at
    OR (m.created_at = $cursor_created_at AND m.id < $cursor_id)
  )
ORDER BY m.created_at DESC, m.id DESC
LIMIT 50;
```

**Note**: Messages เรียงจาก **ใหม่ไปเก่า** (DESC) เพราะ:
- UI แสดงข้อความล่าสุดก่อน
- Infinite scroll โหลดข้อความเก่าๆ เมื่อ scroll ขึ้น

### Response Format
```json
{
  "success": true,
  "data": {
    "messages": [
      {
        "id": "msg-100",
        "conversationId": "conv-001",
        "senderId": "user-456",
        "content": "ข้อความล่าสุด",
        "createdAt": "2024-01-01T11:00:00Z",
        "isRead": false
      },
      {
        "id": "msg-099",
        "content": "ข้อความก่อนหน้า",
        "createdAt": "2024-01-01T10:59:00Z",
        "isRead": true
      }
      /* ... ย้อนหลังไปเรื่อยๆ */
    ],
    "meta": {
      "hasMore": true,
      "nextCursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0wMS0wMVQxMDowMDowMFoiLCJpZCI6Im1zZy0wNTAifQ"
    }
  }
}
```

### Frontend Implementation

```typescript
import { useInfiniteQuery } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';

export function useMessages(conversationId: string) {
  return useInfiniteQuery({
    queryKey: ['messages', conversationId],
    queryFn: async ({ pageParam }) => {
      const params = new URLSearchParams();
      if (pageParam) {
        params.append('cursor', pageParam);
      }
      params.append('limit', '50');

      const response = await fetch(
        `/api/chat/conversations/${conversationId}/messages?${params}`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      return response.json();
    },
    getNextPageParam: (lastPage) => {
      return lastPage.data.meta.hasMore
        ? lastPage.data.meta.nextCursor
        : undefined;
    },
    initialPageParam: undefined,
  });
}

// Component with Reverse Infinite Scroll
function ChatMessages({ conversationId }: { conversationId: string }) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useMessages(conversationId);

  // Auto-scroll to bottom on first load
  useEffect(() => {
    if (scrollRef.current && data?.pages.length === 1) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [data?.pages.length]);

  // Detect scroll to top to load more
  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop } = e.currentTarget;
    const threshold = 100; // pixels from top

    if (scrollTop < threshold && hasNextPage && !isFetchingNextPage) {
      const prevHeight = scrollRef.current?.scrollHeight || 0;

      fetchNextPage().then(() => {
        // Maintain scroll position after loading older messages
        if (scrollRef.current) {
          const newHeight = scrollRef.current.scrollHeight;
          scrollRef.current.scrollTop = newHeight - prevHeight;
        }
      });
    }
  };

  // Flatten all messages from all pages
  const allMessages = data?.pages.flatMap(page => page.data.messages) || [];

  return (
    <ScrollArea ref={scrollRef} onScroll={handleScroll}>
      {isFetchingNextPage && <Loader className="mx-auto my-2" />}

      {/* Reverse order: oldest first */}
      {allMessages.reverse().map((message) => (
        <MessageBubble key={message.id} message={message} />
      ))}
    </ScrollArea>
  );
}
```

### Optimistic UI Update

เมื่อส่งข้อความใหม่:

```typescript
import { useMutation, useQueryClient } from '@tanstack/react-query';

export function useSendMessage(conversationId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (content: string) => {
      // Send via WebSocket or REST
      const response = await fetch(
        `/api/chat/conversations/${conversationId}/messages`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ content }),
        }
      );
      return response.json();
    },
    onMutate: async (content) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries({
        queryKey: ['messages', conversationId],
      });

      // Snapshot previous value
      const previousMessages = queryClient.getQueryData([
        'messages',
        conversationId,
      ]);

      // Optimistically update with temporary message
      queryClient.setQueryData(['messages', conversationId], (old: any) => ({
        ...old,
        pages: [
          {
            ...old.pages[0],
            data: {
              ...old.pages[0].data,
              messages: [
                {
                  id: `temp-${Date.now()}`,
                  conversationId,
                  senderId: currentUserId,
                  content,
                  isRead: false,
                  createdAt: new Date().toISOString(),
                  isPending: true, // Custom flag
                },
                ...old.pages[0].data.messages,
              ],
            },
          },
          ...old.pages.slice(1),
        ],
      }));

      return { previousMessages };
    },
    onError: (err, variables, context) => {
      // Rollback on error
      if (context?.previousMessages) {
        queryClient.setQueryData(
          ['messages', conversationId],
          context.previousMessages
        );
      }
    },
    onSuccess: (data) => {
      // Replace temporary message with real one
      queryClient.setQueryData(['messages', conversationId], (old: any) => ({
        ...old,
        pages: [
          {
            ...old.pages[0],
            data: {
              ...old.pages[0].data,
              messages: old.pages[0].data.messages.map((msg: any) =>
                msg.isPending ? data.data : msg
              ),
            },
          },
          ...old.pages.slice(1),
        ],
      }));
    },
  });
}
```

---

## Backend Cursor Encoding/Decoding

### Go Implementation

```go
package pagination

import (
    "encoding/base64"
    "encoding/json"
    "time"
)

type Cursor struct {
    CreatedAt time.Time `json:"created_at"`
    ID        string    `json:"id"`
}

// EncodeCursor encodes cursor to base64 string
func EncodeCursor(createdAt time.Time, id string) (string, error) {
    cursor := Cursor{
        CreatedAt: createdAt,
        ID:        id,
    }

    jsonBytes, err := json.Marshal(cursor)
    if err != nil {
        return "", err
    }

    encoded := base64.StdEncoding.EncodeToString(jsonBytes)
    return encoded, nil
}

// DecodeCursor decodes base64 cursor string
func DecodeCursor(encodedCursor string) (*Cursor, error) {
    if encodedCursor == "" {
        return nil, nil
    }

    jsonBytes, err := base64.StdEncoding.DecodeString(encodedCursor)
    if err != nil {
        return nil, err
    }

    var cursor Cursor
    err = json.Unmarshal(jsonBytes, &cursor)
    if err != nil {
        return nil, err
    }

    return &cursor, nil
}
```

### Usage in Handler

```go
func (h *ConversationHandler) GetMessages(c *gin.Context) {
    conversationID := c.Param("conversationId")
    encodedCursor := c.Query("cursor")
    limit := c.GetInt("limit")
    if limit == 0 {
        limit = 50
    }

    // Decode cursor
    cursor, err := pagination.DecodeCursor(encodedCursor)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid cursor"})
        return
    }

    // Build query
    query := h.db.Where("conversation_id = ?", conversationID).
        Where("deleted_at IS NULL")

    if cursor != nil {
        query = query.Where(
            "(created_at < ? OR (created_at = ? AND id < ?))",
            cursor.CreatedAt,
            cursor.CreatedAt,
            cursor.ID,
        )
    }

    var messages []Message
    err = query.
        Order("created_at DESC, id DESC").
        Limit(limit + 1). // Fetch one extra to check hasMore
        Find(&messages).Error

    if err != nil {
        c.JSON(500, gin.H{"error": "Database error"})
        return
    }

    // Check if has more
    hasMore := len(messages) > limit
    if hasMore {
        messages = messages[:limit]
    }

    // Generate next cursor
    var nextCursor string
    if hasMore && len(messages) > 0 {
        lastMsg := messages[len(messages)-1]
        nextCursor, _ = pagination.EncodeCursor(lastMsg.CreatedAt, lastMsg.ID)
    }

    c.JSON(200, gin.H{
        "success": true,
        "data": gin.H{
            "messages": messages,
            "meta": gin.H{
                "hasMore":    hasMore,
                "nextCursor": nextCursor,
            },
        },
    })
}
```

---

## Performance Optimization

### 1. Database Indexes

```sql
-- Conversations pagination
CREATE INDEX idx_conversations_pagination
ON conversations(updated_at DESC, id DESC)
WHERE deleted_at IS NULL;

-- Messages pagination
CREATE INDEX idx_messages_pagination
ON messages(conversation_id, created_at DESC, id DESC)
WHERE deleted_at IS NULL;
```

### 2. Query Optimization

**Use LIMIT + 1 Pattern**:
```go
// Fetch one extra to determine hasMore
query.Limit(limit + 1).Find(&results)

hasMore := len(results) > limit
if hasMore {
    results = results[:limit]
}
```

**Benefits**:
- Only one query instead of two (SELECT + COUNT)
- Faster response time

### 3. Caching Strategy

**Cache first page**:
```
Key: "conversations:{user_id}:first_page"
TTL: 30 seconds
```

**Invalidate on**:
- New message received
- Conversation updated

---

## Edge Cases

### 1. Real-time Updates

**Problem**: ข้อความใหม่เข้ามาระหว่าง pagination

**Solution**: ใช้ WebSocket เพิ่มข้อความใหม่ที่ด้านบน (prepend) แทนการ refetch

```typescript
// Listen to WebSocket
useEffect(() => {
  const handleNewMessage = (message: Message) => {
    queryClient.setQueryData(['messages', conversationId], (old: any) => ({
      ...old,
      pages: [
        {
          ...old.pages[0],
          data: {
            ...old.pages[0].data,
            messages: [message, ...old.pages[0].data.messages],
          },
        },
        ...old.pages.slice(1),
      ],
    }));
  };

  socket.on('message.new', handleNewMessage);
  return () => socket.off('message.new', handleNewMessage);
}, [conversationId]);
```

### 2. Scroll Position Restoration

**Problem**: เมื่อโหลดข้อความเก่า scroll position กระโดด

**Solution**: บันทึก `scrollHeight` ก่อนโหลดและปรับ `scrollTop` หลังโหลด

```typescript
const handleLoadMore = async () => {
  const prevHeight = scrollRef.current?.scrollHeight || 0;

  await fetchNextPage();

  // Restore scroll position
  if (scrollRef.current) {
    const newHeight = scrollRef.current.scrollHeight;
    scrollRef.current.scrollTop = newHeight - prevHeight;
  }
};
```

### 3. Initial Load Position

**Options**:
1. **Load most recent** (recommended): ดีสำหรับ active conversations
2. **Load from unread**: ดีสำหรับกลับมาเช็คข้อความเก่า

```typescript
// Option 2: Load from unread
const { data } = useMessages(conversationId, {
  initialCursor: conversation.lastReadMessageId
    ? encodeUnreadCursor(conversation.lastReadMessageId)
    : undefined,
});
```

---

## Testing Strategy

### Unit Tests

```typescript
describe('Cursor Encoding/Decoding', () => {
  it('should encode and decode cursor correctly', () => {
    const original = {
      created_at: '2024-01-01T10:00:00Z',
      id: 'msg-123',
    };

    const encoded = encodeCursor(original);
    const decoded = decodeCursor(encoded);

    expect(decoded).toEqual(original);
  });

  it('should handle empty cursor', () => {
    const decoded = decodeCursor('');
    expect(decoded).toBeNull();
  });
});
```

### Integration Tests

```go
func TestMessagesPagination(t *testing.T) {
    // Setup: Create 100 messages
    for i := 1; i <= 100; i++ {
        createMessage(conversationID, fmt.Sprintf("Message %d", i))
    }

    // Page 1
    resp1 := getMessages(conversationID, "", 20)
    assert.Equal(t, 20, len(resp1.Messages))
    assert.True(t, resp1.Meta.HasMore)

    // Page 2
    resp2 := getMessages(conversationID, resp1.Meta.NextCursor, 20)
    assert.Equal(t, 20, len(resp2.Messages))
    assert.True(t, resp2.Meta.HasMore)

    // Check no duplicates
    ids1 := getIDs(resp1.Messages)
    ids2 := getIDs(resp2.Messages)
    assert.Empty(t, intersection(ids1, ids2))
}
```

---

## Monitoring

### Metrics to Track

1. **Average page load time**: < 100ms
2. **Cache hit rate**: > 80% for first page
3. **Pagination errors**: < 0.1%
4. **Duplicate messages**: 0

### Logging

```go
log.Info("Pagination request",
    zap.String("user_id", userID),
    zap.String("cursor", cursor),
    zap.Int("limit", limit),
    zap.Int("result_count", len(results)),
    zap.Bool("has_more", hasMore),
    zap.Duration("duration", duration),
)
```

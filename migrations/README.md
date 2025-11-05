# Database Migrations

## วิธีรัน Migration สำหรับ Push Subscriptions Unique Constraint

### วิธีที่ 1: ใช้ psql command line

```bash
# Connect to PostgreSQL database
psql -h localhost -p 5432 -U postgres -d gofiber_social

# รัน migration file
\i migrations/add_push_subscriptions_unique_constraint.sql
```

### วิธีที่ 2: ใช้ psql แบบ one-liner

```bash
# Windows (Git Bash / MINGW64)
psql -h localhost -p 5432 -U postgres -d gofiber_social -f migrations/add_push_subscriptions_unique_constraint.sql

# รอใส่ password: n147369
```

### วิธีที่ 3: ใช้ PGPASSWORD environment variable

```bash
# Windows (Git Bash / MINGW64)
PGPASSWORD=n147369 psql -h localhost -p 5432 -U postgres -d gofiber_social -f migrations/add_push_subscriptions_unique_constraint.sql
```

### วิธีที่ 4: ใช้ connection string

```bash
psql "postgresql://postgres:n147369@localhost:5432/gofiber_social" -f migrations/add_push_subscriptions_unique_constraint.sql
```

---

## ตรวจสอบว่า Migration สำเร็จหรือไม่

```sql
-- ตรวจสอบ constraints
SELECT conname, contype, pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conrelid = 'push_subscriptions'::regclass;
```

**Expected output:**
```
        conname         | contype |                    pg_get_constraintdef
------------------------+---------+---------------------------------------------------
 push_subscriptions_pkey| p       | PRIMARY KEY (id)
 idx_user_endpoint      | u       | UNIQUE (user_id, endpoint)
```

---

## Troubleshooting

### Error: "duplicate key value violates unique constraint"

หมายความว่า database มี duplicate (user_id, endpoint) อยู่

**แก้ไข:**
```sql
-- ดูว่ามี duplicates อะไรบ้าง
SELECT user_id, endpoint, COUNT(*)
FROM push_subscriptions
GROUP BY user_id, endpoint
HAVING COUNT(*) > 1;

-- ลบ duplicates (เก็บแค่ record ล่าสุด)
DELETE FROM push_subscriptions
WHERE id NOT IN (
    SELECT DISTINCT ON (user_id, endpoint) id
    FROM push_subscriptions
    ORDER BY user_id, endpoint, updated_at DESC
);
```

### Error: "constraint already exists"

แสดงว่า constraint มีอยู่แล้ว ไม่ต้องทำอะไร

**ตรวจสอบ:**
```sql
SELECT conname FROM pg_constraint
WHERE conrelid = 'push_subscriptions'::regclass
  AND conname = 'idx_user_endpoint';
```

---

## หลังจากรัน Migration

1. **Restart Backend Server**
   ```bash
   # หยุด server (Ctrl+C)
   # รันใหม่
   go run main.go
   ```

2. **ทดสอบ Push Subscription API**
   ```bash
   curl -X POST http://localhost:8080/api/v1/push/subscribe \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "endpoint": "https://fcm.googleapis.com/fcm/send/test123",
       "keys": {
         "p256dh": "test-key",
         "auth": "test-auth"
       }
     }'
   ```

   **Expected response:**
   ```json
   {
     "success": true,
     "message": "Subscription saved successfully"
   }
   ```

3. **ลองส่งซ้ำ (ควร UPDATE แทนที่จะ error)**
   - รัน curl command เดิมอีกครั้ง
   - ควรได้ response success เหมือนเดิม (ไม่ error)

---

## Alternative: ใช้ GORM AutoMigrate

หาก model มี `uniqueIndex` tag ถูกต้องแล้ว สามารถให้ GORM สร้าง constraint ให้อัตโนมัติได้:

```go
// infrastructure/postgres/database.go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // ... other models
        &models.PushSubscription{},
    )
}
```

**Note:** GORM AutoMigrate จะ:
- สร้าง constraint ถ้ายังไม่มี
- ไม่ลบ duplicates ให้ (ต้องลบเอง)
- อาจไม่ได้ rename constraint ถ้ามีอยู่แล้ว

**แนะนำ:** รัน SQL migration เพื่อควบคุมการเปลี่ยนแปลงได้ชัดเจน

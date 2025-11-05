-- Migration: Add UNIQUE constraint to push_subscriptions table
-- Purpose: Fix "there is no unique or exclusion constraint matching the ON CONFLICT specification" error
-- Date: 2025-01-06

-- Step 1: ตรวจสอบ constraints ที่มีอยู่
-- SELECT conname, contype, pg_get_constraintdef(oid)
-- FROM pg_constraint
-- WHERE conrelid = 'push_subscriptions'::regclass;

-- Step 2: ลบ duplicate rows ก่อน (ถ้ามี) เพื่อไม่ให้ ADD CONSTRAINT ล้มเหลว
-- Keep only the most recent record for each (user_id, endpoint) combination
DELETE FROM push_subscriptions
WHERE id NOT IN (
    SELECT DISTINCT ON (user_id, endpoint) id
    FROM push_subscriptions
    ORDER BY user_id, endpoint, updated_at DESC
);

-- Step 3: เพิ่ม UNIQUE constraint (user_id, endpoint)
-- ใช้ชื่อ idx_user_endpoint ตาม GORM model
ALTER TABLE push_subscriptions
ADD CONSTRAINT idx_user_endpoint UNIQUE (user_id, endpoint);

-- Step 4: เพิ่ม index เพิ่มเติมตาม spec (ถ้ายังไม่มี)
CREATE INDEX IF NOT EXISTS idx_push_subscriptions_user_id ON push_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_push_subscriptions_endpoint ON push_subscriptions(endpoint);

-- Verify the constraint was added
SELECT conname, contype, pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conrelid = 'push_subscriptions'::regclass;

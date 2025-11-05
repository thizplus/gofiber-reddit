# Error Codes Reference

## Overview
This document lists all error codes used in the API with their descriptions and recommended client handling.

---

## Error Response Format

```json
{
  "success": false,
  "message": "Human-readable error message",
  "code": "ERROR_CODE",
  "errors": {
    "field": "Field-specific error"
  },
  "details": "Additional error details (development only)"
}
```

---

## Authentication Errors (AUTH_*)

### AUTH_001 - Invalid Credentials
**Message:** ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง
**Status:** 401
**Cause:** Username/password combination doesn't match
**Action:** Ask user to check credentials

### AUTH_002 - Token Missing
**Message:** ไม่พบ Token กรุณาเข้าสู่ระบบ
**Status:** 401
**Cause:** No Authorization header provided
**Action:** Redirect to login

### AUTH_003 - Token Invalid
**Message:** Token ไม่ถูกต้อง กรุณาเข้าสู่ระบบใหม่
**Status:** 401
**Cause:** Token is malformed or tampered
**Action:** Clear token, redirect to login

### AUTH_004 - Token Expired
**Message:** Token หมดอายุ กรุณาเข้าสู่ระบบใหม่
**Status:** 401
**Cause:** Token has expired
**Action:** Attempt token refresh, or redirect to login

### AUTH_005 - User Not Found
**Message:** ไม่พบผู้ใช้
**Status:** 401
**Cause:** User associated with token no longer exists
**Action:** Clear token, redirect to login

### AUTH_006 - Account Suspended
**Message:** บัญชีของคุณถูกระงับ
**Status:** 403
**Cause:** User account has been suspended
**Action:** Show suspension message, contact support

### AUTH_007 - Account Deleted
**Message:** บัญชีนี้ถูกลบแล้ว
**Status:** 410
**Cause:** User account has been deleted
**Action:** Clear token, redirect to login

---

## Validation Errors (VAL_*)

### VAL_001 - Missing Required Field
**Message:** กรุณากรอก {field}
**Status:** 400
**Cause:** Required field is missing
**Action:** Show field-specific error

### VAL_002 - Invalid Format
**Message:** รูปแบบ {field} ไม่ถูกต้อง
**Status:** 400
**Cause:** Field value doesn't match expected format
**Action:** Show format requirements

### VAL_003 - Value Too Short
**Message:** {field} ต้องมีอย่างน้อย {min} ตัวอักษร
**Status:** 400
**Cause:** String value is below minimum length
**Action:** Show minimum length requirement

### VAL_004 - Value Too Long
**Message:** {field} ต้องไม่เกิน {max} ตัวอักษร
**Status:** 400
**Cause:** String value exceeds maximum length
**Action:** Show maximum length requirement

### VAL_005 - Invalid Email
**Message:** รูปแบบอีเมลไม่ถูกต้อง
**Status:** 400
**Cause:** Email format is invalid
**Action:** Show valid email example

### VAL_006 - Invalid URL
**Message:** รูปแบบ URL ไม่ถูกต้อง
**Status:** 400
**Cause:** URL format is invalid
**Action:** Show valid URL example

### VAL_007 - Invalid Enum Value
**Message:** ค่า {field} ไม่ถูกต้อง ต้องเป็น {validValues}
**Status:** 400
**Cause:** Value not in allowed enum values
**Action:** Show valid options

---

## Resource Errors (RES_*)

### RES_001 - Post Not Found
**Message:** ไม่พบโพสต์
**Status:** 404
**Cause:** Post ID doesn't exist
**Action:** Show "post not found" message

### RES_002 - Comment Not Found
**Message:** ไม่พบความคิดเห็น
**Status:** 404
**Cause:** Comment ID doesn't exist
**Action:** Show "comment not found" message

### RES_003 - User Not Found
**Message:** ไม่พบผู้ใช้
**Status:** 404
**Cause:** Username doesn't exist
**Action:** Show "user not found" message

### RES_004 - Media Not Found
**Message:** ไม่พบไฟล์
**Status:** 404
**Cause:** Media ID doesn't exist
**Action:** Show "file not found" message

### RES_005 - Notification Not Found
**Message:** ไม่พบการแจ้งเตือน
**Status:** 404
**Cause:** Notification ID doesn't exist
**Action:** Refresh notification list

---

## Permission Errors (PERM_*)

### PERM_001 - Not Authorized
**Message:** คุณไม่มีสิทธิ์ในการดำเนินการนี้
**Status:** 403
**Cause:** User doesn't have permission for action
**Action:** Show permission denied message

### PERM_002 - Not Post Owner
**Message:** คุณไม่มีสิทธิ์แก้ไขโพสต์นี้
**Status:** 403
**Cause:** User is not the post author
**Action:** Disable edit/delete buttons

### PERM_003 - Not Comment Owner
**Message:** คุณไม่มีสิทธิ์แก้ไขความคิดเห็นนี้
**Status:** 403
**Cause:** User is not the comment author
**Action:** Disable edit/delete buttons

### PERM_004 - Not Media Owner
**Message:** คุณไม่มีสิทธิ์ลบไฟล์นี้
**Status:** 403
**Cause:** User didn't upload the media
**Action:** Disable delete button

---

## Conflict Errors (CONF_*)

### CONF_001 - Username Taken
**Message:** ชื่อผู้ใช้นี้ถูกใช้งานแล้ว
**Status:** 409
**Cause:** Username already exists
**Action:** Ask user to choose different username

### CONF_002 - Email Taken
**Message:** อีเมลนี้ถูกใช้งานแล้ว
**Status:** 409
**Cause:** Email already registered
**Action:** Ask user to use different email or login

### CONF_003 - Already Voted
**Message:** คุณได้โหวตแล้ว
**Status:** 409
**Cause:** User already voted on this post/comment
**Action:** Update vote instead

### CONF_004 - Already Following
**Message:** คุณติดตามผู้ใช้นี้อยู่แล้ว
**Status:** 409
**Cause:** User already follows this user
**Action:** Disable follow button

### CONF_005 - Not Following
**Message:** คุณไม่ได้ติดตามผู้ใช้นี้
**Status:** 409
**Cause:** User doesn't follow this user
**Action:** Show follow button instead

### CONF_006 - Already Saved
**Message:** คุณได้บันทึกโพสต์นี้แล้ว
**Status:** 409
**Cause:** Post already in saved list
**Action:** Show "saved" state

### CONF_007 - Not Saved
**Message:** คุณไม่ได้บันทึกโพสต์นี้
**Status:** 409
**Cause:** Post not in saved list
**Action:** Show "save" button

### CONF_008 - Cannot Self Follow
**Message:** ไม่สามารถติดตามตัวเองได้
**Status:** 409
**Cause:** User trying to follow themselves
**Action:** Hide follow button on own profile

### CONF_009 - Media In Use
**Message:** ไม่สามารถลบไฟล์ที่กำลังใช้งานอยู่
**Status:** 409
**Cause:** Media is attached to posts/comments
**Action:** Remove from posts first

---

## File Upload Errors (FILE_*)

### FILE_001 - Invalid File Type
**Message:** ประเภทไฟล์ไม่ถูกต้อง
**Status:** 400
**Cause:** File type not in allowed list
**Action:** Show allowed file types

### FILE_002 - File Too Large
**Message:** ไฟล์มีขนาดใหญ่เกินไป
**Status:** 400
**Cause:** File exceeds size limit
**Action:** Show maximum file size

### FILE_003 - Too Many Files
**Message:** ไฟล์มากเกินไป
**Status:** 413
**Cause:** More files than allowed per request
**Action:** Show maximum files per upload

### FILE_004 - Image Dimensions Too Large
**Message:** ขนาดรูปภาพใหญ่เกินไป
**Status:** 400
**Cause:** Image dimensions exceed limit
**Action:** Show maximum dimensions

### FILE_005 - Video Duration Too Long
**Message:** วิดีโอยาวเกินไป
**Status:** 400
**Cause:** Video duration exceeds limit
**Action:** Show maximum duration

### FILE_006 - Storage Limit Exceeded
**Message:** พื้นที่จัดเก็บเต็ม
**Status:** 413
**Cause:** User's storage quota is full
**Action:** Ask to delete old files

### FILE_007 - Upload Failed
**Message:** อัปโหลดไฟล์ไม่สำเร็จ
**Status:** 500
**Cause:** Storage service error
**Action:** Ask user to retry

### FILE_008 - Malicious File Detected
**Message:** พบไฟล์ที่เป็นอันตราย
**Status:** 400
**Cause:** File contains malware
**Action:** Reject file, show warning

---

## Rate Limit Errors (RATE_*)

### RATE_001 - Too Many Requests
**Message:** คำขอมากเกินไป กรุณาลองใหม่ในอีก {seconds} วินาที
**Status:** 429
**Cause:** Rate limit exceeded
**Action:** Show countdown, disable actions

### RATE_002 - Too Many Login Attempts
**Message:** พยายามเข้าสู่ระบบมากเกินไป กรุณาลองใหม่ในอีก 15 นาที
**Status:** 429
**Cause:** Too many failed login attempts
**Action:** Show lockout message

### RATE_003 - Too Many Uploads
**Message:** อัปโหลดไฟล์มากเกินไป กรุณารออีก {seconds} วินาที
**Status:** 429
**Cause:** Upload rate limit exceeded
**Action:** Disable upload, show countdown

### RATE_004 - Too Many Posts
**Message:** สร้างโพสต์มากเกินไป กรุณารอสักครู่
**Status:** 429
**Cause:** Post creation rate limit exceeded
**Action:** Disable post button temporarily

---

## Server Errors (SRV_*)

### SRV_001 - Internal Server Error
**Message:** เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์
**Status:** 500
**Cause:** Unhandled server error
**Action:** Show generic error, ask to retry

### SRV_002 - Database Error
**Message:** ไม่สามารถเชื่อมต่อฐานข้อมูล
**Status:** 503
**Cause:** Database connection failed
**Action:** Show maintenance message

### SRV_003 - External Service Error
**Message:** บริการภายนอกไม่พร้อมใช้งาน
**Status:** 503
**Cause:** Third-party service unavailable
**Action:** Ask to retry later

### SRV_004 - Maintenance Mode
**Message:** ระบบอยู่ระหว่างปรับปรุง
**Status:** 503
**Cause:** Server in maintenance mode
**Action:** Show maintenance notice

---

## Business Logic Errors (BIZ_*)

### BIZ_001 - Cannot Delete Post With Comments
**Message:** ไม่สามารถลบโพสต์ที่มีความคิดเห็นได้ (future feature)
**Status:** 409
**Cause:** Post has comments
**Action:** Ask to delete comments first

### BIZ_002 - Max Nesting Depth Reached
**Message:** ไม่สามารถตอบกลับลึกกว่านี้ได้
**Status:** 400
**Cause:** Comment nesting depth limit reached
**Action:** Disable reply button

### BIZ_003 - Account Too New
**Message:** บัญชีใหม่ต้องรออีก {hours} ชั่วโมงก่อนทำการนี้ (future feature)
**Status:** 403
**Cause:** Account age restriction
**Action:** Show waiting period

### BIZ_004 - Karma Too Low
**Message:** Karma ต่ำเกินไป ต้องมีอย่างน้อย {required} (future feature)
**Status:** 403
**Cause:** Karma requirement not met
**Action:** Explain karma system

---

## Search Errors (SRCH_*)

### SRCH_001 - Query Too Short
**Message:** คำค้นหาต้องมีอย่างน้อย 2 ตัวอักษร
**Status:** 400
**Cause:** Search query < 2 characters
**Action:** Show minimum length

### SRCH_002 - Query Too Long
**Message:** คำค้นหายาวเกินไป
**Status:** 400
**Cause:** Search query > 100 characters
**Action:** Show maximum length

### SRCH_003 - No Results Found
**Message:** ไม่พบผลลัพธ์
**Status:** 200 (not an error, empty results)
**Cause:** No matches for query
**Action:** Show "no results" message

---

## Client Error Handling

### Recommended Approach

```typescript
async function apiCall(endpoint: string, options: RequestInit) {
  try {
    const response = await fetch(endpoint, options);
    const data = await response.json();

    if (!data.success) {
      switch (data.code) {
        case 'AUTH_002':
        case 'AUTH_003':
        case 'AUTH_004':
          // Redirect to login
          router.push('/login');
          break;

        case 'RATE_001':
          // Show rate limit message
          toast.error(data.message);
          // Retry after X seconds
          break;

        case 'VAL_001':
        case 'VAL_002':
          // Show field errors
          setErrors(data.errors);
          break;

        case 'PERM_001':
        case 'PERM_002':
          // Show permission denied
          toast.error(data.message);
          break;

        default:
          // Generic error handling
          toast.error(data.message);
      }
    }

    return data;
  } catch (error) {
    // Network error
    toast.error('ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์');
    throw error;
  }
}
```

---

## Error Logging

All errors should be logged with:
- Timestamp
- User ID (if authenticated)
- Endpoint
- Request payload
- Error code
- Stack trace (server-side only)

Sensitive data (passwords, tokens) should NOT be logged.

---
description: Create Bruno Documentation
---

1. **Hierarchy**: Save to `docs/bruno/[Module Name]/[Request Name].bru`.
2. **Standard Variables**:
   - URL: Always use `{{host}}/path`.
   - Auth: Use `auth: bearer` with `token: {{auth_token}}`.
3. **Login Automation**:
   - For Login requests, add `script:post-response`:
     ```javascript
     if (res.status === 200 && res.body.data && res.body.data.access_token) {
       bru.setVar("auth_token", res.body.data.access_token);
     }
     ```
4. **Validation**: Use `body:json` for POST/PUT. Use `params:path` or `params:query` for parameters.
5. **Sync**: Ensure the Bruno collection is updated whenever a new endpoint is added.

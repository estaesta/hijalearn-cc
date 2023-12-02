API (by assumption)

TODO
---
- [ ] Update users profile
- [ ] Deploy

Installation
---
Pre-requisite: Docker

```bash
$ git clone
$ cd hijalearn-cc
$ docker build -t hijalearn-cc .
$ docker run -d -p 8080:8080 hijalearn-cc
```

PUT /api/v1/progress
---
Update the user's progress.
Include subab to update the subab progress.
Do not include subab to update the bab progress (all subab completed).

Header

Authorization: Bearer \<token\>

Body

multipart/form-data
- bab: bab number
- subab: subab number


POST /api/v1/register
---
Register a new user.

Body

multipart/form-data
- email: email
- password: password
- username: username

GET /api/v1/progress
---
Get the user's progress.

Header

Authorization: Bearer \<token\>

Response

```json
{
  "bab": {
    "1": {
      "selesai": true,
      "subab": {
        "1": true,
        "2": true
      }
    }
  }
}
```

API

TODO
---
- [ ] Update users profile

Installation
---
Pre-requisite: Docker

```bash
$ git clone
$ cd hijalearn-cc
$ docker build -t hijalearn-cc .
$ docker run -d -p 8080:8080 hijalearn-cc
```

POST /api/v1/register
---
Register a new user.

Body

multipart/form-data
- email: email
- password: password
- username: username


POST /api/v1/prediction
---
Make a prediction. The backend service will also handle the update progress of the user.

header:

Authorization: Bearer <token>

Body

multipart/form-data
- audio: audio file
- caraEja: label hijaiyah
- ~~moduleId: module id (1-5)~~ ( gk perlu, dihandle backend, ambil dari last module)
- done: boolean (true/false)

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
- moduleId: module id (1-4)
- done: boolean (true/false)

Response
```
"Correct answer"
```

GET /api/v1/progress
---
Get user progress

header:

Authorization: Bearer <token>

Response
```
{
	"last_module": 2,
	"module": [
		{
			"completed": true,
			"module_id": 1,
			"subModuleDone": 30,
			"totalSubModule": 30
		},
		{
			"completed": false,
			"module_id": 2,
			"subModuleDone": 1,
			"totalSubModule": 28
		},
		{
			"completed": false,
			"module_id": 3,
			"subModuleDone": 0,
			"totalSubModule": 28
		},
		{
			"completed": false,
			"module_id": 4,
			"subModuleDone": 0,
			"totalSubModule": 28
		}
	]
}
```

POST /api/v1/progress
---
Initialize user progress (if needed). Already handled by the backend service when user register.

header:

Authorization: Bearer <token>


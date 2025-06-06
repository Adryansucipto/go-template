# go-template-2025/06/03

Make high-level guide on structuring a Clean Architecture in Go for building a RESTful API.

-------------------------------------------------------------------------------------------
```
Depedencies Injection Flows

CONFIG        ← Parses & provides configuration
   ↓
RESOURCE      ← Depends on CONFIG (initializes DB, Redis, S3, etc.)
   ↓
REPOSITORY    ← Depends on RESOURCE (interacts with DB, Redis, etc.)
   ↓
USECASE       ← Depends on REPOSITORY (contains business logic)
   ↓
CONTROLLER    ← Depends on USECASE (handles HTTP requests & responses)
```



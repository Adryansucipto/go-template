# go-template-2025/06/03

Make high-level guide on structuring a Clean Architecture in Go for building a RESTful API.

-------------------------------------------------------------------------------------------

Depedencies Injection Flows

  CONFIG        ← parse & provide configuration
     ||
  RESOURCE      ← membutuhkan CONFIG (untuk inisialisasi: DB, Redis, S3)
     ||
  REPOSITORY    ← membutuhkan RESOURCE (akses DB, Redis, dll)
     ||
  USECASE       ← membutuhkan REPOSITORY dan terkadang RESOURCE langsung
     ||
  CONTROLLER    ← membutuhkan USECASE

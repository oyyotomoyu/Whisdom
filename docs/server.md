# Server

Whisdom's server does most of the platform work. The React frontend is only the user interface; the Go backend owns authentication, permissions, configuration, material handling, model communication, RAG orchestration, logging, and API responses.

The backend is a **Go web server** that exposes REST APIs to the React frontend.

---

# Core Responsibilities

The server is responsible for:

* Connecting frontend requests to backend services
* User authentication
* User roles and permissions
* System configuration
* Training material path configuration
* Conversation storage
* Recent conversation lookup
* Passing user messages to the model
* Returning AI answers to the frontend
* Uploading and processing training materials
* Storing corrected answers
* RAG retrieval and source selection
* Training and inference job coordination
* Audit and system logging

Frontend permission checks are only for user experience. The server must always validate permissions before performing protected actions.

---

# Server Structure

Recommended first structure:

```text
server/
├── main.go
├── apis/
│   ├── router.go
│   ├── auth.go
│   ├── users.go
│   ├── roles.go
│   ├── conversations.go
│   ├── chat.go
│   ├── materials.go
│   ├── corrections.go
│   ├── config.go
│   ├── logs.go
│   └── health.go
├── system/
│   ├── auth.go
│   ├── users.go
│   ├── roles.go
│   ├── permissions.go
│   ├── conversations.go
│   ├── materials.go
│   ├── corrections.go
│   ├── config.go
│   └── database.go
├── ai/
│   ├── model.go
│   ├── rag.go
│   ├── embeddings.go
│   └── training.go
├── logs/
│   ├── service.go
│   └── middleware.go
├── storage/
│   ├── files.go
│   └── paths.go
└── config/
    └── default.json
```

Guidelines:

* `apis` handles HTTP requests and responses.
* `system` owns business rules and database operations.
* `ai` owns model, RAG, embedding, and training orchestration.
* `logs` owns audit/system logging.
* `storage` owns file paths and file persistence.
* API handlers should stay thin and call service functions.

---

# API Connection With Frontend

The frontend communicates with the backend through REST APIs.

Base path:

```text
/api/v1
```

The server should return JSON for normal APIs and support `FormData` for uploads.

Common API behavior:

* Use JSON request and response bodies.
* Return clear error objects.
* Validate authentication before protected routes.
* Validate permissions before admin routes.
* Log important user actions.
* Never return secrets, password hashes, or internal file paths.

Example error:

```json
{
  "error": "permission denied"
}
```

---

# User System

The server owns the user system.

User responsibilities:

* Login
* Logout
* Token refresh or session renewal
* Current user lookup
* User creation
* User update
* User deletion or disable
* Role assignment
* Permission validation

Initial API:

```http
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/auth/refresh
GET  /api/v1/auth/me

GET    /api/v1/users
GET    /api/v1/users/{id}
POST   /api/v1/users
PATCH  /api/v1/users/{id}
DELETE /api/v1/users/{id}

GET   /api/v1/roles
POST  /api/v1/roles
PATCH /api/v1/roles/{id}

GET /api/v1/permissions
```

Possible permissions:

```text
chat.use

conversations.read
conversations.delete

materials.read
materials.upload
materials.delete
materials.path.edit

corrections.create
corrections.read
corrections.update
corrections.delete

users.read
users.manage
roles.read
roles.manage

config.read
config.manage

logs.read

models.read
models.manage
training.run

system.manage
```

---

# Conversation And Chat

The server owns the complete chat flow.

Flow:

```text
Frontend sends message
    ↓
Backend validates user and permission
    ↓
Backend stores user message
    ↓
Backend retrieves relevant company knowledge
    ↓
Backend sends prompt/context to model
    ↓
Model returns answer
    ↓
Backend stores AI answer
    ↓
Backend returns answer to frontend
```

Initial APIs:

```http
POST   /api/v1/conversations
GET    /api/v1/conversations
GET    /api/v1/conversations/{id}
DELETE /api/v1/conversations/{id}

POST /api/v1/conversations/{id}/messages
```

`POST /api/v1/conversations/{id}/messages` is a server job, not a frontend job. The frontend sends the message; the server handles RAG, model calls, response storage, and returned answer data.

Example request:

```json
{
  "message": "How should we handle a refund request?"
}
```

Example response:

```json
{
  "message_id": "msg_123",
  "answer": "The company should first offer a replacement. A refund should only be offered when replacement is not possible.",
  "sources": [
    {
      "material_id": "mat_123",
      "title": "Refund SOP",
      "chunk_id": "chunk_456"
    }
  ]
}
```

The server may later support streaming responses, but the first API can return a normal JSON answer.

---

# Model And RAG

The default model target is **Gemma 4**.

The server should access the model through an internal model interface so the runtime can change later.

Model interface responsibilities:

* Generate answer
* Summarize content
* Use retrieved context
* Return model errors clearly
* Support timeout and cancellation

RAG responsibilities:

* Search vector database
* Select relevant chunks
* Build model context
* Attach source references
* Avoid sending unrelated material to the model

Training responsibilities:

* Prepare approved corrections
* Prepare selected materials
* Trigger training or adapter jobs
* Store training job status
* Keep training separate from normal document upload

Uploading material should not automatically fine-tune the model. It should normally enter the RAG knowledge layer first.

---

# Material And Training Path Config

The server must store a configuration value for the path where model training material is saved or sent.

This config controls the destination used by admin material uploads and training jobs.

Example config:

```json
{
  "training_material_path": "/training/company-default/"
}
```

Requirements:

* Only admin users with `config.manage` or `materials.path.edit` can edit the path.
* The path must be validated before saving.
* The path change must be logged.
* The server must not allow path traversal.
* The server should store previous path changes in logs.

Validation rules:

* Path is required.
* Path must start with `/`.
* Path cannot contain `..`.
* Path cannot contain null bytes.
* Path should be normalized before use.

Config APIs:

```http
GET   /api/v1/config
PATCH /api/v1/config
```

Example update:

```json
{
  "training_material_path": "/training/customer-support/"
}
```

---

# Material APIs

Admin users can upload and manage materials.

```http
GET    /api/v1/materials
GET    /api/v1/materials/{id}
POST   /api/v1/materials
PATCH  /api/v1/materials/{id}
DELETE /api/v1/materials/{id}
POST   /api/v1/materials/{id}/reprocess
```

Upload metadata should include:

* File name
* File type
* Uploaded by
* Upload time
* Destination path
* Processing status
* RAG status
* Training status

The server stores original files unless an authorized user deletes them.

---

# Correct Answer APIs

Admin users can teach the model correct answers.

```http
GET    /api/v1/corrections
GET    /api/v1/corrections/{id}
POST   /api/v1/corrections
PATCH  /api/v1/corrections/{id}
DELETE /api/v1/corrections/{id}
```

Correction data:

```json
{
  "question": "Can a customer receive a refund immediately?",
  "original_answer": "Yes. The customer can request a refund.",
  "corrected_answer": "The company should first offer a replacement. A refund should only be offered when replacement is not possible.",
  "usage": ["rag", "evaluation", "training"]
}
```

Corrections may be used for:

* RAG
* Evaluation datasets
* Prompt examples
* Supervised fine-tuning
* Preference datasets

---

# Log System

The log system is one of the most important server features. Whisdom should follow NexGestion's log system definition.

The server must record system events and user actions consistently.

## Storage

Default log directory:

```text
log/
```

Production may override this with:

```text
LOG_DIR
```

The server must create the directory automatically at startup:

```go
os.MkdirAll("log", 0755)
```

Logs use one JSON Lines file per day:

```text
log/
├── 2026-08-27.log
└── 2026-08-28.log
```

Only the log system may generate filenames. Clients must never provide log filenames.

## Record Format

Each line is one JSON record:

```json
{"timestamp":"2026-08-27 14:35:21 +08:00","status":"info","ip":"192.168.1.20","user_id":"user_123","content":"updated training material path"}
```

Fields:

| Field | Required | Description |
| --- | --- | --- |
| `timestamp` | Yes | `YYYY-MM-DD HH:mm:ss ±HH:MM`, 24-hour time with timezone |
| `status` | Yes | `info`, `warning`, or `error` |
| `ip` | Yes | Client IP address, empty for background jobs |
| `user_id` | Yes | Authenticated user ID, empty for background jobs |
| `content` | Yes | Safe log message without secrets |

## Status Levels

Only these statuses are accepted:

* `info`
* `warning`
* `error`

Invalid status values must return an error.

## Request Logger

After authentication middleware validates the user, log middleware binds:

* Client IP
* User ID

Flow:

```text
HTTP Request
    ↓
Authentication Middleware
    ↓
Log Middleware
    ↓
API Handler
    ↓
logger.Log(status, content)
    ↓
log/YYYY-MM-DD.log
```

API handlers should call:

```go
logger := logs.FromContext(r.Context())
logger.Log("info", "uploaded material mat_123")
```

The caller only provides status and content. IP and user ID come from request context.

Background jobs use system logs with empty IP and user ID.

## Events To Log

The server must log at least:

* Server startup and shutdown
* Successful login
* Failed login
* Logout
* Token refresh success and failure
* User creation, update, disable, and deletion
* Role and permission changes
* Conversation creation and deletion
* Model request failures
* RAG retrieval failures
* Material upload
* Material processing success and failure
* Material deletion
* Training material path changes
* Correction creation, update, and deletion
* Training job start, success, and failure
* Internal API errors

Never log:

* Plaintext passwords
* Password hashes
* Access tokens
* Refresh tokens
* Cookies
* Authorization headers
* JWT signing secret
* Complete sensitive request bodies
* Full private document contents

## Read Log API

```http
GET /api/v1/logs
```

Required permission:

```text
logs.read
```

Query parameters:

| Parameter | Required | Description |
| --- | --- | --- |
| `start` | No | RFC 3339 start time |
| `end` | No | RFC 3339 end time |
| `status` | No | `info`, `warning`, `error`, or comma-separated statuses |
| `limit` | No | Default 100, max 1000 |
| `cursor` | No | Pagination cursor |
| `page` | No | Page number for UI pagination |
| `page_size` | No | Page size for UI pagination |
| `keyword` | No | Search content, user ID, and IP |
| `sort` | No | `timestamp` or `status` |
| `order` | No | `asc` or `desc` |

Example response:

```json
{
  "logs": [
    {
      "timestamp": "2026-08-27 14:35:21 +08:00",
      "status": "warning",
      "ip": "192.168.1.20",
      "user_id": "user_123",
      "content": "model request failed"
    }
  ],
  "next_cursor": ""
}
```

## Retention

Logs must be kept for no more than seven days.

Cleanup runs:

* At server startup
* Every hour while the server is running
* Before creating a new daily log file

Cleanup may delete only log files that match:

```text
YYYY-MM-DD.log
```

---

# Health API

```http
GET /api/v1/health
```

The health API should return:

* Server status
* Database status
* Storage status
* Model service status
* Vector database status

Example:

```json
{
  "status": "ok",
  "database": "ok",
  "storage": "ok",
  "model": "ok",
  "vector_database": "ok"
}
```

---

# Implementation Order

1. Create Go server entry point and centralized router.
2. Add config loading and storage.
3. Add log service based on NexGestion's definition.
4. Add authentication and request logger middleware.
5. Add user, role, and permission services.
6. Add conversation APIs.
7. Add chat API that sends messages to the model.
8. Add material upload and training path config.
9. Add correction APIs.
10. Add RAG and model service interfaces.
11. Add log read API.
12. Add tests for permissions, logs, config validation, and chat flow.

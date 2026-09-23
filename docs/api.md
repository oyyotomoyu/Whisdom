# API List

This document is the quick reference for Whisdom backend APIs.

Base path:

```text
/api/v1
```

Normal APIs should use JSON request and response bodies. Upload APIs may use `FormData`.

## Common Behavior

| Behavior | Rule |
| --- | --- |
| Authentication | Protected APIs must validate the current user. |
| Authorization | Admin or management APIs must validate permissions on the backend. |
| Errors | Return clear JSON error objects. |
| Secrets | Never return secrets, password hashes, internal file paths, or private runtime details. |
| Logs | Important user actions should be logged. |

Example error:

```json
{
  "error": "permission denied"
}
```

## Authentication

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `POST` | `/api/v1/auth/login` | Log in and create a session or token. | Public |
| `POST` | `/api/v1/auth/logout` | Log out and clear the current session or token. | Authenticated |
| `POST` | `/api/v1/auth/refresh` | Refresh or renew authentication. | Authenticated |
| `GET` | `/api/v1/auth/me` | Get the current authenticated user. | Authenticated |

## Users

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/users` | List users. | `users.read` or `users.manage` |
| `GET` | `/api/v1/users/{id}` | Get one user. | `users.read` or `users.manage` |
| `POST` | `/api/v1/users` | Create a user. | `users.manage` |
| `PATCH` | `/api/v1/users/{id}` | Update a user. | `users.manage` |
| `DELETE` | `/api/v1/users/{id}` | Delete or disable a user. | `users.manage` |

## Roles And Permissions

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/roles` | List roles. | `roles.read` or `roles.manage` |
| `POST` | `/api/v1/roles` | Create a role. | `roles.manage` |
| `PATCH` | `/api/v1/roles/{id}` | Update a role. | `roles.manage` |
| `GET` | `/api/v1/permissions` | List available permissions. | `roles.read` or `roles.manage` |

## Conversations And Chat

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `POST` | `/api/v1/conversations` | Create a conversation. | `chat.use` |
| `GET` | `/api/v1/conversations` | List conversations visible to the current user. | `chat.use` or `conversations.read` |
| `GET` | `/api/v1/conversations/{id}` | Read a conversation and its messages. | `chat.use` or `conversations.read` |
| `DELETE` | `/api/v1/conversations/{id}` | Delete a conversation. | Owner or `conversations.delete` |
| `POST` | `/api/v1/conversations/{id}/messages` | Send a user message and receive the AI answer. | `chat.use` |

The message API is owned by the server. The frontend sends the user message; the backend handles storage, RAG retrieval, model calls, answer storage, and source references.

Example request:

```json
{
  "message": "What is our refund policy?"
}
```

Example response:

```json
{
  "message_id": "msg_123",
  "answer": "Replacement should be attempted before a refund.",
  "sources": [
    {
      "material_id": "mat_123",
      "title": "Refund SOP",
      "chunk_id": "chunk_456"
    }
  ]
}
```

## Corrections

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/corrections` | List corrected answers. | `corrections.read` |
| `GET` | `/api/v1/corrections/{id}` | Read one correction. | `corrections.read` |
| `POST` | `/api/v1/corrections` | Create a correction. | `corrections.create` |
| `PATCH` | `/api/v1/corrections/{id}` | Update a correction. | `corrections.update` |
| `DELETE` | `/api/v1/corrections/{id}` | Delete a correction. | `corrections.delete` |
| `POST` | `/api/v1/messages/{id}/corrections` | Create a correction from a specific AI message. | `corrections.create` |

`POST /api/v1/messages/{id}/corrections` is a convenience route for correcting an existing message. The canonical correction resource is `/api/v1/corrections`.

## Materials

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/materials` | List training or knowledge materials. | `materials.read` |
| `GET` | `/api/v1/materials/{id}` | Read material metadata. | `materials.read` |
| `POST` | `/api/v1/materials` | Upload material. | `materials.upload` |
| `PATCH` | `/api/v1/materials/{id}` | Update material metadata. | `materials.upload` or `materials.path.edit` |
| `DELETE` | `/api/v1/materials/{id}` | Delete material. | `materials.delete` |
| `POST` | `/api/v1/materials/{id}/reprocess` | Reprocess material for extraction, chunking, embeddings, or RAG. | `materials.upload` |
| `GET` | `/api/v1/materials/{id}/status` | Read processing, RAG, or training status. | `materials.read` |

Use `/api/v1/materials/{id}/reprocess` for reprocessing. Older notes may call this `/api/v1/materials/{id}/process`; that route should be treated as a legacy alias if implemented.

## Config

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/config` | Read system configuration exposed to admins. | `config.read` or `config.manage` |
| `PATCH` | `/api/v1/config` | Update editable system configuration. | `config.manage` |

Config updates include settings such as the training material path. Path values must be validated and must not allow traversal.

## Models

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/models` | List available model targets. | `models.read` or `models.manage` |
| `GET` | `/api/v1/models/current` | Read the active model target. | `models.read` or `models.manage` |
| `POST` | `/api/v1/models/{id}/activate` | Activate a model target. | `models.manage` |

The API should stay model-independent. Handlers should call an internal model service instead of embedding runtime-specific logic.

## Logs

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/logs` | Read audit or system logs. | `logs.read` |

Supported query parameters:

| Parameter | Required | Description |
| --- | --- | --- |
| `start` | No | RFC 3339 start time. |
| `end` | No | RFC 3339 end time. |
| `status` | No | `info`, `warning`, `error`, or comma-separated statuses. |
| `limit` | No | Default `100`, max `1000`. |
| `cursor` | No | Pagination cursor. |
| `page` | No | Page number for client pagination. |
| `page_size` | No | Page size for client pagination. |
| `keyword` | No | Search content, user ID, and IP. |
| `sort` | No | `timestamp` or `status`. |
| `order` | No | `asc` or `desc`. |

## Health

| Method | Path | Purpose | Permission |
| --- | --- | --- | --- |
| `GET` | `/api/v1/health` | Read server and dependency health. | Public or internal |

The health response should include server, database, storage, model service, and vector database status.

## Permission Names

Initial permissions:

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

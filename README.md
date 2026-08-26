# Whisdom
An open-source, self-hosted AI platform that helps businesses build and train private AI with their own knowledge.
# Private Enterprise AI Platform

A lightweight, self-hosted AI platform designed for small and medium-sized businesses.

The platform allows organizations to deploy a customizable language model inside their own VM or private infrastructure. Company knowledge can be continuously expanded through internal documents, images, audio recordings, and corrections provided by authorized users.

The primary goals are:

* **Private by default** — company data stays inside the organization's infrastructure.
* **Low deployment cost** — designed for small VMs and affordable hardware.
* **Easy to deploy** — minimize infrastructure and AI/ML expertise requirements.
* **Customizable** — organizations can build their own AI knowledge base.
* **Human-guided improvement** — authorized users can correct AI responses and provide better answers.
* **Model-independent** — the platform should not be permanently tied to a single LLM provider.
* **Cross-device** — users access the system through a responsive web interface.

---

# Use Cases

The platform can be customized for different internal business scenarios.

Examples include:

* Customer service assistants
* Employee training
* Internal knowledge assistants
* SOP assistants
* Product support
* Technical support
* Internal Q&A
* Onboarding assistants
* Sales training
* Company policy assistants

---

# Core Features

## AI Conversation

Regular users interact with the system through a chat interface.

The AI generates responses based on:

* Base language model knowledge
* Company knowledge
* Uploaded training materials
* Retrieved internal information
* Human corrections and approved examples

The interface should support multi-turn conversations and conversation history.

---

## User Authentication

Users must authenticate before accessing the system.

Each account is assigned a permission level that determines which operations are available.

Example roles:

| Role          | Chat | Upload Materials | Correct AI | Delete Materials | Administration |
| ------------- | ---: | ---------------: | ---------: | ---------------: | -------------: |
| User          |    ✓ |                  |            |                  |                |
| Trainer       |    ✓ |                ✓ |          ✓ |                  |                |
| Manager       |    ✓ |                ✓ |          ✓ |                ✓ |                |
| Administrator |    ✓ |                ✓ |          ✓ |                ✓ |              ✓ |

The exact role and permission model should remain configurable.

Backend APIs must always validate permissions.

Frontend permission checks are used only for user experience and must never be treated as the security boundary.

---

# AI Correction

Authorized users can correct AI-generated responses.

Example:

```text
User:
Can a customer receive a refund immediately?

AI:
Yes. The customer can request a refund.

Trainer:
Incorrect.

Correct answer:
The company should first offer a replacement.
A refund should only be offered when replacement is not possible.
```

The correction can be stored as structured training data.

For example:

```json
{
  "question": "Can a customer receive a refund immediately?",
  "original_answer": "Yes. The customer can request a refund.",
  "corrected_answer": "The company should first offer a replacement. A refund should only be offered when replacement is not possible."
}
```

Corrections may later be used for:

* Knowledge retrieval
* Evaluation datasets
* Prompt examples
* Supervised fine-tuning
* LoRA / QLoRA training
* Preference datasets

This allows domain experts to improve the system without requiring them to understand machine learning.

---

# Training Material Library

Authorized users can upload internal company materials.

Supported material types may include:

### Documents

* PDF
* DOCX
* TXT
* Markdown
* CSV
* XLSX
* PPTX

### Images

* PNG
* JPEG
* WebP

Images may be processed by OCR or multimodal models.

### Audio

* MP3
* WAV
* M4A

Audio can be converted into text using speech-to-text models before entering the knowledge pipeline.

### Other Data

The architecture should allow additional data processors to be added in the future.

---

# Training Material Management

Uploaded materials are stored in a dedicated **Training Material Library**.

Authorized users can:

* Upload materials
* View materials
* Search materials
* Inspect processing status
* Reprocess materials
* Correct extracted information
* Delete materials

Deletion requires elevated permissions.

Example material lifecycle:

```text
Upload
  ↓
Validation
  ↓
File Storage
  ↓
Content Extraction
  ↓
Normalization
  ↓
Chunking
  ↓
Embedding
  ↓
Vector Database
  ↓
Available to AI
```

For audio:

```text
Audio
  ↓
Speech-to-Text
  ↓
Transcript
  ↓
Chunking
  ↓
Embedding
  ↓
Vector Database
```

For images:

```text
Image
  ↓
OCR / Vision Model
  ↓
Extracted Information
  ↓
Chunking
  ↓
Embedding
  ↓
Vector Database
```

Original files should be preserved unless explicitly deleted by an authorized user.

---

# Knowledge vs Model Training

Uploading a document should **not automatically trigger model fine-tuning**.

The platform separates company knowledge from model behavior.

## Knowledge Layer

Frequently changing company information should normally use Retrieval-Augmented Generation (RAG).

Examples:

* SOPs
* Product manuals
* Company policies
* Internal documentation
* Price information
* Customer service procedures

```text
Company Documents
        ↓
Content Extraction
        ↓
Embedding
        ↓
Vector Database
        ↓
Relevant Knowledge
        ↓
LLM
        ↓
Answer
```

This allows new information to become available without retraining the model.

## Training Layer

Fine-tuning should be used when the organization wants to modify model behavior.

Examples:

* Preferred response style
* Company terminology
* Customer service behavior
* Domain-specific reasoning patterns
* Repeated human corrections

Potential techniques include:

* LoRA
* QLoRA
* Supervised Fine-Tuning
* Preference Optimization

---

# System Architecture

The platform uses a frontend/backend architecture.

```text
              Smartphone
                  │
              Tablet
                  │
              Desktop
                  │
                  ▼
        ┌───────────────────┐
        │  React Frontend   │
        │ Responsive Web UI │
        └─────────┬─────────┘
                  │
               REST API
                  │
                  ▼
        ┌───────────────────┐
        │  Backend Server   │
        └─────────┬─────────┘
                  │
       ┌──────────┼──────────────┐
       │          │              │
       ▼          ▼              ▼
 Authentication Knowledge     AI Service
 & Permissions   Service
       │          │              │
       ▼          ▼              ▼
   Database   Vector DB      Model Server
                  │              │
                  │              ▼
                  │         Open-Weight LLM
                  │
                  ▼
            File Storage
```

---

# Frontend

The frontend is built with **React**.

The application must use responsive web design and support:

* Desktop browsers
* Laptop browsers
* Tablets
* Smartphones

The frontend should not directly access databases, model files, vector databases, or training files.

All operations must communicate with the backend through APIs.

```text
React
  │
  │ HTTPS / REST
  ▼
Backend API
```

The frontend is responsible for:

* Login
* Chat interface
* Conversation history
* Material upload
* Material management
* AI correction interface
* User management
* Permission-aware UI
* Administration interface

---

# Backend

The backend acts as the central control layer.

Responsibilities include:

* Authentication
* Authorization
* User management
* Session/token management
* Conversation management
* File management
* Material processing
* RAG
* Model inference
* Training data management
* AI correction management
* Model management
* Audit logging

The backend should expose all functionality through versioned APIs.

Example:

```text
/api/v1/
```

---

# API Design

The following endpoints represent the initial API structure.

## Authentication

```http
POST /api/v1/auth/login
POST /api/v1/auth/logout
POST /api/v1/auth/refresh
GET  /api/v1/auth/me
```

Example:

```http
POST /api/v1/auth/login
```

```json
{
  "username": "trainer@example.com",
  "password": "********"
}
```

---

# Users

```http
GET    /api/v1/users
GET    /api/v1/users/{id}
POST   /api/v1/users
PATCH  /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

Administrative permissions are required for user-management endpoints.

---

# Roles and Permissions

```http
GET   /api/v1/roles
POST  /api/v1/roles
PATCH /api/v1/roles/{id}

GET   /api/v1/permissions
```

Possible permissions include:

```text
chat.use

materials.read
materials.upload
materials.delete

corrections.create
corrections.read
corrections.delete

users.read
users.manage

models.read
models.manage

system.manage
```

This allows roles to be changed without redesigning the entire authorization system.

---

# Conversations

Create a conversation:

```http
POST /api/v1/conversations
```

List conversations:

```http
GET /api/v1/conversations
```

Read a conversation:

```http
GET /api/v1/conversations/{id}
```

Delete a conversation:

```http
DELETE /api/v1/conversations/{id}
```

---

# Chat

Send a message:

```http
POST /api/v1/conversations/{id}/messages
```

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
  "answer": "According to the company policy, replacement should be attempted before a refund.",
  "sources": [
    {
      "material_id": "mat_123",
      "name": "customer-service-policy.pdf"
    }
  ]
}
```

Streaming responses should also be supported.

Possible implementations include:

* Server-Sent Events
* WebSocket
* Streaming HTTP responses

---

# AI Corrections

Authorized trainers can correct a response.

```http
POST /api/v1/messages/{id}/corrections
```

Example:

```json
{
  "corrected_answer": "Replacement must be attempted before offering a refund.",
  "comment": "Updated according to the current customer service policy."
}
```

Retrieve corrections:

```http
GET /api/v1/corrections
```

Retrieve a specific correction:

```http
GET /api/v1/corrections/{id}
```

Delete a correction:

```http
DELETE /api/v1/corrections/{id}
```

Deletion should require elevated permissions.

---

# Training Materials

Upload:

```http
POST /api/v1/materials
```

List:

```http
GET /api/v1/materials
```

Read metadata:

```http
GET /api/v1/materials/{id}
```

Delete:

```http
DELETE /api/v1/materials/{id}
```

Reprocess:

```http
POST /api/v1/materials/{id}/process
```

Processing status:

```http
GET /api/v1/materials/{id}/status
```

Example status:

```json
{
  "id": "mat_123",
  "filename": "employee-handbook.pdf",
  "status": "ready",
  "processing": {
    "extraction": "completed",
    "chunking": "completed",
    "embedding": "completed"
  }
}
```

---

# Model API

The backend should provide an abstraction between the application and the actual language model.

```text
Application
     │
     ▼
Model Service Interface
     │
 ┌───┼─────────────┐
 ▼   ▼             ▼
Gemma gpt-oss    Other Models
```

The application should therefore avoid model-specific logic whenever possible.

Example endpoints:

```http
GET  /api/v1/models
GET  /api/v1/models/current
POST /api/v1/models/{id}/activate
```

Model-management APIs should require administrative permissions.

---

# Model Strategy

The platform is designed to support open-weight language models.

Potential model families include:

* Gemma
* gpt-oss
* Qwen
* Llama
* Other compatible open-weight models

The initial implementation should select a lightweight default model while maintaining a model abstraction layer.

This prevents vendor lock-in and allows organizations to select models based on:

* Available RAM
* Available VRAM
* CPU performance
* GPU availability
* Language requirements
* Model quality
* Licensing requirements

---

# Deployment

The primary deployment target is an **internal company VM**.

Example:

```text
Company Network

Employees
    │
    │ Browser
    ▼
https://internal-ai.company.local
    │
    ▼
┌─────────────────────────────┐
│        Internal VM          │
│                             │
│ React Frontend              │
│ Backend API                 │
│ Database                    │
│ Vector Database             │
│ File Storage                │
│ Model Server                │
│ Open-Weight LLM             │
│                             │
└─────────────────────────────┘
```

The system should not require public Internet access for normal AI conversations after installation and model preparation.

This allows sensitive company information to remain inside the organization's infrastructure.

---

# Container Deployment

Docker should be the preferred deployment method.

The long-term goal is to make installation as simple as:

```bash
docker compose up -d
```

Example services:

```yaml
services:

  frontend:
    # React web application

  backend:
    # API server

  database:
    # Application database

  vector-db:
    # Knowledge embeddings

  model:
    # LLM inference server

  worker:
    # Document / image / audio processing
```

Individual components should remain replaceable.

---

# Security

Because the platform may process sensitive company information, security must be considered a core feature.

The system should support:

* Local deployment
* HTTPS
* Authentication
* Role-based access control
* API authorization
* Secure password hashing
* File access control
* Audit logs
* Upload validation
* File size limits
* File type validation
* Model access isolation
* Training-data access control
* Secure deletion procedures

Every privileged backend operation must independently validate the authenticated user's permissions.

---

# Audit Logging

Important administrative and training actions should be recorded.

Examples:

```text
USER_LOGIN

MATERIAL_UPLOAD
MATERIAL_DELETE

AI_CORRECTION_CREATE
AI_CORRECTION_DELETE

USER_CREATE
USER_UPDATE
USER_DELETE

ROLE_UPDATE

MODEL_ACTIVATE
MODEL_TRAIN
```

Example:

```json
{
  "event": "MATERIAL_DELETE",
  "user_id": "usr_102",
  "material_id": "mat_501",
  "timestamp": "2026-08-24T15:30:00Z"
}
```

Audit logs help organizations understand who modified the AI's knowledge and training data.

---

# Design Principles

## Privacy First

Company information should remain inside the company's infrastructure whenever possible.

## Human-Guided AI

AI behavior should be improvable by authorized domain experts rather than only ML engineers.

## Lightweight

The platform should remain usable by SMEs without requiring enterprise-scale GPU infrastructure.

## API First

All application functionality should be available through backend APIs.

## Responsive

The same React application should work across desktop and mobile devices.

## Model Independent

Business logic should not depend on one specific LLM.

## Knowledge Is Not Training

Frequently changing information should use RAG.

Model training should be reserved for behavioral and domain adaptation.

## Security by Backend Enforcement

Frontend controls improve usability.

Backend authorization provides security.

---

# Initial Product Scope

The first version should focus on:

1. Authentication
2. Role-based permissions
3. Responsive React frontend
4. AI chat
5. Conversation history
6. Training material upload
7. Training material management
8. Document processing
9. Image processing
10. Audio transcription
11. RAG
12. AI response correction
13. Correction dataset management
14. Model abstraction
15. Local VM deployment
16. Docker-based installation
17. Audit logging

Advanced model training can be introduced after the knowledge and correction pipelines are stable.

---

# Future Development

Potential future features include:

* LoRA / QLoRA training
* Automatic training dataset generation
* Training approval workflows
* Model version management
* Model rollback
* AI evaluation
* Knowledge versioning
* Material categories
* Department-specific knowledge bases
* Department-specific permissions
* Multiple AI assistants
* Multiple models
* Multimodal conversations
* Speech conversations
* API integrations
* Webhooks
* SSO
* LDAP / Active Directory integration
* Backup and restore
* High availability
* GPU cluster support

---

# Project Vision

The goal is not simply to provide another chatbot.

The platform aims to give small and medium-sized organizations the ability to build and continuously improve their own private AI.

```text
Company Knowledge
       +
Internal Documents
       +
Images & Audio
       +
Expert Corrections
       +
Open-Weight LLM
       ↓
────────────────────
   Company Private AI
────────────────────
```

The organization controls:

**its data, its knowledge, its training process, and its AI.**

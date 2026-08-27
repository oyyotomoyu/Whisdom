# Architecture

Whisdom is a self-hosted private AI platform for businesses that want to build an internal AI assistant using their own documents, media, and human corrections.

The system is designed around four core goals:

* **Private deployment** - company data stays inside the organization's own server or private infrastructure.
* **Simple operation** - small teams should be able to run the platform without deep AI infrastructure knowledge.
* **Expandable knowledge** - documents, media, and corrections can continuously improve the assistant.
* **Model flexibility** - the system should support local open models and avoid permanent dependency on one provider.

---

# Technology Stack

## UI

The frontend is built with **React**.

The UI must support:

* Mobile browsers
* Tablet browsers
* Desktop browsers
* Responsive layouts for both small and large screens

The frontend communicates only with the backend through REST APIs. It must not directly access the database, model runtime, vector database, or file storage.

## Server

The backend is a **Go web server**.

The server is responsible for:

* Account authentication
* Permission checks
* REST API routing
* Material upload and management
* Knowledge processing
* RAG orchestration
* AI inference requests
* Training and correction workflows
* Audit logging

## Model

The default local model target is **Gemma 4**.

The model layer should be abstracted so the platform can later support different model backends, model sizes, or inference runtimes without rewriting the whole application.

Possible model responsibilities include:

* Chat response generation
* Document-aware question answering
* Summarization
* Classification
* Human correction usage
* Optional fine-tuning or adapter training

---

# System Structure

```text
Mobile / PC
    ↓
React Frontend
    ↓
REST API
    ↓
Backend Server
    ↓
┌──────┼─────────┐
│      │         │
帳號   素材庫     AI / RAG
權限   文件媒體   訓練 / 推論
```

## Main Areas

### 帳號 / 權限

The account and permission system controls who can access each function.

Responsibilities:

* Login and logout
* User identity
* Role management
* Permission validation
* Session or token handling
* API-level authorization

Permission checks must always happen on the backend. Frontend checks are only for user experience.

### 素材庫 / 文件媒體

The material library stores company knowledge sources.

Supported material types may include:

* PDF
* DOCX
* TXT
* Markdown
* CSV
* XLSX
* PPTX
* Images
* Audio files

Material processing may include:

* File validation
* File storage
* Text extraction
* OCR for images
* Speech-to-text for audio
* Normalization
* Chunking
* Embedding
* Indexing into a vector database

### AI / RAG / 訓練 / 推論

The AI layer combines the local model with business knowledge.

RAG should be used for frequently changing information such as:

* SOPs
* Product manuals
* Company policies
* Customer support procedures
* Internal documentation

Training or fine-tuning should be used only when the company wants to change model behavior, such as:

* Preferred answer style
* Company terminology
* Repeated correction patterns
* Domain-specific response behavior

---

# Minimum Server Hardware

The platform should define a lowest practical server PC target that can run the web server, database, file processing, vector search, and local Gemma 4 inference.

## Minimum Practical Target

```text
CPU:      8-core modern x86_64 or ARM64 processor
RAM:      32 GB
GPU:      NVIDIA GPU with 12 GB VRAM
Storage:  1 TB NVMe SSD
Network:  Gigabit Ethernet
OS:       Linux server distribution
```

This target is intended for small teams, light concurrent usage, and quantized local inference.

## Recommended Target

```text
CPU:      12-core or better modern processor
RAM:      64 GB
GPU:      NVIDIA GPU with 16 GB or more VRAM
Storage:  2 TB NVMe SSD
Network:  Gigabit Ethernet or better
OS:       Linux server distribution
```

The recommended target gives more room for larger documents, more users, faster indexing, and smoother AI responses.

## Notes

Actual hardware requirements depend on:

* Gemma 4 model size
* Quantization level
* Context length
* Number of concurrent users
* Document volume
* Embedding model choice
* Whether training runs on the same machine

For the lowest-cost deployment, training and heavy indexing can be scheduled during off-hours, while regular chat inference remains available during business hours.

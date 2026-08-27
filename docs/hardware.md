# Hardware

Whisdom is designed to run as a self-hosted private AI platform on a company-owned server PC or private VM. The hardware target should be low enough for small businesses to afford, but strong enough to run the Go server, React frontend hosting, database, file processing, vector search, and local model inference.

The default local model target is **Gemma 4**.

---

# Minimum Practical Server PC

This is the lowest practical target for a small team using quantized local inference and light concurrent usage.

```text
CPU:      8-core modern x86_64 or ARM64 processor
RAM:      32 GB
GPU:      NVIDIA GPU with 12 GB VRAM
Storage:  1 TB NVMe SSD
Network:  Gigabit Ethernet
OS:       Linux server distribution
```

Expected use:

* Small business deployment
* Light concurrent chat usage
* Small to medium document library
* Quantized Gemma 4 inference
* RAG-based knowledge retrieval
* Training or heavy processing scheduled during off-hours

This target should avoid running large training jobs while many users are actively chatting.

---

# Recommended Server PC

This target gives the platform more room for smoother AI responses, larger document libraries, and more concurrent usage.

```text
CPU:      12-core or better modern x86_64 processor
RAM:      64 GB
GPU:      NVIDIA GPU with 16 GB or more VRAM
Storage:  2 TB NVMe SSD
Network:  Gigabit Ethernet or better
OS:       Linux server distribution
```

Expected use:

* Small to medium business deployment
* More concurrent users
* Larger material library
* Faster document processing
* Better model response latency
* More comfortable RAG and embedding workloads
* Limited training or adapter jobs

---

# Higher Performance Server

For heavier workloads, use a stronger GPU and more memory.

```text
CPU:      16-core or better processor
RAM:      128 GB
GPU:      NVIDIA GPU with 24 GB or more VRAM
Storage:  4 TB NVMe SSD or larger
Network:  2.5 GbE / 10 GbE
OS:       Linux server distribution
```

Expected use:

* Larger number of users
* Large internal document collections
* Faster embedding and indexing
* Longer model context
* Larger local model variants
* Training jobs on the same machine

---

# Hardware Responsibilities

Whisdom server hardware must support several workloads at the same time:

* Go backend API server
* Static frontend hosting
* User database
* File storage
* Material processing
* OCR or speech-to-text when enabled
* Embedding generation
* Vector database search
* RAG context assembly
* Gemma 4 inference
* Log writing and log search
* Optional training or adapter jobs

The AI model is normally the heaviest part of the system. GPU VRAM has the biggest effect on model size, context length, and response speed.

---

# CPU

Minimum:

```text
8 cores
```

Recommended:

```text
12 cores or more
```

CPU is used for:

* API handling
* Authentication
* File processing
* Text extraction
* Database work
* Vector database coordination
* Background jobs
* Compression and cleanup

If OCR, audio transcription, or CPU-based embedding is used, more CPU cores will help.

---

# Memory

Minimum:

```text
32 GB RAM
```

Recommended:

```text
64 GB RAM
```

Higher performance:

```text
128 GB RAM
```

RAM is used for:

* Operating system
* Go backend
* Database
* Vector database
* File processing
* Embedding batches
* Model runtime buffers
* Concurrent requests

Running the model, vector database, and document processing on the same machine requires enough RAM headroom. The server should avoid using swap during normal chat usage.

---

# GPU

Minimum:

```text
NVIDIA GPU with 12 GB VRAM
```

Recommended:

```text
NVIDIA GPU with 16 GB or more VRAM
```

Higher performance:

```text
NVIDIA GPU with 24 GB or more VRAM
```

GPU is used for:

* Local Gemma 4 inference
* Faster embedding generation when supported
* Optional fine-tuning or adapter training

VRAM affects:

* Model size
* Quantization choice
* Context length
* Response speed
* Number of concurrent model requests
* Whether training can run locally

If the GPU has limited VRAM, use quantized inference and keep training jobs separate from live chat usage.

---

# Storage

Minimum:

```text
1 TB NVMe SSD
```

Recommended:

```text
2 TB NVMe SSD
```

Storage is used for:

* Original uploaded files
* Extracted text
* Chunked material
* Embeddings
* Vector database files
* Conversation history
* Correction datasets
* Training material path
* Model files
* Logs
* Backups

NVMe SSD is strongly preferred because material processing, vector indexing, and model loading can be slow on HDD storage.

Recommended storage layout:

```text
/opt/whisdom/
├── data/
│   ├── database/
│   ├── vector/
│   ├── materials/
│   ├── extracted/
│   ├── corrections/
│   └── training/
├── models/
├── log/
└── backup/
```

The training material path stored in server config should point to a controlled location such as:

```text
/opt/whisdom/data/training/
```

The server must validate and normalize this path before using it.

---

# Network

Minimum:

```text
Gigabit Ethernet
```

Recommended:

```text
Gigabit Ethernet or better
```

Network is used for:

* Browser access from mobile and PC
* Uploading training materials
* Admin access
* Backups to another machine
* Optional updates or external integrations

For a private company deployment, the server should normally run inside the company's trusted network or behind a secure VPN.

---

# Operating System

Recommended:

```text
Linux server distribution
```

Examples:

* Ubuntu Server
* Debian
* Rocky Linux

The OS should support:

* Go server runtime
* NVIDIA GPU drivers
* CUDA-compatible model runtime when needed
* System service management
* Firewall configuration
* Scheduled backups
* Log rotation or monitoring

---

# Deployment Sizes

## Small Team

```text
Users:      1-10
Documents:  Small to medium
Hardware:   Minimum practical server PC
```

Use RAG for knowledge. Schedule heavy indexing and training jobs outside work hours.

## Growing Team

```text
Users:      10-50
Documents:  Medium to large
Hardware:   Recommended server PC
```

Use background queues for material processing. Monitor GPU memory, response latency, and storage growth.

## Larger Internal Deployment

```text
Users:      50+
Documents:  Large
Hardware:   Higher performance server or split services
```

Consider separating:

* API server
* Database
* Vector database
* Model inference server
* Training machine
* Backup storage

---

# Scaling Options

Whisdom should start as a single-server deployment, but the architecture should allow services to split later.

Possible split:

```text
React Frontend + Go API
        ↓
Database Server
        ↓
Vector Database
        ↓
Model Inference Server
        ↓
Training Server
```

Single-server deployment is simpler and cheaper. Multi-server deployment improves performance and reliability when the company grows.

---

# Backup Requirements

The server should back up:

* User database
* Conversation history
* Uploaded materials
* Extracted text
* Vector index or rebuild metadata
* Corrections
* Config
* Logs within retention policy

Backups should not include:

* Plaintext passwords
* Temporary upload files
* Expired tokens
* Unneeded cache files

Backups should be encrypted when stored outside the server.

---

# Monitoring

The server should monitor:

* CPU usage
* RAM usage
* GPU utilization
* GPU VRAM usage
* Disk usage
* API latency
* Model response latency
* Failed logins
* Material processing failures
* Training job failures
* Log write failures

The log system is important for operational debugging and audit history, but hardware-level monitoring should also be available.

---

# Hardware Decision Rule

Choose hardware based on:

* Gemma 4 model size
* Quantization level
* Expected concurrent users
* Document volume
* Required response speed
* Whether training runs on the same server
* Whether OCR or audio transcription is enabled

When uncertain, prioritize:

1. More GPU VRAM
2. More RAM
3. NVMe SSD capacity
4. More CPU cores

GPU VRAM is usually the first limit for local AI inference.

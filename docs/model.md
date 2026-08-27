# Model

Whisdom uses a local language model to answer company questions with private business knowledge. The default model target is **Gemma 4**, running inside the company's own infrastructure.

The model is not accessed directly by the React frontend. All model requests go through the Go backend server.

---

# Core Idea

The model layer has two separate responsibilities:

* **Inference** - answer user questions during conversations.
* **Training or improvement** - use approved material and corrections to improve future behavior.

Uploading a document should not automatically fine-tune the model. Most company knowledge should first be used through RAG.

---

# Server-Controlled Model Access

The frontend sends a user message to the backend.

The backend then:

* Validates the user
* Checks permissions
* Stores the message
* Finds relevant company knowledge
* Builds the model prompt
* Sends the request to the model runtime
* Receives the model answer
* Stores the answer
* Returns the answer to the frontend
* Logs important events and failures

Flow:

```text
React Frontend
    ↓
POST /api/v1/conversations/{id}/messages
    ↓
Go Backend Server
    ↓
RAG Search
    ↓
Prompt Builder
    ↓
Gemma 4 Runtime
    ↓
Answer
    ↓
Backend stores and returns response
```

---

# Model Interface

The backend should use an internal interface for model access.

This keeps the platform model-independent and allows future replacement of Gemma 4 or the runtime.

Example interface shape:

```go
type ModelClient interface {
    Generate(ctx context.Context, request GenerateRequest) (*GenerateResponse, error)
}
```

Request fields may include:

* Conversation ID
* User message
* System prompt
* Retrieved context
* Previous conversation messages
* Temperature
* Max tokens
* Timeout

Response fields may include:

* Answer text
* Model name
* Token usage
* Finish reason
* Runtime duration
* Error details when failed

---

# Default Model

Default target:

```text
Gemma 4
```

The exact model size and quantization level should be configurable because hardware capacity will differ between deployments.

Config examples:

```json
{
  "model": {
    "provider": "local",
    "name": "gemma-4",
    "runtime_url": "http://127.0.0.1:11434",
    "context_length": 8192,
    "temperature": 0.2,
    "max_output_tokens": 1024
  }
}
```

The server must not hard-code model runtime details inside API handlers.

---

# RAG Layer

RAG is the first choice for company knowledge.

Use RAG for:

* SOPs
* Product manuals
* Company policies
* Customer support rules
* Internal documentation
* Frequently changing information

RAG flow:

```text
User question
    ↓
Embedding search
    ↓
Relevant chunks
    ↓
Prompt context
    ↓
Gemma 4
    ↓
Answer with sources
```

The RAG layer should return source references so the UI can show where the answer came from.

Source reference data:

```json
{
  "material_id": "mat_123",
  "title": "Refund SOP",
  "chunk_id": "chunk_456",
  "score": 0.87
}
```

The server should avoid sending unrelated chunks to the model. More context is not always better.

---

# Training Layer

Training should be separate from normal upload.

Use training when the company wants to change model behavior, such as:

* Preferred answer style
* Company terminology
* Repeated correction patterns
* Domain-specific response behavior
* Classification behavior

Possible techniques:

* Prompt examples
* Evaluation datasets
* Supervised fine-tuning
* LoRA
* QLoRA
* Preference datasets

Training jobs should be admin-only and logged.

Training job flow:

```text
Approved materials and corrections
    ↓
Dataset preparation
    ↓
Training or adapter job
    ↓
Validation / evaluation
    ↓
Admin approval
    ↓
Model or adapter activation
```

---

# Correct Answer Data

Admin users can teach the model correct answers when the AI response is wrong.

Correction record:

```json
{
  "question": "Can a customer receive a refund immediately?",
  "original_answer": "Yes. The customer can request a refund.",
  "corrected_answer": "The company should first offer a replacement. A refund should only be offered when replacement is not possible.",
  "usage": ["rag", "evaluation", "training"],
  "created_by": "user_123"
}
```

Corrections can be used for:

* RAG retrieval
* Prompt examples
* Evaluation tests
* Fine-tuning datasets
* Preference datasets

Correction creation, update, deletion, and training usage must be logged.

---

# Material Processing

Training and RAG materials come from the material library.

Processing flow:

```text
Uploaded file
    ↓
Validation
    ↓
File storage
    ↓
Text extraction
    ↓
Normalization
    ↓
Chunking
    ↓
Embedding
    ↓
Vector database
```

Material types:

* PDF
* DOCX
* TXT
* Markdown
* CSV
* XLSX
* PPTX
* Images with OCR
* Audio with speech-to-text

Original files should be preserved unless deleted by an authorized admin.

---

# Training Material Path

The backend config stores where model training material is saved or sent.

Example:

```json
{
  "training_material_path": "/opt/whisdom/data/training/"
}
```

This path is used by:

* Material upload workflows
* Dataset preparation
* Training job input
* Admin destination path settings

Rules:

* Only admin users can edit the path.
* Path changes must be logged.
* Path must be validated before saving.
* Path traversal must be blocked.
* The server must normalize paths before use.

---

# Prompt Construction

The backend builds prompts for the model.

Prompt inputs:

* System instruction
* User question
* Recent conversation messages
* Retrieved RAG context
* Relevant corrections
* Role or permission context when needed

Prompt rules:

* Keep company knowledge separate from system instructions.
* Do not include unrelated documents.
* Do not include secrets, tokens, or hidden admin config.
* Prefer concise context over excessive context.
* Include source metadata outside the answer when possible.

---

# Model Safety And Privacy

The model must run under server control.

Requirements:

* Do not expose model runtime directly to the browser.
* Do not send private company data to external services unless explicitly configured.
* Do not log full prompts by default.
* Do not log full private documents.
* Do not log secrets or tokens.
* Apply request timeouts.
* Apply maximum context and output limits.
* Validate uploaded material before processing.

---

# Model Config

Model behavior should be configurable.

Possible config:

```json
{
  "model": {
    "provider": "local",
    "name": "gemma-4",
    "runtime_url": "http://127.0.0.1:11434",
    "context_length": 8192,
    "temperature": 0.2,
    "max_output_tokens": 1024,
    "request_timeout_seconds": 120
  },
  "rag": {
    "top_k": 5,
    "min_score": 0.65,
    "include_sources": true
  },
  "training": {
    "training_material_path": "/opt/whisdom/data/training/",
    "allow_runtime_training": false
  }
}
```

Only users with model or system management permission can edit model config.

---

# Logging Requirements

The model layer must log important events through the server log system.

Log these events:

* Model runtime connection failure
* Model request timeout
* Model answer generation failure
* RAG retrieval failure
* Material embedding failure
* Training job start
* Training job success
* Training job failure
* Model config change
* Training material path change
* Correction used for training

Do not log:

* Full prompts
* Full answers when they may contain private data
* Full document contents
* Secrets
* Tokens
* Passwords

Log content should be short and operational, for example:

```text
model request failed conversation=conv_123
training job started job=train_123
updated training material path
```

---

# First Implementation Order

1. Define model config structure.
2. Define `ModelClient` interface.
3. Implement local runtime client.
4. Add conversation message API integration.
5. Add basic prompt builder.
6. Add RAG retrieval interface.
7. Add source references to chat responses.
8. Add correction data model.
9. Add training material path config.
10. Add training job model.
11. Add model and training logs.
12. Add tests for prompt building, config validation, and model error handling.

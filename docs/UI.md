# UI

Whisdom uses a React frontend. The main user experience is a conversation workspace similar to common GPT applications: users open the app, see recent conversations, select or start a chat, and ask questions against company knowledge.

Admin-only tools are available only to authorized users. Regular users should not see material upload, model teaching, or training path configuration in the UI.

---

# Reference Pattern

The UI structure should follow the important frontend pattern used in the NexGestion project:

```text
client/src/views
client/src/requests
client/src/locales
```

The same separation should be used in Whisdom:

* `views` contains route-level pages and screen composition.
* `requests` contains typed API request modules.
* `locales` contains translation setup and language keys.

Views must not call `fetch` directly. API access belongs in `requests`. User-visible text must come from `locales`.

---

# Main Experience

## Conversation First

The default authenticated screen is the conversation area.

The UI should feel like a private company GPT app:

```text
┌────────────────────────────────────────────┐
│ Sidebar                                    │
│ - New chat          Conversation Area      │
│ - Recent chats      ┌───────────────────┐  │
│ - Search            │ Messages          │  │
│ - Admin tools       │                   │  │
│                     │ Input composer    │  │
└────────────────────────────────────────────┘
```

The conversation area should support:

* New conversation
* Recent conversation list
* Search recent conversations
* Conversation title
* Multi-turn messages
* Streaming or loading response state
* Message retry
* Copy answer
* Source references when RAG is used
* Clear error state when the model or backend fails

The input composer should support:

* Text input
* Submit button
* Keyboard submit behavior
* Disabled state while sending
* Optional file attachment in the future

## Recent Conversations

Recent conversations are visible in the main sidebar or mobile drawer.

Each recent conversation item should show:

* Conversation title
* Last updated time
* Active state
* Delete or archive action when allowed

The most recent conversations should be easy to find without opening an admin page.

---

# User Roles

## Regular User

Regular users can:

* Chat with the AI
* View their recent conversations
* Continue previous conversations
* View answer sources when available

Regular users cannot:

* Upload training materials
* Teach the model corrected answers
* Edit material destination paths
* Manage users or permissions
* Trigger training or reprocessing jobs

## Admin User

Admin users can access all regular user functions plus admin tools.

Admin-only UI areas include:

* User and role management
* Training material upload
* Training material library
* Correct answer teaching
* Material destination path editing
* RAG and training status
* Model or system settings

Admin navigation must be hidden from users without permission, but backend API permission checks remain the real security boundary.

---

# Admin UI

## Training Material Upload

Admins can upload company knowledge and training materials.

The upload UI should support:

* Drag-and-drop upload
* File picker
* Multiple files
* Upload progress
* Validation errors
* Processing status
* Material type detection
* Destination path selection

Supported material examples:

* PDF
* DOCX
* TXT
* Markdown
* CSV
* XLSX
* PPTX
* Images
* Audio

## Material Library

Admins can inspect uploaded material.

The material library should show:

* File name
* File type
* Upload time
* Uploaded by
* Processing status
* Destination path
* RAG availability
* Training availability
* Reprocess action
* Delete action

## Destination Path Editing

Admins can edit where uploaded material is sent.

This path represents the model training or knowledge processing destination.

Examples:

```text
/training/company-policy/
/training/customer-support/
/training/product-manuals/
/rag/sop/
/rag/internal-docs/
```

The UI should validate destination paths before saving.

Validation rules:

* Path is required.
* Path must start with `/`.
* Path cannot contain `..`.
* Path should use clear domain names.
* Path changes must be audit logged.

## Teach Correct Answer

Admins can teach the language model the correct answer when an AI response is wrong.

The correction UI should allow an admin to record:

* Original user question
* Original AI answer
* Correct answer
* Reason or note
* Related material or source
* Whether the correction is used for RAG, evaluation, training, or all

Example correction flow:

```text
User asks a question
    ↓
AI gives an incorrect answer
    ↓
Admin opens correction panel
    ↓
Admin writes the correct answer
    ↓
System stores correction data
    ↓
Correction becomes available for RAG, evaluation, or training
```

The correction action may appear inside a conversation message only for admins.

---

# Recommended Client Structure

```text
client/
├── src/
│   ├── components/
│   ├── hooks/
│   ├── layouts/
│   │   ├── AppLayout/
│   │   └── AuthLayout/
│   ├── locales/
│   │   ├── i18n.ts
│   │   ├── languages.ts
│   │   └── key/
│   │       ├── Global/
│   │       ├── Login/
│   │       ├── Conversation/
│   │       └── Admin/
│   ├── requests/
│   │   ├── core/
│   │   ├── auth/
│   │   ├── conversations/
│   │   ├── materials/
│   │   ├── corrections/
│   │   ├── users/
│   │   └── system/
│   ├── store/
│   ├── theme/
│   └── views/
│       ├── Login/
│       ├── Conversation/
│       ├── Materials/
│       ├── Corrections/
│       ├── Settings/
│       └── index.tsx
```

---

# Views

## Login

Unauthenticated users start at login.

The login view should support:

* Username or email
* Password
* Error state
* Loading state
* Language selector

## Conversation

This is the main application view.

Route:

```text
/chat
/chat/:conversationId
```

The conversation view owns:

* Recent conversation sidebar
* Current message thread
* Message composer
* Source reference display
* Admin-only correction action on AI messages

## Materials

Admin-only.

Route:

```text
/admin/materials
```

The materials view owns:

* Upload panel
* Material table or list
* Processing status
* Destination path editor
* Delete and reprocess actions

## Corrections

Admin-only.

Route:

```text
/admin/corrections
```

The corrections view owns:

* Correction list
* Create correction form
* Edit correction form
* Usage target selection
* Review status

## Settings

Admin-only, except for personal profile settings if needed.

Route:

```text
/settings
```

The settings view owns:

* User management
* Role management
* Permission management
* System configuration
* Model configuration

---

# Requests

Request modules should centralize API access.

Recommended modules:

```text
requests/auth
requests/conversations
requests/materials
requests/corrections
requests/users
requests/roles
requests/system
requests/core
```

Guidelines:

* Views should call typed request functions.
* Token refresh and request errors belong in `requests/core`.
* Uploads should use `FormData`.
* Conversation APIs should support streaming or polling when needed.
* Admin APIs must return permission errors clearly.

---

# Locales

Whisdom should use locale files from the start.

Recommended languages:

* English
* Traditional Chinese
* Japanese

Recommended key groups:

```text
locales/key/Global
locales/key/Login
locales/key/Conversation
locales/key/Admin
locales/key/Materials
locales/key/Corrections
locales/key/Settings
```

All visible UI text should come from translation keys, including:

* Navigation labels
* Buttons
* Placeholders
* Empty states
* Error messages
* Accessibility labels
* Document titles

---

# Responsive Layout

The UI must support mobile and PC.

## Desktop

Desktop layout should use:

* Persistent sidebar
* Recent conversation list
* Large conversation panel
* Admin navigation when allowed

## Mobile

Mobile layout should use:

* Top bar
* Navigation drawer
* Full-width conversation view
* Collapsible recent conversation list
* Touch-friendly controls

The app must be usable at 320px width without document-level horizontal scrolling.

---

# Permission Rules

UI permission rules:

* Hide admin navigation from non-admin users.
* Hide correction actions from non-admin users.
* Hide material upload and destination path editing from non-admin users.
* Disable actions while permissions or session state are loading.
* Show a clear unauthorized page if a user opens an admin route directly.

Backend permission rules:

* Every admin API must validate permissions.
* Frontend checks are not security boundaries.
* Permission failures should return consistent API errors.

---

# First UI Build Order

1. Login view
2. Auth provider and protected route
3. App layout with responsive navigation
4. Conversation view
5. Recent conversation sidebar
6. Conversation request module
7. Admin-only materials view
8. Admin-only correction flow
9. Destination path editor
10. Locale keys for all visible text

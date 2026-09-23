package system

// CreateConversation starts a new, empty conversation owned by userID.
func (s *Store) CreateConversation(userID string) *Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := timeNow()
	conv := &Conversation{
		ID:        NewID("conv"),
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.conversations[conv.ID] = conv
	return conv
}

// ListConversations returns userID's conversations, most recently updated
// first.
func (s *Store) ListConversations(userID string) []*Conversation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*Conversation, 0)
	for _, c := range s.conversations {
		if c.UserID == userID {
			out = append(out, c)
		}
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].UpdatedAt.After(out[j-1].UpdatedAt); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// GetConversation returns a conversation, scoped to its owner.
func (s *Store) GetConversation(userID, id string) (*Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conv, ok := s.conversations[id]
	if !ok || conv.UserID != userID {
		return nil, ErrNotFound
	}
	return conv, nil
}

// DeleteConversation removes a conversation, scoped to its owner.
func (s *Store) DeleteConversation(userID, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[id]
	if !ok || conv.UserID != userID {
		return ErrNotFound
	}
	delete(s.conversations, id)
	return nil
}

// AppendMessage adds a message to a conversation, scoped to its owner, and
// bumps its UpdatedAt. If the conversation has no title yet and this is the
// first user message, the title is derived from it.
func (s *Store) AppendMessage(userID, conversationID string, role MessageRole, content string, sources []MessageSource) (*Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[conversationID]
	if !ok || conv.UserID != userID {
		return nil, ErrNotFound
	}

	msg := Message{
		ID:             NewID("msg"),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		Sources:        sources,
		CreatedAt:      timeNow(),
	}
	conv.Messages = append(conv.Messages, msg)
	conv.UpdatedAt = msg.CreatedAt

	if conv.Title == "" && role == MessageRoleUser {
		conv.Title = deriveTitle(content)
	}

	return &msg, nil
}

// SetConversationTitle overwrites a conversation's title, scoped to its
// owner. Used to replace the naive truncated title with a model-summarized
// one once the first exchange completes.
func (s *Store) SetConversationTitle(userID, id, title string) (*Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conv, ok := s.conversations[id]
	if !ok || conv.UserID != userID {
		return nil, ErrNotFound
	}
	conv.Title = title
	return conv, nil
}

func deriveTitle(content string) string {
	const maxLen = 60
	runes := []rune(content)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "…"
	}
	return content
}

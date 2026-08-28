export default {
  title: "AI Corrections",
  subtitle: "Teach the model the correct answer when a response is wrong.",
  new: "New correction",
  form: {
    question: "Original question",
    originalAnswer: "Original AI answer",
    correctedAnswer: "Correct answer",
    note: "Reason or note",
    relatedMaterial: "Related material or source",
    usage: "Use this correction for",
    usageOptions: {
      rag: "RAG",
      evaluation: "Evaluation",
      training: "Training",
      all: "All",
    },
    submit: "Save correction",
    submitting: "Saving...",
  },
  list: {
    empty: "No corrections recorded yet.",
    createdAt: "Created",
    status: "Review status",
  },
  status: {
    pending: "Pending review",
    approved: "Approved",
    rejected: "Rejected",
  },
  deleteConfirm: "Delete this correction? This cannot be undone.",
};

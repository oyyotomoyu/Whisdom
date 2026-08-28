export default {
  title: "AI 答案修正",
  subtitle: "當回答錯誤時，教導模型正確的答案。",
  new: "新增修正",
  form: {
    question: "原始問題",
    originalAnswer: "原始 AI 回答",
    correctedAnswer: "正確答案",
    note: "原因或備註",
    relatedMaterial: "相關資料或來源",
    usage: "此修正用於",
    usageOptions: {
      rag: "RAG",
      evaluation: "評估",
      training: "訓練",
      all: "全部",
    },
    submit: "儲存修正",
    submitting: "儲存中...",
  },
  list: {
    empty: "尚無修正紀錄。",
    createdAt: "建立時間",
    status: "審核狀態",
  },
  status: {
    pending: "審核中",
    approved: "已核准",
    rejected: "已拒絕",
  },
  deleteConfirm: "確定要刪除此修正嗎？此操作無法復原。",
};

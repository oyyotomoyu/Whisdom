export default {
  title: "AI 修正",
  subtitle: "回答が誤っている場合に、正しい答えをモデルに教えます。",
  new: "新しい修正",
  form: {
    question: "元の質問",
    originalAnswer: "元の AI 回答",
    correctedAnswer: "正しい回答",
    note: "理由・メモ",
    relatedMaterial: "関連資料・出典",
    usage: "この修正の用途",
    usageOptions: {
      rag: "RAG",
      evaluation: "評価",
      training: "トレーニング",
      all: "すべて",
    },
    submit: "修正を保存",
    submitting: "保存中...",
  },
  list: {
    empty: "まだ修正はありません。",
    createdAt: "作成日時",
    status: "レビュー状況",
  },
  status: {
    pending: "レビュー待ち",
    approved: "承認済み",
    rejected: "却下",
  },
  deleteConfirm: "この修正を削除しますか？この操作は取り消せません。",
};

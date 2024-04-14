import React, { createContext, useContext, useReducer, ReactNode, useEffect } from 'react';
import { getCommentList } from './api';

// 定义评论类型
export type CommentType = {
  id: string;
  nickname: string;
  avatar: string;
  website?: string;
  content: string;
  comment_at: string;
  reply_user?: string;
};

export type CommentListItemType = {
  comment: CommentType;
  replies?: CommentType[];
}

// 定义状态类型
export type CommentsStateType = {
  comments: CommentListItemType[];
  article_id: string;
};

// 定义action类型
type Action =
  | { type: 'SET_COMMENTS'; payload: CommentListItemType[] };

function getArticleID (): string {
  const path = window.location.pathname; // 获得路径部分，例如: "/archives/35"
  const segments = path.split('/');      // 将路径分割为段

  // 检查是否存在第三个段并且是否为数字
  if (segments.length > 2 && /^\d+$/.test(segments[2])) {
    return segments[2];
  }

  return '';
}

// 初始化状态
const initialState: CommentsStateType = {
  comments: [],
  article_id: getArticleID()
};

// 创建Context
const CommentsContext = createContext<{
  state: CommentsStateType;
  dispatch: React.Dispatch<Action>;
}>({
  state: initialState,
  dispatch: () => null
});

// 创建Provider组件
type CommentsProviderProps = { children: ReactNode };

export const CommentsProvider: React.FC<CommentsProviderProps> = ({ children }) => {
  // 使用useReducer管理状态
  const [state, dispatch] = useReducer((state: CommentsStateType, action: Action) => {
    switch (action.type) {
      case 'SET_COMMENTS':
        return { ...state, comments: action.payload };
      default:
        return state;
    }
  }, initialState);

  useEffect(() => {
    async function fetchComments() {
      try {
        const articleId = getArticleID();
        const resp = await getCommentList(articleId)
        if (resp.data.code===0){
          dispatch({ type: 'SET_COMMENTS', payload: resp.data.data.comments });
        }
      } catch (error) {
        console.error('Failed to fetch comments:', error);
      }
    }

    fetchComments();
  }, []); // 依赖数组为空，表示仅在组件挂载时执行

  return (
    <CommentsContext.Provider value={{ state, dispatch }}>
      {children}
    </CommentsContext.Provider>
  );
};

// 导出自定义hook
export const useComments = () => useContext(CommentsContext);

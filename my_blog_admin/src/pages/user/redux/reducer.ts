import { PaginationProps } from '@arco-design/web-react/es/Pagination/pagination';
import {
  ADD_LIST_ITEM,
  CLEAR_EDIT_MODAL_CONTENT,
  TOGGLE_PASSWORD_VISIBLE,
  TOGGLE_VISIBLE,
  UPDATE_EDIT_MODAL_CONTENT,
  UPDATE_LIST,
  UPDATE_LIST_ITEM,
  UPDATE_LOADING,
} from './actionTypes';

const initialState = {
  data: [],
  loading: true,
  visible: false,
  passwordVisible: false,
  confirmLoading: false,
  pagination: {
    sizeCanChange: true,
    showTotal: true,
    pageSize: 10,
    current: 1,
    total: 0,
    pageSizeChangeResetCurrent: true,
  },
  editModalContent: {
    id: 0,
    username: '',
    nickname: '',
    email: '',
    role: '',
    avatar: '',
  },
};

interface EditModalContent {
  id?: number;
  username?: string;
  nickname?: string;
  email?: string;
  role?: string;
  avatar?: string;
}

export interface UserState {
  data?: any[];
  pagination?: PaginationProps;
  loading?: boolean;
  visible?: boolean;
  confirmLoading?: boolean;
  editModalContent?: EditModalContent;
  passwordVisible?: boolean;
}

export default function(state = initialState, action) {
  switch (action.type) {
    case UPDATE_LOADING: {
      const { loading } = action.payload;
      return {
        ...state,
        loading,
      };
    }
    case UPDATE_LIST: {
      const { data, pagination } = action.payload;
      const p = state.pagination;
      return {
        ...state,
        data,
        pagination:
          pagination != null
            ? {
                sizeCanChange: p.sizeCanChange,
                showTotal: p.showTotal,
                pageSize: pagination.pageSize,
                current: pagination.current,
                total: pagination.total,
                pageSizeChangeResetCurrent: p.pageSizeChangeResetCurrent,
              }
            : state.pagination,
      };
    }
    case TOGGLE_VISIBLE: {
      const { visible } = action.payload;
      return {
        ...state,
        visible,
      };
    }
    case TOGGLE_PASSWORD_VISIBLE: {
      const { passwordVisible } = action.payload;
      return {
        ...state,
        passwordVisible,
      };
    }
    case ADD_LIST_ITEM: {
      const { item } = action.payload;

      return {
        ...state,
        data: [
          ...state.data,
          {
            id: item.id,
            username: item.username,
            nickname: item.nickname,
            email: item.email,
            role: item.role,
            avatar: item.avatar,
            post_count: item.post_count,
            page_count: item.page_count,
          },
        ],
      };
    }
    case UPDATE_LIST_ITEM: {
      const { userInfo } = action.payload;
      const list = state.data;
      for (let i = 0; i < list.length; i++) {
        if (list[i].id === userInfo.id) {
          list[i].nickname = userInfo.nickname;
          list[i].email = userInfo.email;
          list[i].role = userInfo.role;
        }
      }
      return {
        ...state,
        data: list,
      };
    }
    case CLEAR_EDIT_MODAL_CONTENT: {
      return {
        ...state,
        editModalContent: {
          id: 0,
          username: '',
          nickname: '',
          email: '',
          role: '',
          avatar: '',
        },
      };
    }
    case UPDATE_EDIT_MODAL_CONTENT: {
      const { content } = action.payload;
      return {
        ...state,
        editModalContent: {
          id: content.id,
          username: content.username,
          nickname: content.nickname,
          email: content.email,
          role: content.role,
          avatar: content.avatar,
        },
      };
    }

    default:
      return state;
  }
}

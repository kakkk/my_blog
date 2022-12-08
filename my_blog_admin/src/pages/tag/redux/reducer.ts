import { PaginationProps } from '@arco-design/web-react/es/Pagination/pagination';
import {
  ADD_LIST_ITEM,
  CLEAR_EDIT_MODAL_CONTENT,
  TOGGLE_VISIBLE,
  UPDATE_EDIT_MODAL_CONTENT,
  UPDATE_LIST,
  UPDATE_LOADING,
  UPDATE_SEARCH_KEYWORD,
} from './actionTypes';

const initialState = {
  data: [],
  pagination: {
    sizeCanChange: true,
    showTotal: true,
    pageSize: 10,
    current: 1,
    total: 0,
    pageSizeChangeResetCurrent: true,
  },
  loading: true,
  searchKeyWord: '',
  editModalContent: {
    id: 0,
    name: '',
  },
  visible: false,
  confirmLoading: false,
};

interface EditModalContent {
  id?: number;
  name?: string;
}

export interface TagState {
  data?: any[];
  pagination?: PaginationProps;
  loading?: boolean;
  searchKeyWord?: string;
  editModalContent?: EditModalContent;
  visible?: boolean;
  confirmLoading?: boolean;
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
    case UPDATE_SEARCH_KEYWORD: {
      const { searchKeyWord } = action.payload;
      return {
        ...state,
        searchKeyWord,
      };
    }
    case TOGGLE_VISIBLE: {
      const { visible } = action.payload;
      return {
        ...state,
        visible,
      };
    }
    case CLEAR_EDIT_MODAL_CONTENT: {
      return {
        ...state,
        editModalContent: {
          id: 0,
          name: '',
        },
      };
    }
    case UPDATE_EDIT_MODAL_CONTENT: {
      const { content } = action.payload;
      return {
        ...state,
        editModalContent: {
          id: content.id,
          name: content.name,
        },
      };
    }
    case ADD_LIST_ITEM: {
      const { item } = action.payload;
      const p = state.pagination;
      return {
        ...state,
        data: [
          {
            id: item.id,
            name: item.name,
            count: 0,
          },
          ...state.data,
        ],
        pagination: {
          sizeCanChange: p.sizeCanChange,
          showTotal: p.showTotal,
          pageSize: p.pageSize,
          current: p.current,
          total: p.total + 1,
          pageSizeChangeResetCurrent: p.pageSizeChangeResetCurrent,
        },
      };
    }
    default:
      return state;
  }
}

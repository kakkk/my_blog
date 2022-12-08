import {
  ADD_LIST_ITEM,
  CLEAR_EDIT_MODAL_CONTENT,
  TOGGLE_VISIBLE,
  UPDATE_EDIT_MODAL_CONTENT,
  UPDATE_LIST,
  UPDATE_LOADING,
} from './actionTypes';

const initialState = {
  data: [],
  loading: true,
  formParams: {},
  visible: false,
  confirmLoading: false,
  editModalContent: {
    id: 0,
    name: '',
    slug: '',
    description: '',
    parentId: 0,
  },
};

interface FormParams {
  [key: string]: string;
}

interface EditModalContent {
  id?: number;
  name?: string;
  slug?: string;
  description?: string;
  parentId?: number;
}

export interface CategoryState {
  data?: any[];
  formParams?: FormParams;
  loading?: boolean;
  visible?: boolean;
  confirmLoading?: boolean;
  editModalContent?: EditModalContent;
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
      const { data } = action.payload;
      return {
        ...state,
        data,
      };
    }
    case TOGGLE_VISIBLE: {
      const { visible } = action.payload;
      return {
        ...state,
        visible,
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
            name: item.name,
            slug: item.slug,
            count: 0,
            children: 0,
          },
        ],
      };
    }
    case CLEAR_EDIT_MODAL_CONTENT: {
      return {
        ...state,
        editModalContent: {
          id: 0,
          name: '',
          slug: '',
          description: '',
          parentId: 0,
        },
      };
    }
    case UPDATE_EDIT_MODAL_CONTENT: {
      const { content } = action.payload;
      console.log(content);
      return {
        ...state,
        editModalContent: {
          id: content.id,
          name: content.name,
          slug: content.slug,
          description: content.description,
          parentId: content.parent.id,
        },
      };
    }

    default:
      return state;
  }
}

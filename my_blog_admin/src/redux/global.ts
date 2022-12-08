import defaultSettings from '../settings.json';
import { UPDATE_USER_INFO } from './actionTypes';

const defaultTheme = localStorage.getItem('arco-theme') || 'light';

function changeTheme(newTheme?: 'string') {
  if ((newTheme || defaultTheme) === 'dark') {
    document.body.setAttribute('arco-theme', 'dark');
  } else {
    document.body.removeAttribute('arco-theme');
  }
}

// init page theme
changeTheme();

export interface GlobalState {
  theme?: string;
  settings?: typeof defaultSettings;
  userInfo?: {
    id?: number;
    username?: string;
    nickname?: string;
    email?: string;
    role?: string;
    avatar?: string;
    description?: string;
  };
}

const initialState: GlobalState = {
  theme: defaultTheme,
  settings: defaultSettings,
  userInfo: {
    id: 0,
    username: '',
    nickname: '',
    email: '',
    role: '',
    avatar: '',
    description: '',
  },
};

export default function(state = initialState, action) {
  switch (action.type) {
    case 'toggle-theme': {
      const { theme } = action.payload;
      if (theme === 'light' || theme === 'dark') {
        localStorage.setItem('arco-theme', theme);
        changeTheme(theme);
      }

      return {
        ...state,
        theme,
      };
    }
    case 'update-settings': {
      const { settings } = action.payload;
      return {
        ...state,
        settings,
      };
    }
    case UPDATE_USER_INFO:
      const { userInfo } = action.payload;
      return {
        ...state,
        userInfo,
      };
    default:
      return state;
  }
}

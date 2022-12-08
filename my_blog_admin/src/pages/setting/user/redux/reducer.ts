import { UPDATE_PASSWORD_RETYPE_OK } from './actionType';

const initialState = {
  passwordRetypeOk: true,
};

export interface UserSettingState {
  passwordRetypeOk?: boolean;
}

export default function(state = initialState, action) {
  switch (action.type) {
    case UPDATE_PASSWORD_RETYPE_OK: {
      const { passwordRetypeOk } = action.payload;
      return {
        ...state,
        passwordRetypeOk,
      };
    }
    default:
      return state;
  }
}

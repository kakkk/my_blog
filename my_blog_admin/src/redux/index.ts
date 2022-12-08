import { combineReducers } from 'redux';
import global, { GlobalState } from './global';
import tag, { TagState } from '../pages/tag/redux/reducer';
import category, { CategoryState } from '../pages/category/redux/reducer';
import user, { UserState } from '../pages/user/redux/reducer';
import userSetting, { UserSettingState } from '../pages/setting/user/redux/reducer';

export interface ReducerState {
  global: GlobalState;
  tag: TagState;
  category: CategoryState;
  user: UserState;
  userSetting: UserSettingState;
}

export default combineReducers({
  global,
  tag,
  category,
  user,
  userSetting,
});

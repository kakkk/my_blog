import React, { useState, useEffect } from 'react';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import { ConfigProvider } from '@arco-design/web-react';
import zhCN from '@arco-design/web-react/es/locale/zh-CN';
import enUS from '@arco-design/web-react/es/locale/en-US';
import ReactDOM from 'react-dom';
import { Router, Switch, Route } from 'react-router-dom';
import rootReducer from './redux';
import history from './history';
import PageLayout from './layout/page-layout';
import { GlobalContext } from './context';
import './style/index.less';
import Login from './pages/login';
import checkLogin from './utils/checkLogin';
import { getUserInfo } from './api/user';

const store = createStore(rootReducer);

function Index() {
  const localeName = localStorage.getItem('arco-lang') || 'zh-CN';

  if (!localStorage.getItem('arco-lang')) {
    localStorage.setItem('arco-lang', localeName);
  }

  const [locale, setLocale] = useState();

  async function fetchLocale(ln?: string) {
    const locale = (await import(`./locale/${ln || localeName}`)).default;
    setLocale(locale);
  }

  function getArcoLocale() {
    switch (localeName) {
      case 'zh-CN':
        return zhCN;
      case 'en-US':
        return enUS;
      default:
        return zhCN;
    }
  }

  async function fetchUserInfo() {
    try {
      const res: any = await getUserInfo();
      console.log(res);
      if (res.data) {
        if (res.code === 0) {
          store.dispatch({
            type: 'update-userInfo',
            payload: { userInfo: res.data },
          });
        }
      }
    } catch (error) {}
  }

  useEffect(() => {
    fetchLocale();
  }, []);

  useEffect(() => {
    if (checkLogin()) {
      fetchUserInfo();
    } else {
      history.push('/user/login');
    }
  }, []);

  const contextValue = {
    locale,
  };

  return locale ? (
    <Router history={history}>
      <ConfigProvider locale={getArcoLocale()}>
        <Provider store={store}>
          <GlobalContext.Provider value={contextValue}>
            <Switch>
              <Route path="/user/login" component={Login} />
              <Route path="/" component={PageLayout} />
            </Switch>
          </GlobalContext.Provider>
        </Provider>
      </ConfigProvider>
    </Router>
  ) : null;
}

ReactDOM.render(<Index />, document.getElementById('root'));

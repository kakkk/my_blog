import axios from 'axios';
import { Notification } from '@arco-design/web-react';
import history from '../history';

function getBaseURL(): string {
  if (process.env.NODE_ENV === 'development') {
    return 'http://localhost:8888/api';
  }
  const protocol = window.location.protocol;
  const host = window.location.hostname;
  const port = window.location.port;
  return `${protocol}//${host}${port ? `:${port}` : ''}/api`;
}

export const request = (config) => {
  const http = axios.create({
    baseURL: getBaseURL(),
    headers: {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Credentials': 'true',
      'Content-Type': 'application/json',
    },
    withCredentials: true,
    // timeout: 5000,
  });

  // 请求拦截
  http.interceptors.request.use(
    (config) => {
      console.log('config', config);
      return config;
    },
    () => {}
  );

  // 响应拦截
  http.interceptors.response.use(
    (res) => {
      console.log('res-------', res);
      if (!(res.data.code === 0)) {
        switch (res.data.code) {
          case 400100:
            Notification.error({ title: '登录已过期', content: res.data.msg });
            history.push('/user/login');
            break;
          default:
            Notification.error({ title: '操作失败', content: res.data.msg });
            break;
        }
      }
      return res.data ? res.data : res;
    },
    (error) => {
      console.log('error===', error.response); // 注意这里必须打印error.response
      const response = error.response;
      if (response && response.status) {
        if (response.status === 403) {
          // location.href = '/403';
          history.push('/user/login');
          Notification.error({ title: '权限错误', content: response.data.msg });
        }
        if (response.status === 401) {
          localStorage.removeItem('login');
          history.push('/user/login');
          Notification.error({ title: '登陆过期', content: '登陆过期，请重新登录' });
        }
      }
    }
  );

  return http(config);
};

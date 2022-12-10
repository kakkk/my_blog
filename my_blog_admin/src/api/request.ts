import axios from 'axios';
import { Notification } from '@arco-design/web-react';

export const request = (config) => {
  const http = axios.create({
    baseURL: 'http://localhost:8888/api',
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
        Notification.error({ title: '权限错误', content: res.data.msg });
        switch (res.data.code) {
          case 400100:
            location.href = '/login';
            break;
          default:
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
          location.href = '/login';
          Notification.error({ title: '权限错误', content: response.data.msg });
        }
        if (response.status === 401) {
          localStorage.removeItem('login');
          location.href = '/login';
          Notification.error({ title: '登陆过期', content: '登陆过期，请重新登录' });
        }
      }
    }
  );

  return http(config);
};

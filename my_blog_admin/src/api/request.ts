import axios from 'axios';
import { Notification } from '@arco-design/web-react';

export const request = (config) => {
  const http = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
    // timeout: 5000,
  });

  // 请求拦截
  http.interceptors.request.use(
    (config) => {
      console.log('config', config);
      const token = localStorage.getItem('token');
      config.headers = {
        Authorization: `${token}`,
      };
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
          case -4:
            location.href = '/403';
            break;
          case -3:
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
          localStorage.removeItem('token');
          location.href = '/login';
          Notification.error({ title: 'Token错误', content: 'token过期，请重新登录' });
        }
      }
    }
  );

  return http(config);
};

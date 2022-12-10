import { request } from './request';

export async function login(username: string, password: string) {
  return request({
    url: '/admin/login',
    method: 'POST',
    data: {
      username,
      password,
    },
  });
}

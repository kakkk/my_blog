import { request } from './request';

export async function login(data) {
  return request({
    url: '/admin/login',
    method: 'POST',
    data,
  });
}

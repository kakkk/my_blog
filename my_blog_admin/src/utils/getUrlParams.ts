// 仅用于线上预览，实际使用中可以将此逻辑删除
import qs from 'query-string';

export type ParamsType = Record<string, any>;

export default function getUrlParams(): ParamsType {
  const params = qs.parseUrl(window.location.href).query;
  const returnParams: ParamsType = {};
  Object.keys(params).forEach((p) => {
    returnParams[p] = params[p];
  });
  return returnParams;
}

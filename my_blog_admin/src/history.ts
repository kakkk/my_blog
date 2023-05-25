/**
 * history.push(path, [state])
 * history.replace(path, [state])
 * history.go(n)
 * history.goBack()
 * history.goForward()
 *
 * history.listen(func) // listen for changes to the current location
 *
 */
import { createBrowserHistory } from 'history';

let basename: string;
if (process.env.NODE_ENV === 'development') {
  basename = '/';
}
if (process.env.NODE_ENV === 'production') {
  basename = '/admin/';
}

const HISTORY = createBrowserHistory({
  basename,
});

export default HISTORY;

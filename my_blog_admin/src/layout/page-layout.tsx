import React, { useState, useRef, useMemo } from 'react';
import { Switch, Route, Link, Redirect } from 'react-router-dom';
import { Layout, Menu } from '@arco-design/web-react';
import { IconMenuFold, IconMenuUnfold } from '@arco-design/web-react/icon';
import qs from 'query-string';
import Navbar from '../components/NavBar';
import Footer from '../components/Footer';
import LoadingBar from '../components/LoadingBar';
import { routes, defaultRoute } from '../routes';
import { isArray } from '../utils/is';
import history from '../history';
import useLocale from '../utils/useLocale';
import lazyload from '../utils/lazyload';
import styles from './style/layout.module.less';

const MenuItem = Menu.Item;
const SubMenu = Menu.SubMenu;

const Sider = Layout.Sider;
const Content = Layout.Content;

function getFlattenRoutes() {
  const res = [];

  function travel(_routes) {
    _routes.forEach((route) => {
      if (route.componentPath) {
        route.component = lazyload(() => import(`../pages/${route.componentPath}`));
        res.push(route);
      } else if (isArray(route.children) && route.children.length) {
        travel(route.children);
      }
    });
  }

  travel(routes);
  return res;
}

function renderRoutes(locale) {
  const nodes = [];
  function travel(_routes, level) {
    return _routes.map((route) => {
      const titleDom = (
        <>
          {route.icon} {locale[route.name] || route.name}
        </>
      );
      if (
        route.component &&
        (!isArray(route.children) || (isArray(route.children) && !route.children.length))
      ) {
        if (level > 1) {
          return <MenuItem key={route.key}>{titleDom}</MenuItem>;
        }
        nodes.push(
          <MenuItem key={route.key}>
            <Link to={`/${route.key}`}>{titleDom}</Link>
          </MenuItem>
        );
      }
      if (isArray(route.children) && route.children.length) {
        if (level > 1) {
          return (
            <SubMenu key={route.key} title={titleDom}>
              {travel(route.children, level + 1)}
            </SubMenu>
          );
        }
        nodes.push(
          <SubMenu key={route.key} title={titleDom}>
            {travel(route.children, level + 1)}
          </SubMenu>
        );
      }
    });
  }

  travel(routes, 1);
  return nodes;
}

function PageLayout() {
  const pathname = history.location.pathname;
  const currentComponent = qs.parseUrl(pathname).url.slice(1);
  const defaultSelectedKeys = [currentComponent || defaultRoute];

  const locale = useLocale();
  const [collapsed, setCollapsed] = useState<boolean>(false);
  const [selectedKeys, setSelectedKeys] = useState<string[]>(defaultSelectedKeys);
  const loadingBarRef = useRef(null);

  const navbarHeight = 60;
  const menuWidth = collapsed ? 48 : 220;

  const flattenRoutes = useMemo(() => getFlattenRoutes() || [], []);

  function onClickMenuItem(key) {
    const currentRoute = flattenRoutes.find((r) => r.key === key);
    const component = currentRoute.component;
    const preload = component.preload();
    loadingBarRef.current.loading();
    preload.then(() => {
      setSelectedKeys([key]);
      history.push(currentRoute.path ? currentRoute.path : `/${key}`);
      loadingBarRef.current.success();
    });
  }

  function toggleCollapse() {
    setCollapsed((collapsed) => !collapsed);
  }

  const paddingLeft = { paddingLeft: menuWidth };
  const paddingTop = { paddingTop: navbarHeight };
  const paddingStyle = { ...paddingLeft, ...paddingTop };

  return (
    <Layout className={styles.layout}>
      <LoadingBar ref={loadingBarRef} />
      <div className={styles.layoutNavbar}>
        <Navbar />
      </div>
      <Layout>
        <Sider
          className={styles.layoutSider}
          width={menuWidth}
          collapsed={collapsed}
          onCollapse={setCollapsed}
          trigger={null}
          collapsible
          breakpoint="xl"
          style={paddingTop}
        >
          <div className={styles.menuWrapper}>
            <Menu
              collapse={collapsed}
              onClickMenuItem={onClickMenuItem}
              selectedKeys={selectedKeys}
            >
              {renderRoutes(locale)}
            </Menu>
          </div>
          <div className={styles.collapseBtn} onClick={toggleCollapse}>
            {collapsed ? <IconMenuUnfold /> : <IconMenuFold />}
          </div>
        </Sider>
        <Layout className={styles.layoutContent} style={paddingStyle}>
          <Content>
            <Switch>
              {flattenRoutes.map((route, index) => {
                return <Route key={index} path={`/${route.key}`} component={route.component} />;
              })}
              <Redirect push to={`/${defaultRoute}`} />
            </Switch>
          </Content>
          <Footer />
        </Layout>
      </Layout>
    </Layout>
  );
}

export default PageLayout;

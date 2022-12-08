import { Typography } from '@arco-design/web-react';
import React from 'react';
import { useSelector } from 'react-redux';
import { ReducerState } from '../../redux';
import styles from './style/index.module.less';

export default function Dashboard() {
  const userInfo = useSelector((state: ReducerState) => state.global.userInfo) || {};
  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Typography.Title heading={5} style={{ marginTop: 0 }}>
          欢迎
        </Typography.Title>
        <Typography.Text type="secondary">
          {userInfo.nickname}, 这里是dBlog后台管理系统
        </Typography.Text>
      </div>
      <div className={styles.content}>
        <div />
      </div>
    </div>
  );
}

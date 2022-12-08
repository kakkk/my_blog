import React from 'react';
import { Breadcrumb, Grid } from '@arco-design/web-react';
import styles from './style/index.module.less';
import history from '../../../history';
import UserInfo from './userInfo';
import BasicInfo from './basicInfo';
import Password from './password';

export default function UserSetting() {
  const Row = Grid.Row;
  const Col = Grid.Col;
  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>设置</Breadcrumb.Item>
        <Breadcrumb.Item>
          <a
            onClick={() => {
              history.push(`/setting/user`);
              // fetchData();
            }}
            href="#"
          >
            用户设置
          </a>
        </Breadcrumb.Item>
      </Breadcrumb>

      <Row>
        <Col flex="350px">
          <UserInfo />
        </Col>
        <Col flex="auto">
          <Row>
            <Col span={12}>
              <BasicInfo />
            </Col>
            <Col span={12}>
              <Password />
            </Col>
          </Row>
        </Col>
      </Row>
    </div>
  );
}

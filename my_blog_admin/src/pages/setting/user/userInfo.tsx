import React from 'react';
import { Avatar, Card, Grid, Tag } from '@arco-design/web-react';
import { IconEdit, IconFile, IconUser } from '@arco-design/web-react/icon';
import styles from './style/index.module.less';

export default function UserInfo() {
  const Row = Grid.Row;
  const Col = Grid.Col;
  return (
    <Card title="基本信息" style={{ margin: '5px', height: '500px' }}>
      <div className={styles['avatar-avatar']}>
        <Avatar triggerIcon={<IconEdit />} style={{ backgroundColor: '#14C9C9' }} size={150}>
          <IconUser />
        </Avatar>
      </div>
      <div className={styles['avatar-info']}>
        <Tag style={{ fontSize: 'larger' }} size="large" color="blue" icon={<IconUser />}>
          kakkk
        </Tag>
      </div>
      <div className={styles['avatar-summary']}>
        <Row style={{ padding: '10px' }}>
          <Col span={6}>
            <div>
              <div style={{ fontSize: 'medium' }}>
                <IconEdit />
                文章:
              </div>
            </div>
          </Col>
          <Col span={6}>
            <Tag color="green">15</Tag>
          </Col>
          <Col span={6}>
            <div>
              <div style={{ fontSize: 'medium' }}>
                <IconFile />
                页面:
              </div>
            </div>
          </Col>
          <Col span={6}>
            <Tag color="green">15</Tag>
          </Col>
        </Row>
        <Row style={{ padding: '10px' }}>
          <Col span={9}>
            <div style={{ fontSize: 'medium' }}>
              <IconUser />
              上次登录：
            </div>
          </Col>
          <Col span={15}>
            <Tag>2022-03-13 17:00:00</Tag>
          </Col>
        </Row>
      </div>
    </Card>
  );
}

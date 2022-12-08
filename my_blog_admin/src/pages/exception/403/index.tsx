import React from 'react';
import { Result, Button } from '@arco-design/web-react';
import styles from './style/index.module.less';

function Exception403() {
  return (
    <div className={styles.container}>
      <div className={styles.wrapper}>
        <Result
          className={styles.result}
          status="403"
          subTitle="对不起，您没有访问该资源的权限"
          extra={
            <Button key="back" type="primary">
              返回
            </Button>
          }
        />
      </div>
    </div>
  );
}

export default Exception403;

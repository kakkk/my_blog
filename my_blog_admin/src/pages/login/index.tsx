import React, { useEffect } from 'react';
import ParticlesBg from 'particles-bg';
import Footer from '../../components/Footer';
import LoginForm from './form';

import styles from './style/index.module.less';

export default () => {
  useEffect(() => {
    document.body.setAttribute('arco-theme', 'light');
  }, []);
  return (
    <>
      <div className={styles.container}>
        <div className={styles.content}>
          <div className={styles['content-inner']}>
            <LoginForm />
          </div>
          <div className={styles.footer}>
            <Footer />
          </div>
        </div>
      </div>
      <ParticlesBg type="circle" bg />
    </>
  );
};

import { Form, Input, Checkbox, Button, Space } from '@arco-design/web-react';
import { FormInstance } from '@arco-design/web-react/es/Form';
import { IconLock, IconUser } from '@arco-design/web-react/icon';
import React, { useEffect, useRef, useState } from 'react';
import { login as adminLogin } from '../../api/login';
import styles from './style/index.module.less';
import history from '../../history';

export default function LoginForm() {
  const formRef = useRef<FormInstance>();
  const [errorMessage, setErrorMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [rememberPassword, setRememberPassword] = useState(false);

  function afterLoginSuccess(params, resp) {
    // 记住密码
    if (rememberPassword) {
      localStorage.setItem('loginParams', JSON.stringify(params));
    } else {
      localStorage.removeItem('loginParams');
    }
    localStorage.setItem('token', resp.data.token);
    localStorage.setItem('role', resp.data.role);
    // 跳转首页
    window.location.href = history.createHref({
      pathname: '/',
    });
  }

  async function login(params) {
    setErrorMessage('');
    setLoading(true);
    try {
      const res: any = await adminLogin(params);
      console.log(res);
      if (res.data) {
        if (res.code === 0) {
          afterLoginSuccess(params, res);
        }
      } else {
        setErrorMessage(res.msg);
      }
    } catch (error) {
    } finally {
      setLoading(false);
    }
  }

  function onSubmitClick() {
    formRef.current.validate().then((values) => {
      login(values);
    });
  }

  // 读取 localStorage，设置初始值
  useEffect(() => {
    const params = localStorage.getItem('loginParams');
    const rememberPassword = !!params;
    setRememberPassword(rememberPassword);
    if (formRef.current && rememberPassword) {
      const parseParams = JSON.parse(params);
      formRef.current.setFieldsValue(parseParams);
    }
  }, []);

  return (
    <div className={styles['login-form-wrapper']}>
      <div className={styles['login-form-title']}>dBlog Login</div>
      <div className={styles['login-form-error-msg']}>{errorMessage}</div>
      <Form className={styles['login-form']} layout="vertical" ref={formRef}>
        <Form.Item field="userName" rules={[{ required: true, message: '用户名不能为空' }]}>
          <Input
            size="large"
            prefix={<IconUser />}
            placeholder="用户名或邮箱"
            onPressEnter={onSubmitClick}
          />
        </Form.Item>
        <Form.Item field="password" rules={[{ required: true, message: '密码不能为空' }]}>
          <Input.Password
            size="large"
            prefix={<IconLock />}
            placeholder="密码"
            onPressEnter={onSubmitClick}
          />
        </Form.Item>
        <Space size={16} direction="vertical">
          <div className={styles['login-form-password-actions']}>
            <Checkbox checked={rememberPassword} onChange={setRememberPassword}>
              记住密码
            </Checkbox>
          </div>
          <Button size="large" type="primary" long onClick={onSubmitClick} loading={loading}>
            登录
          </Button>
        </Space>
      </Form>
    </div>
  );
}

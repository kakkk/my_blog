import React, { useEffect, useState } from 'react';
import {
  Breadcrumb,
  Button,
  Card,
  Form,
  Input,
  Message,
  Popconfirm,
  Popover,
} from '@arco-design/web-react';
import { IconClose } from '@arco-design/web-react/icon';
import Editor, { Themes } from 'md-editor-rt';
import { useSelector } from 'react-redux';
import styles from './style/index.module.less';
import 'md-editor-rt/lib/style.css';
import { addPage, editPage, getPage } from '../../../api/page';
import getUrlParams from '../../../utils/getUrlParams';
import history from '../../../history';
import { ReducerState } from '../../../redux';

export default function PageEdit(props) {
  const [text, setText] = useState('');
  const [visible, setVisible] = React.useState(false);
  const [slug, setSlug] = useState('');
  const [title, setTitle] = useState('');
  const theme = useSelector((state: ReducerState) => state.global.theme);
  const FormItem = Form.Item;
  const fetchData = async (id: number) => {
    try {
      const res: any = await getPage(id);
      if (res.code === 0) {
        if (res.data !== null) {
          setTitle(res.data.title);
          setSlug(res.data.slug);
          setText(res.data.content);
        } else {
          setTitle('');
          setSlug('');
          setText('');
        }
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };
  useEffect(() => {
    if (getUrlParams().id) {
      fetchData(getUrlParams().id);
    }
  }, [props.location]);
  const createPage = async (title: string, content: string, slug: string) => {
    try {
      const res: any = await addPage(title, content, slug);
      if (res.code === 0) {
        Message.success('发布成功');
        history.push(`/page/edit?id=${res.data.id}`);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const onCreate = async () => {
    await createPage(title, text, slug);
  };

  const onUpdate = async () => {
    try {
      const res: any = await editPage(getUrlParams().id, title, text, slug);
      if (res.code === 0) {
        Message.success('发布成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const onPublish = async () => {
    if (title === '') {
      Message.error('标题不能为空!');
      return;
    }
    if (slug === '') {
      Message.error('缩略名不能为空!');
      return;
    }
    if (text === '') {
      Message.error('内容不能为空');
      return;
    }
    if (getUrlParams().id) {
      await onUpdate();
    } else {
      await onCreate();
    }
    setVisible(false);
  };

  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>页面</Breadcrumb.Item>
        {getUrlParams().id && getUrlParams().id !== 0 && (
          <Breadcrumb.Item>
            <a
              onClick={() => {
                history.push(`/page/list`);
              }}
              href="#"
            >
              页面管理
            </a>
          </Breadcrumb.Item>
        )}
        {getUrlParams().id && getUrlParams().id !== 0 ? (
          <Breadcrumb.Item>编辑页面</Breadcrumb.Item>
        ) : (
          <Breadcrumb.Item>添加页面</Breadcrumb.Item>
        )}
      </Breadcrumb>
      <Card>
        <div style={{ display: 'flex', paddingBottom: '15px' }}>
          <div style={{ flex: 'auto' }}>
            <Input size="large" placeholder="请输入标题..." value={title} onChange={setTitle} />
          </div>
          <div style={{ textAlign: 'right', paddingLeft: '15px' }}>
            <Popover
              title={
                <div style={{ display: 'flex', paddingBottom: '5px' }}>
                  <div style={{ marginRight: 'auto' }}>页面设置</div>
                  <Button icon={<IconClose />} type="text" onClick={() => setVisible(false)} />
                </div>
              }
              popupVisible={visible}
              onVisibleChange={(visible) => {
                if (visible) {
                  setVisible(true);
                }
              }}
              style={{ width: '500px' }}
              content={
                <div>
                  <Form labelCol={{ span: 6 }} wrapperCol={{ span: 18 }}>
                    <FormItem label="缩略名">
                      <Input value={slug} onChange={setSlug} />
                    </FormItem>
                  </Form>
                  <div style={{ textAlign: 'right' }}>
                    <Popconfirm title="确定发布？" onOk={onPublish}>
                      <Button type="primary" style={{ marginLeft: '10px' }}>
                        发布
                      </Button>
                    </Popconfirm>
                  </div>
                </div>
              }
              position="br"
            >
              <Button size="large" type="primary">
                发布
              </Button>
            </Popover>
          </div>
        </div>
        <div>
          <Editor
            modelValue={text}
            onChange={setText}
            style={{ height: '450px' }}
            theme={theme as Themes}
          />
        </div>
      </Card>
    </div>
  );
}

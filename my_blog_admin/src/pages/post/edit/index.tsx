import React, { useEffect, useState } from 'react';
import { Dayjs } from 'dayjs';
import {
  Breadcrumb,
  Button,
  Card,
  Checkbox,
  DatePicker,
  Form,
  Input,
  Message,
  Popover,
  Select,
} from '@arco-design/web-react';
import Editor, { Themes } from 'md-editor-rt';
import { useSelector } from 'react-redux';
import { IconClose } from '@arco-design/web-react/icon';
import { dayjs } from '@arco-design/web-react/es/_util/dayjs';
import styles from './style/index.module.less';
import 'md-editor-rt/lib/style.css';

import { ReducerState } from '../../../redux';
import { getCategoryList } from '../../../api/category';
import { getTagList } from '../../../api/tag';
import { addPost, getPost, updatePost, updatePostPublish } from '../../../api/post';
import getUrlParams from '../../../utils/getUrlParams';
import history from '../../../history';

export default function PostEdit(props) {
  const [text, setText] = useState('');
  const [visible, setVisible] = React.useState(false);
  const theme = useSelector((state: ReducerState) => state.global.theme);
  const FormItem = Form.Item;
  const [categoryOptions, setCategoryOptions] = useState([] as { label: string; value: number }[]);
  const [labelOptions, setLabelOptions] = useState([] as string[]);
  const [title, setTitle] = useState('');
  const [categories, setCategories] = useState([] as number[]);
  const [label, setLabel] = useState([] as string[]);
  const [publishAt, setPublishAt] = useState(null as Dayjs);
  const [publish, setPublish] = useState(false);
  const fetchData = async (id: number) => {
    try {
      const res: any = await getPost(id);
      if (res.code === 0) {
        setTitle(res.data.title);
        setText(res.data.content);
        const cts: number[] = new Array(0);
        for (const category of res.data.categories) {
          cts.push(category.id);
        }
        setCategories(cts);
        const ls: string[] = new Array(0);
        for (const tag of res.data.tags) {
          ls.push(tag.name);
        }
        setLabel(ls);
        setPublishAt(dayjs.unix(res.data.publish_at));
        setPublish(res.data.publish);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const getCategory = async () => {
    try {
      const res: any = await getCategoryList(-1);
      if (res.code === 0) {
        const options: { label: string; value: number }[] = new Array(0);
        for (const c of res.data) {
          options.push({ label: c.name, value: c.id });
        }
        setCategoryOptions(options);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  useEffect(() => {
    getCategory();
    if (getUrlParams().id) {
      fetchData(getUrlParams().id);
    }
  }, [props.location]);

  const onUpdate = async () => {
    try {
      const res: any = await updatePost(
        getUrlParams().id,
        title,
        text,
        categories,
        label,
        publishAt.unix()
      );
      if (res.code === 0) {
        Message.success('保存成功');
        setVisible(false);
        await updatePostPublish(getUrlParams().id, publish);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const onPublish = async () => {
    if (getUrlParams().id) {
      setPublish(true);
      await onUpdate();
      return;
    }
    try {
      const res: any = await addPost(title, text, true, categories, label, publishAt.unix());
      if (res.code === 0) {
        Message.success('发布成功');
        setVisible(false);
        history.push(`/post/edit?id=${res.data.id}`);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const onSave = async () => {
    if (getUrlParams().id) {
      await onUpdate();
      return;
    }
    try {
      const res: any = await addPost(title, text, false, categories, label, publishAt.unix());
      if (res.code === 0) {
        Message.success('保存成功');
        setVisible(false);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const searchLabel = async (value: string) => {
    try {
      const res: any = await getTagList(value, 1, 15);
      if (res.code === 0) {
        if (res.data === null) {
          setLabelOptions([]);
        }
        const options: string[] = new Array(0);
        for (const tag of res.data) {
          options.push(tag.name);
        }
        setLabelOptions(options);
      } else {
        Message.error(res.message);
      }
    } finally {
    }
  };

  const CheckboxGroup = Checkbox.Group;
  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>文章</Breadcrumb.Item>
        {getUrlParams().id && getUrlParams().id !== 0 && (
          <Breadcrumb.Item>
            <a
              onClick={() => {
                history.push(`/post/list`);
              }}
              href="#"
            >
              文章管理
            </a>
          </Breadcrumb.Item>
        )}
        {getUrlParams().id && getUrlParams().id !== 0 ? (
          <Breadcrumb.Item>编辑文章</Breadcrumb.Item>
        ) : (
          <Breadcrumb.Item>添加文章</Breadcrumb.Item>
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
                  <div style={{ marginRight: 'auto' }}>文章设置</div>
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
                    <FormItem label="分类">
                      <div style={{ height: 120, overflow: 'auto' }}>
                        <CheckboxGroup
                          direction="vertical"
                          defaultValue={categories}
                          options={categoryOptions}
                          value={categories}
                          onChange={setCategories}
                        />
                      </div>
                    </FormItem>
                    <FormItem label="标签">
                      <Select
                        allowCreate
                        showSearch
                        filterOption={false}
                        mode="multiple"
                        placeholder="Please select"
                        options={labelOptions}
                        allowClear
                        value={label}
                        onChange={setLabel}
                        onSearch={(value) => {
                          searchLabel(value);
                        }}
                        onVisibleChange={(visible) => {
                          if (!visible) {
                            setLabelOptions([]);
                          }
                        }}
                      />
                    </FormItem>
                    <FormItem label="发布日期">
                      <DatePicker
                        style={{ width: '100%' }}
                        showTime
                        value={publishAt}
                        onChange={(_, date) => {
                          setPublishAt(date);
                        }}
                      />
                    </FormItem>
                  </Form>
                  <div style={{ textAlign: 'right' }}>
                    <Button type="primary" onClick={onSave}>
                      保存
                    </Button>
                    <Button type="primary" style={{ marginLeft: '10px' }} onClick={onPublish}>
                      发布
                    </Button>
                  </div>
                </div>
              }
              position="br"
            >
              <Button size="large" type="primary">
                发布/保存
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

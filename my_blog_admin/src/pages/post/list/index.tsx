import React, { useEffect, useState } from 'react';
import {
  Breadcrumb,
  Button,
  Card,
  Form,
  Grid,
  Input,
  Message,
  Popconfirm,
  Select,
  Table,
  Tag,
} from '@arco-design/web-react';
import { ColumnProps } from '@arco-design/web-react/es/Table/interface';
import { IconLink, IconSearch, IconUser } from '@arco-design/web-react/icon';
import { PaginationProps } from '@arco-design/web-react/es/Pagination/pagination';
import history from '../../../history';
import styles from './style/index.module.less';
import { formatDate } from '../../../utils/formatDate';
import getUrlParams from '../../../utils/getUrlParams';
import { deletePost, getPostList, updatePostStatus } from '../../../api/post';
import { getCategoryList } from '../../../api/category';
import { getTagList } from '../../../api/tag';

export default function PostList(props) {
  const options = ['草稿', '已发布', '下架'];
  const { Row, Col } = Grid;
  const [colData, setColData] = useState(
    new Array(0) as {
      id: number;
      title: string;
      categories: string;
      editor: string;
      pv: number;
      status: number;
      updateAt: number;
      publishAt: number;
    }[]
  );
  const [categoryList, setCategoryList] = useState(new Array(0) as string[]);
  const [selectCategory, setSelectCategory] = useState(new Array(0) as string[]);
  const [labelOptions, setLabelOptions] = useState([] as string[]);
  const [label, setLabel] = useState([] as string[]);
  const [searchTitle, setSearchTitle] = useState('');
  const [pagination, setPagination] = useState({
    sizeCanChange: true,
    showTotal: true,
    pageSize: 10,
    current: 1,
    total: 0,
    pageSizeChangeResetCurrent: true,
  } as PaginationProps);
  const [loading, setLoading] = useState(false);
  const makeCategory = (item) => {
    let str: string = item[0];
    for (let i = 1; i < item.length; i++) {
      str = `${str}、${item[i]}`;
    }
    return str;
  };
  const searchLabel = async (value: string) => {
    try {
      const res: any = await getTagList(value, 1, 15);
      if (res.code === 0) {
        if (res.data === null) {
          setLabelOptions([]);
        }
        const options: string[] = new Array(0);
        for (const tag of res.data.tag_list) {
          options.push(tag.name);
        }
        setLabelOptions(options);
      } else {
        Message.error(res.message);
      }
    } finally {
    }
  };

  const getCategory = async () => {
    try {
      const res: any = await getCategoryList();
      if (res.code === 0) {
        const options: string[] = new Array(0);
        for (const c of res.data.category_list) {
          options.push(c.name);
        }
        setCategoryList(options);
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const fetchData = async () => {
    setLoading(true);
    const urlParams = getUrlParams();
    const page = urlParams.page ? Number(urlParams.page) : 0;
    const size = urlParams.size ? Number(urlParams.size) : 10;
    const titleFromParams = urlParams.title ? urlParams.title : '';
    const categoriesFromParams = urlParams.categories ? urlParams.categories : '[]';
    const tagsFromParams = urlParams.tags ? urlParams.tags : '[]';
    setSearchTitle(titleFromParams);
    const cList: string[] = JSON.parse(categoriesFromParams);
    setSelectCategory(cList);
    const tList: string[] = JSON.parse(tagsFromParams);
    setLabel(tList);
    await getCategory();
    await searchLabel('');
    try {
      const res: any = await getPostList(searchTitle, cList, tList, page, size);
      if (res.code === 0) {
        const newColData: {
          id: number;
          title: string;
          categories: string;
          editor: string;
          pv: number;
          status: number;
          updateAt: number;
          publishAt: number;
        }[] = new Array(0);
        if (res.data.post_list) {
          for (const datum of res.data.post_list) {
            const categoryStr = makeCategory(datum.category_list);
            const d = {
              id: datum.id,
              categories: categoryStr,
              pv: datum.pv,
              editor: datum.editor,
              status: datum.status,
              publishAt: datum.publish_at,
              title: datum.title,
              updateAt: datum.update_at,
            };
            newColData.push(d);
          }
        }

        setColData(newColData);
        setPagination({
          sizeCanChange: pagination.sizeCanChange,
          showTotal: pagination.showTotal,
          pageSize: pagination.pageSize,
          current: res.data.pagination.page,
          total: res.data.pagination.total,
          pageSizeChangeResetCurrent: pagination.pageSizeChangeResetCurrent,
        });
      } else {
        Message.error(res.msg);
      }
    } finally {
      setLoading(false);
    }
  };
  const onDelete = async (id: number) => {
    try {
      const res: any = await deletePost(id);
      if (res.code === 0) {
        fetchData();
        Message.success('删除成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const onSearch = async () => {
    let categoryStr = '[]';
    let tagStr = '[]';
    if (selectCategory.length > 0) {
      categoryStr = JSON.stringify(selectCategory);
    }
    if (label.length > 0) {
      tagStr = JSON.stringify(label);
    }
    history.push(
      `/post/list?page=0&size=10&title=${searchTitle}&categories=${categoryStr}&tags=${tagStr}`
    );
  };

  const onChangeStatus = async (id: number, value: string) => {
    let status: number;
    if (value === '下架') {
      status = 3;
    } else if (value === '已发布') {
      status = 2;
    } else {
      status = 1;
    }
    try {
      const res: any = await updatePostStatus(id, status);
      if (res.code === 0) {
        fetchData();
        Message.success('操作成功');
      } else {
        Message.error(res.msg);
      }
    } finally {
    }
  };

  const columns: ColumnProps[] = [
    {
      title: '标题',
      dataIndex: 'title',
      width: 280,
      render: (_, record) => {
        return (
          <div>
            {record.title}
            <Button
              type="text"
              icon={<IconLink />}
              onClick={() => {
                window.open(`http://localhost:8888/archives/${record.id}`);
              }}
            />
          </div>
        );
      },
    },
    {
      title: '作者',
      dataIndex: 'editor',
      align: 'center',
      render: (item) => {
        return <Tag icon={<IconUser />}>{item}</Tag>;
      },
    },
    {
      title: '分类',
      dataIndex: 'categories',
      align: 'center',
    },
    {
      title: '阅读',
      dataIndex: 'pv',
      align: 'center',
    },
    {
      title: '状态',
      dataIndex: 'status',
      align: 'center',
      render: (_, record) => {
        return (
          <Select
            options={options}
            defaultValue={record.status === 2 ? '已发布' : record.status === 1 ? '草稿' : '下架'}
            onChange={(value) => {
              onChangeStatus(record.id, value);
            }}
            style={{ width: 100 }}
          />
        );
      },
    },
    {
      title: '更新时间',
      dataIndex: 'updateAt',
      align: 'center',
      render: (item) => {
        return <div>{formatDate(item * 1000, 'yyyy-MM-dd HH:ss:mm')}</div>;
      },
    },
    {
      title: '操作',
      width: 200,
      align: 'center',
      dataIndex: 'operations',
      render: (_, record) => (
        <div className={styles.operations}>
          <Button
            type="text"
            size="small"
            onClick={() => {
              history.push(`/post/edit?id=${record.id}`);
            }}
          >
            修改
          </Button>
          <Popconfirm
            title="确认删除？"
            onOk={() => {
              onDelete(record.id);
            }}
          >
            <Button type="text" size="small">
              删除
            </Button>
          </Popconfirm>
        </div>
      ),
    },
  ];
  useEffect(() => {
    fetchData();
  }, [props.location]);

  const onChangeTable = (pagination) => {
    setPagination(pagination);
    history.push(
      `/post/list?page=${pagination.current}&size=${pagination.pageSize}${
        searchTitle === '' ? `&${searchTitle}` : ''
      }`
    );
  };
  return (
    <div className={styles.container}>
      <Breadcrumb style={{ marginBottom: 20 }}>
        <Breadcrumb.Item>文章</Breadcrumb.Item>
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
      </Breadcrumb>
      <Card bordered={false}>
        <div className={styles.toolbar}>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'space-between',
              paddingRight: '20px',
              marginBottom: '20px',
              borderRight: '1px solid var(--color-border-2)',
              boxSizing: 'border-box',
            }}
          >
            <Button
              type="primary"
              onClick={() => {
                history.push('/post/edit');
              }}
            >
              添加文章
            </Button>
          </div>
          <Form
            labelAlign="left"
            labelCol={{ span: 4 }}
            wrapperCol={{ span: 20 }}
            style={{ paddingLeft: '20px' }}
          >
            <Row>
              <Col span={8} style={{ paddingRight: '20px' }}>
                <Form.Item label="分类">
                  <Select
                    allowCreate={false}
                    filterOption={false}
                    mode="multiple"
                    placeholder="Please select"
                    options={categoryList}
                    allowClear
                    value={selectCategory}
                    onChange={setSelectCategory}
                  />
                </Form.Item>
              </Col>
              <Col span={8} style={{ paddingRight: '20px' }}>
                <Form.Item label="标签">
                  <Select
                    allowCreate={false}
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
                  />
                </Form.Item>
              </Col>
              <Col span={8} style={{ paddingRight: '20px' }}>
                <Form.Item label="标题">
                  <Input placeholder="Title" onChange={setSearchTitle} value={searchTitle} />
                </Form.Item>
              </Col>
            </Row>
          </Form>
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'space-between',
              paddingLeft: '20px',
              marginBottom: '20px',
              borderLeft: '1px solid var(--color-border-2)',
              boxSizing: 'border-box',
            }}
          >
            <Button type="primary" icon={<IconSearch />} onClick={onSearch}>
              查询
            </Button>
          </div>
        </div>
        <Table
          rowKey="id"
          columns={columns}
          border={false}
          data={colData}
          pagination={pagination}
          loading={loading}
          onChange={onChangeTable}
        />
      </Card>
    </div>
  );
}

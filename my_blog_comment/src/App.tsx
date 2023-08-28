import React from 'react';
import { Comment, Input, Button } from '@arco-design/web-react';
import {IconMessage } from '@arco-design/web-react/icon';
import { Grid } from '@arco-design/web-react';
import './App.css'

const Row = Grid.Row;
const Col = Grid.Col;

const TextArea = Input.TextArea;

export default function App () {
  const actions = (
    <span className='custom-comment-action'>
      <IconMessage/> Reply
    </span>
  );
  return (
    <div>
      <Comment
        align='right'
        avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/e278888093bef8910e829486fb45dd69.png~tplv-uwbnlip3yd-webp.webp'
        content={
          <div>
            <Row>
              <Col sm={8} xs={24}>
                <Input  addBefore='昵称' className="info-input"/>
              </Col>
              <Col sm={8} xs={24}>
                <Input addBefore='邮箱' className="info-input"/>
              </Col>
              <Col sm={8} xs={24}>
                <Input addBefore='网址' className="info-input" />
              </Col>
            </Row>
            <div className="info-input">
              <Input.TextArea placeholder='Here is you content.' style={{minHeight: 128}}/>
            </div>
            <div className="info-input"  >
              <div style={{display:'flex',justifyContent:'right'}}>
                < div style={{paddingRight:5}}>
                  <Button key='0' type='secondary'>
                    Cancel
                  </Button>
                </div>
                < div style={{paddingLeft:5}}>
                  <Button key='1' type='primary'>
                    Reply
                  </Button>
                </div>
              </div>
            </div>
          </div>
        }
        style={
          {
            paddingTop: '20px'
          }
        }
      />
      <Comment
        actions={actions}
        author={'Socrates'}
        avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/e278888093bef8910e829486fb45dd69.png~tplv-uwbnlip3yd-webp.webp'
        content={<div>Comment body content.</div>}
        datetime='1 hour'
      >
        <Comment
          actions={actions}
          author='Balzac'
          avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/9eeb1800d9b78349b24682c3518ac4a3.png~tplv-uwbnlip3yd-webp.webp'
          content={<div>Comment body content.</div>}
          datetime='1 hour'
        />
      </Comment>
      <Comment
        actions={actions}
        author={'Socrates'}
        avatar='//p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/e278888093bef8910e829486fb45dd69.png~tplv-uwbnlip3yd-webp.webp'
        content={<div>Comment body content.</div>}
        datetime='1 hour'
      />
    </div>
  )
}

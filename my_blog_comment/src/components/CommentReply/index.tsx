import React from 'react';
import { Button, Comment, Grid, Input } from '@arco-design/web-react';

const Row = Grid.Row;
const Col = Grid.Col;
const TextArea = Input.TextArea;

export default function CommentReply(){
  return (<Comment
    align='right'
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
          <TextArea placeholder='Here is you content.' style={{minHeight: 128}}/>
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
  />)
}

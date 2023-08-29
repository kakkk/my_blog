import React from 'react';
import { Comment } from '@arco-design/web-react';
import { IconMessage } from '@arco-design/web-react/icon';
import './index.css'
import CommentReply from '../CommentReply';

export default function CommentItem () {
  const [showReply, setShowReply] = React.useState(false);
  return (
    <Comment
      actions={
        <button className='custom-comment-action' onClick={()=>{setShowReply(!showReply)}}>
          <IconMessage/> Reply
        </button>
      }
      author='kakkk'
      avatar='//dn-qiniu-avatar.qbox.me/avatar/381206956e1d704103be7530dadd2e90'
      content={<div>Comment body content.</div>}
      datetime='1 hour'
    >
      {showReply?<CommentReply/>:<></>}
    </Comment>
  )
}

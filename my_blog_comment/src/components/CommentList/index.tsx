import React from 'react';
import CommentItem from '../CommonItem';
import { useComments, CommentType, CommentListItemType } from '../../CommentsContext';

export default function CommentList () {
  const { state } = useComments();

  const renderReplies = (replies: CommentType[]) => {
    return replies && replies.length > 0 ? (
      replies.map((reply) => (
        <CommentItem
          key={reply.id}
          id={reply.id}
          nickname={reply.nickname}
          avatar={reply.avatar}
          website={reply.website}
          content={reply.content}
          comment_at={reply.comment_at}
          reply_user={reply.reply_user}
        />
      ))
    ) : null;
  };

  const renderCommentWithReplies = (commentWithReplies: CommentListItemType) => {
    const { comment, replies } = commentWithReplies;

    return (
      <CommentItem
        key={comment.id}
        id={comment.id}
        nickname={comment.nickname}
        avatar={comment.avatar}
        website={comment.website}
        content={comment.content}
        comment_at={comment.comment_at}
        reply_user={comment.reply_user}
      >
        {renderReplies(replies || [])}
      </CommentItem>
    );
  };

  return (
    <div>
      {state.comments.map(renderCommentWithReplies)}
    </div>
  )
}

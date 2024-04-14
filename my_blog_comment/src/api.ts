import { httpClient } from './http';

export async function getCommentList (article_id: string) {
  return httpClient.get(`/comment/list?article_id=${article_id}`)
}

export async function commentArticle (article_id: string, nickname: string, email: string, website: string, content: string) {
  return httpClient.post(`/comment/article`, {
    article_id,
    nickname,
    email,
    website,
    content
  })
}

export async function replyComment (reply_id: string,article_id: string, nickname: string, email: string, website: string, content: string) {
  return httpClient.post(`/comment/reply`, {
    reply_id,
    article_id,
    nickname,
    email,
    website,
    content
  })
}

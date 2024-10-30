export type VideoModel = {
  id: number
  title: string
  thumbnail: string
  slug: string
  published_at: string
  likes: number
  views: number
  tags: string[]
  // author: {
  //   id: number
  //   name: string
  //   avatar: string
  // }
  video_url: string
}

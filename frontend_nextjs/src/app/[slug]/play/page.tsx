import { getVideo } from './getVideo'

export default function VideoPlayPage({
  params
}: {
  params: { slug: string }
}) {
  const video = getVideo(params.slug)

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24"></main>
  )
}

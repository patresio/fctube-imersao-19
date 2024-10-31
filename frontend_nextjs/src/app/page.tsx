import VideoCardSkeleton from '@/components/VideoCardSkeleton'
import { VideosList } from '@/components/VideosList'
import { Suspense } from 'react'

export default async function HomePage({
  searchParams: { search }
}: {
  searchParams: { search: string }
}) {
  return (
    <main className="container mx-auto px-4 py-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 md:grid-cols-3 gap-6">
        <Suspense
          fallback={Array.from({ length: 15 }).map((_, index) => (
            <VideoCardSkeleton key={index} />
          ))}
        >
          <VideosList search={search} />
        </Suspense>
      </div>
    </main>
  )
}

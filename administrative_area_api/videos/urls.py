from django.urls import path

from videos.api import videos_detail_by_id, videos_detail_by_slug, videos_list

urlpatterns = [
    path("api/videos/", videos_list, name="api_videos_list"),
    path("api/videos/<int:id>/", videos_detail_by_id, name="api_videos_detail_id"),
    path(
        "api/videos/<slug:slug>/", videos_detail_by_slug, name="api_videos_detail_slug"
    ),
]

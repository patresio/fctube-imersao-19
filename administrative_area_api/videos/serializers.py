from rest_framework import serializers

from videos.models import Video


class VideoSerializer(serializers.ModelSerializer):
    id = serializers.IntegerField()
    title = serializers.CharField()
    description = serializers.CharField()
    slug = serializers.CharField()
    published_at = serializers.DateTimeField()
    views = serializers.IntegerField(source="num_views")
    likes = serializers.IntegerField(source="num_likes")
    tags = serializers.PrimaryKeyRelatedField(many=True, read_only=True)
    thumbnail = serializers.SerializerMethodField()
    video_url = serializers.SerializerMethodField()

    def get_thumbnail(self, obj):
        # return f"{settings.ASSET_URL}{obj.thumbnail}"
        return f"http://localhost:9000{obj.thumbnail}"

    def get_video_url(self, obj):
        # return f"{settings.ASSET_URL}{obj.video_media.video_path}"
        return f"http://localhost:9000{obj.video_media.video_path}"

    class Meta:
        model = Video
        fields = [
            "id",
            "title",
            "description",
            "slug",
            "published_at",
            "views",
            "likes",
            "tags",
            "thumbnail",
            "video_url",
        ]

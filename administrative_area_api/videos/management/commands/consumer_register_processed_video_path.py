from django.core.management import BaseCommand
from kombu import Exchange, Queue

from videos.rabbitmq import create_rabbitmq_connection
from videos.services import create_video_service_factory


class Command(BaseCommand):
    help = "Register processed video path"

    def handle(self, *args, **options):
        self.stdout.write(self.style.SUCCESS("Starting consumer ..."))
        exchange = Exchange("conversion_exchange", type="direct", auto_delete=True)
        queue = Queue("finish-conversion", exchange, routing_key="finish-conversion")

        with create_rabbitmq_connection() as conn:
            with conn.Consumer(queue, callbacks=[self.process_message]):
                while True:
                    self.stdout.write(self.style.WARNING("Waiting for messages ..."))
                    conn.drain_events()

    def process_message(self, body, message):
        self.stdout.write(self.style.SUCCESS(f"Received message: {body}"))
        create_video_service_factory().register_processed_video_path(
            body["video_id"], body["path"]
        )
        self.stdout.write(
            self.style.SUCCESS(f"Deleting message: {message.delivery_tag}")
        )
        message.ack()

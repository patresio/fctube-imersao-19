from django.core.management import BaseCommand
from kombu import Exchange, Queue

from videos.rabbitmq import create_rabbitmq_connection
from videos.services import create_video_service_factory


class Command(BaseCommand):
    help = "Upload chunks to external storage"

    def handle(self, *args, **options):
        self.stdout.write(self.style.SUCCESS("Starting consumer ..."))
        exchange = Exchange("conversion_exchange")
        queue = Queue("chunks", exchange, routing_key="chunks")

        with create_rabbitmq_connection() as conn:
            with conn.Consumer(queue, callbacks=[self.process_message]):
                while True:
                    self.stdout.write(self.style.NOTICE("Waiting for messages ..."))
                    conn.drain_events()

    def process_message(self, body, message):
        self.stdout.write(self.style.SUCCESS(f"Received message: {body}"))
        create_video_service_factory().upload_chunks_to_external_storage(body["video_id"])
        self.stdout.write(
            self.style.SUCCESS(f"Deleting message: {message.delivery_tag}")
        )
        message.ack()

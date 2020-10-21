import os
from datetime import datetime
from datetime import datetime
import threading
import time
import docker

from aion.logger import initialize_logger, lprint

SERVICE_NAME = "container-image-sweeper"
CONTAINER_SIZE = 3
INTERVAL_TIME_SECOND = 3000 if os.environ.get("INTERVAL_TIME_SECOND") is None else int(os.environ.get("INTERVAL_TIME_SECOND"))

class DockerClient:
    def __init__(self):
        self.client = docker.from_env()

    def remove_image(self, image):
        for tag in image.tags:
            self.client.images.remove(tag)

    def check_exist_image(self, image_name):
        try:
            return self.client.images.get(image_name)
        except docker.errors.ImageNotFound:
            return None
    
    def images_prune(self):
        self.client.images.prune()

    def get_uniq_image_names(self, images):
        tags = []
        for image in images:
            tags.extend(image.tags)
        uniq_images = [ image.split(':')[0] for image in tags ]
        return set(uniq_images)

    def get_service_images(self):
        images = self.client.images.list()
        micro_service_images = list(filter(lambda image: 
                                    any([ "/microservice" in tag.split(':')[0] for tag in image.tags]), 
                                    images))
        return micro_service_images

    def get_base_images(self):
        images = self.client.images.list()
        base_images = list(filter(lambda image: 
                            all([ "/microservice" not in tag.split(':')[0] for tag in image.tags]),
                            images))
        return base_images
    
    def get_images(self):
        return self.client.images.list()

    def get_image_except_latest_n(self, contaierers_list, n):
        sorted_images = sorted(contaierers_list,
                            key=(lambda x: x.attrs.get('Created')), reverse=True)
        return sorted_images[n:]

    def remove_image_old(self, images):
        if not images:
            return None
        uniq_image_names = self.get_uniq_image_names(images)
        for image_name in uniq_image_names:
            image_list = self.client.images.list(image_name)
            old_images = self.get_image_except_latest_n(image_list, CONTAINER_SIZE)
            if old_images == []:
                print(image_name + f" has no old images")
            for old_image in old_images:
                try:
                    self.remove_image(old_image)
                    print(f"old_image:{old_image} removed")
                except docker.errors.APIError as e:
                    lprint(f"old_image: ${old_image} can't remove for using")

class setInterval :
    def __init__(self,interval,action) :
        self.interval=interval
        self.action=action
        self.stopEvent=threading.Event()
        thread=threading.Thread(target=self.__setInterval)
        thread.start()

    def __setInterval(self) :
        nextTime=time.time()+self.interval
        self.action()
        while not self.stopEvent.wait(nextTime-time.time()) :
            nextTime+=self.interval
            self.action()

def sweep_container_image():
    lprint(f"start on " + datetime.now().strftime("%Y/%m/%d %H:%M:%S"))
    # initialize docker client
    docker_client = DockerClient()
    # prune container images untagged
    docker_client.images_prune()
    try:
        images = docker_client.get_images()
        lprint(f"image list:{images}")
        docker_client.remove_image_old(images)
    except docker.errors.APIError as e:
        lprint(e)

def main():
    initialize_logger(SERVICE_NAME)
    setInterval(INTERVAL_TIME_SECOND, sweep_container_image)

if __name__ == "__main__":
    main()


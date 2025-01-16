import asyncio
import time
import threading
from nats.aio.client import Client as NATS

# Configuration
NATS_SERVER = "nats://localhost:4222"
SUBJECT = "functions.execute"
PAYLOAD = b"{\"image_reference\": \"fliuzzi02/custom-echo-image:latest\", \"parameter\": \"Hello, world!\"}"
TIMEOUT = 10000 # 10 seconds
NUM_THREADS = 100
NUM_REQUESTS = 1

# Statistics
total_requests = 0
successful_requests = 0
failed_requests = 0
total_latency = 0
lock = threading.Lock()

async def send_request(nc):
    global total_requests, successful_requests, failed_requests, total_latency
    start_time = time.time()
    try:
        response = await nc.request(SUBJECT, PAYLOAD, timeout=TIMEOUT)
        text = response.data.decode()
        with lock:
            # if the text contains "Hello, world!", then the request was successful
            if("Hello, world!" in text):
                successful_requests += 1
            else:
                failed_requests += 1
    except Exception as e:
        with lock:
            failed_requests += 1
    finally:
        with lock:
            total_latency += time.time() - start_time
            total_requests += 1

def worker():
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    nc = NATS()
    try:
        loop.run_until_complete(nc.connect(servers=[NATS_SERVER]))
    except Exception as e:
        print(f"Failed to connect to NATS server: {e}")
        return
    tasks = [send_request(nc) for _ in range(NUM_REQUESTS)]
    loop.run_until_complete(asyncio.gather(*tasks))
    loop.run_until_complete(nc.close())
    loop.close()

def main():
    threads = []
    start_time = time.time()

    for _ in range(NUM_THREADS):
        thread = threading.Thread(target=worker)
        thread.start()
        threads.append(thread)

    print(f"Stress testing {NUM_THREADS} threads with {NUM_REQUESTS} requests each")
    # Print a point every second until all requests are completed
    while True:
        time.sleep(1)
        print(".", end="", flush=True)
        if total_requests == NUM_THREADS * NUM_REQUESTS:
            print()
            break

    for thread in threads:
        thread.join()

    end_time = time.time()
    duration = end_time - start_time

    print(f"Total requests: {total_requests}")
    print(f"Successful requests: {successful_requests}")
    print(f"Failed requests: {failed_requests}")
    print(f"Duration: {duration:.2f} seconds")
    print(f"Success rate: {successful_requests / total_requests:.2f}")
    print(f"Requests per second: {successful_requests / duration:.2f}")
    print(f"Average seconds per request: {duration / successful_requests:.2f}")
    print(f"Average latency: {total_latency / total_requests:.2f} seconds")

if __name__ == "__main__":
    main()
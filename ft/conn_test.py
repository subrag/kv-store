import redis
import asyncio
import time
import logging
logging.basicConfig(level=logging.DEBUG)

async def call_redis(i:int):
    r = redis.Redis(host="localhost", port=8987)
    print(r.get("A"))



async def main():
    r = redis.Redis(host="localhost", port=6379)
    r.set("A", 10)
    st = time.time()
    task_list = []
    for i in range(21):
        task_list.append(loop.create_task(call_redis(str(i))))
    logging.info(f"Created task list, Total time:{time.time()-st}")
    await asyncio.wait(task_list)
    logging.info(f"Total time:{time.time()-st}")

loop = asyncio.get_event_loop()
loop.run_until_complete(main())
# asyncio.run(main())